package vultr

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/vultr/govultr/v3"
)

type stepCreateServer struct {
	client *govultr.Client
}

// Run provides the step create server run functionality
func (s *stepCreateServer) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	c := state.Get("config").(*Config)
	ui := state.Get("ui").(packer.Ui)

	ui.Say("Creating Vultr instance...")

	sshKeys := c.SSHKeyIDs
	key, keyOK := state.GetOk("temp_ssh_key_id")
	if keyOK {
		sshKeys = append(sshKeys, key.(string))
	}

	// check if ISO ID should be populated by the result of creating an ISO in step_create_iso via 'iso_url'
	isoID := c.ISOID
	iso, isoOK := state.GetOk("iso")
	if isoOK {
		isoID = iso.(*govultr.ISO).ID
	}

	instanceReq := &govultr.InstanceCreateReq{
		ISOID:                isoID,
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
		SSHKeys:              sshKeys,
		UserData:             c.UserData,
		ActivationEmail:      govultr.BoolToBoolPtr(false),
		Hostname:             c.Hostname,
		Tag:                  c.Tag,
	}

	instance, _, err := s.client.Instance.Create(ctx, instanceReq)
	if err != nil {
		errOut := errors.New("Error creating server: " + err.Error())
		state.Put("error", errOut)
		ui.Error(errOut.Error())
		return multistep.ActionHalt
	}
	state.Put("default_password", instance.DefaultPassword)

	// wait until server is running
	ui.Say(fmt.Sprintf("Waiting %ds for server %s to power on...",
		int(c.stateTimeout/time.Second), instance.ID))

	if err = waitForServerState("active", "running", instance.ID, s.client, c.stateTimeout); err != nil {
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	if instance, _, err = s.client.Instance.Get(context.Background(), instance.ID); err != nil {
		errOut := fmt.Errorf("error getting server: %s", err)
		state.Put("error", errOut)
		ui.Error(errOut.Error())
		return multistep.ActionHalt
	}
	state.Put("server", instance)
	state.Put("server_ip", instance.MainIP)
	state.Put("server_id", instance.ID)

	return multistep.ActionContinue
}

// Cleanup provides the step create server cleanup functionality
func (s *stepCreateServer) Cleanup(state multistep.StateBag) {
	server, ok := state.GetOk("server")
	if !ok {
		return
	}

	ui := state.Get("ui").(packer.Ui)
	instance := server.(*govultr.Instance)

	// If an ISO was uploaded as part of this build, detach from the instance before destroying
	if iso, ok := state.GetOk("iso"); ok {
		isoStatus, _, err := s.client.Instance.ISOStatus(context.Background(), instance.ID)
		if err != nil {
			state.Put("error", err)
		}

		if isoStatus.State == "isomounted" && isoStatus.IsoID == iso.(*govultr.ISO).ID {
			if _, err := s.client.Instance.DetachISO(context.Background(), instance.ID); err != nil {
				state.Put("error", err)
			}
		}
	}

	ui.Say("Destroying server " + instance.ID)
	if err := s.client.Instance.Delete(context.Background(), instance.ID); err != nil {
		state.Put("error", err)
	}
}
