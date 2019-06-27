package vultr

import (
	"context"
	"fmt"
	"time"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

type stepCreate struct {
	v *lib.Client
}

func (s *stepCreate) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	c := state.Get("config").(Config)
	ui := state.Get("ui").(packer.Ui)

	opts := &lib.ServerOptions{
		Script:               c.ScriptID,
		DontNotifyOnActivate: true,
	}
	if c.OSID == SnapshotOSID {
		opts.Snapshot = c.SnapshotID
	}

	serverName := "packing..." + c.Description
	server, err := s.v.CreateServer(serverName, c.RegionID, c.PlanID, c.OSID, opts)
	if err != nil {
		err := fmt.Errorf("Error creating server: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	state.Put("server-creation-time", time.Now())
	serverID := server.ID
	state.Put("server_id", serverID)

	ui.Say(fmt.Sprintf("Waiting %ds for server %s to power on...",
		int(c.stateTimeout/time.Second), serverID))

	err = waitForServerState("active", "running", serverID, s.v, c.stateTimeout)
	if err != nil {
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	server, err = s.v.GetServer(serverID)
	if err != nil {
		err := fmt.Errorf("Error querying server: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	state.Put("server", server)
	return multistep.ActionContinue
}

func (s *stepCreate) Cleanup(state multistep.StateBag) {
	ui := state.Get("ui").(packer.Ui)
	serverID := state.Get("server_id").(string)
	startTime := state.Get("server-creation-time").(time.Time)

	wait := (5 * time.Minute) - time.Since(startTime)
	if wait > 0 {
		ui.Say("Vultr requires you to wait 5 minutes before destroying a server, we have " + wait.String() + " left...")
		time.Sleep(wait)
	}

	ui.Say("Destroying server " + serverID)
	if err := s.v.DeleteServer(serverID); err != nil {
		state.Put("error", err)
	}
}
