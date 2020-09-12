package vultr

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
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

	instanceReq := &govultr.InstanceReq{
		ISOID:                c.ISOID,
		SnapshotID:           c.SnapshotID,
		OsID:                 c.OSID,
		Region:               c.RegionID,
		Plan:                 c.PlanID,
		AppID:                c.AppID,
		ScriptID:             c.ScriptID,
		EnableIPv6:           c.EnableIPV6,
		EnablePrivateNetwork: c.EnablePrivateNetwork,
		Label:                c.Label,
		SSHKey:               keys,
		UserData:             c.UserData,
		ActivationEmail:      false,
		Hostname:             c.Hostname,
		Tag:                  c.Tag,
	}
	ui.Say(fmt.Sprintf("%v", instanceReq))
	instance, err := s.client.Instance.Create(ctx, instanceReq)
	if err != nil {
		err = errors.New("Error creating server: " + err.Error())
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

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
	instanceID := server.(*govultr.Instance).ID

	ui.Say("Destroying server " + instanceID)
	if err := s.client.Instance.Delete(context.Background(), instanceID); err != nil {
		state.Put("error", err)
	}
}
