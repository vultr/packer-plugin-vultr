package vultr

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/vultr/govultr/v2"
)

type stepCreateServer struct {
	client *govultr.Client
}

func (s *stepCreateServer) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	c := state.Get("config").(*Config)
	ui := state.Get("ui").(packer.Ui)

	ui.Say("Creating Vultr instance...")

	tempKey := state.Get("temp_ssh_key_id").(string)
	keys := append(c.SSHKeyIDs, tempKey)

	// check if ISO ID should be populated by the result of creating an ISO in step_create_iso via 'iso_url'
	iso_id := c.ISOID

	iso, iso_ok := state.GetOk("iso")
	if iso_ok {
		iso_id = iso.(*govultr.ISO).ID
	}

	instanceReq := &govultr.InstanceCreateReq{
		ISOID:                iso_id,
		SnapshotID:           c.SnapshotID,
		OsID:                 c.OSID,
		Region:               c.RegionID,
		Plan:                 c.PlanID,
		AppID:                c.AppID,
		ScriptID:             c.ScriptID,
		ImageID:              c.ImageID,
		EnableIPv6:           govultr.BoolToBoolPtr(c.EnableIPV6),
		EnablePrivateNetwork: govultr.BoolToBoolPtr(c.EnablePrivateNetwork),
		Label:                c.Label,
		SSHKeys:              keys,
		UserData:             c.UserData,
		ActivationEmail:      govultr.BoolToBoolPtr(false),
		Hostname:             c.Hostname,
		Tag:                  c.Tag,
	}

	instance, err := s.client.Instance.Create(ctx, instanceReq)
	if err != nil {
		err = errors.New("Error creating server: " + err.Error())
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	state.Put("default_password", instance.DefaultPassword)

	// wait until server is running
	ui.Say(fmt.Sprintf("Waiting %ds for server %s to power on...",
		int(c.stateTimeout/time.Second), instance.ID))

	err = waitForServerState("active", "running", instance.ID, s.client, c.stateTimeout)
	if err != nil {
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	if instance, err = s.client.Instance.Get(context.Background(), instance.ID); err != nil {
		err := fmt.Errorf("Error getting server: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	state.Put("server", instance)
	state.Put("server_ip", instance.MainIP)
	state.Put("server_id", instance.ID)

	return multistep.ActionContinue
}

func (s *stepCreateServer) Cleanup(state multistep.StateBag) {
	server, ok := state.GetOk("server")
	if !ok {
		return
	}

	ui := state.Get("ui").(packer.Ui)
	instance := server.(*govultr.Instance)

	// If an ISO was uploaded as part of this build, detach from the instance before destroying
	if iso, ok := state.GetOk("iso"); ok {
		iso_status, iso_status_err := s.client.Instance.ISOStatus(context.Background(), instance.ID)
		if iso_status_err != nil {
			state.Put("error", iso_status_err)
		}

		if iso_status.State == "isomounted" && iso_status.IsoID == iso.(*govultr.ISO).ID {
			if detach_err := s.client.Instance.DetachISO(context.Background(), instance.ID); detach_err != nil {
				state.Put("error", detach_err)
			}
		}
	}

	ui.Say("Destroying server " + instance.ID)
	if err := s.client.Instance.Delete(context.Background(), instance.ID); err != nil {
		state.Put("error", err)
	}
}
