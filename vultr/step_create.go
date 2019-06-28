package vultr

import (
	"context"
	"fmt"
	"time"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

// Delays
const (
	StartupScriptDelaySec = 20
	ServerDestroyWaitMin  = 3
)

type stepCreate struct {
	v *lib.Client
}

func (s *stepCreate) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	c := state.Get("config").(Config)
	ui := state.Get("ui").(packer.Ui)

	opts := &lib.ServerOptions{
		Script:            c.ScriptID,
		IPV6:              c.IPv6,
		PrivateNetworking: c.PrivateNetworking,
		//DontNotifyOnActivate: true,
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

	if c.ScriptID != 0 {
		ui.Say(fmt.Sprintf("Delay %ds for startup script %v to complete...", StartupScriptDelaySec, c.ScriptID))
		time.Sleep(StartupScriptDelaySec * time.Second)
	}

	return multistep.ActionContinue
}

func (s *stepCreate) Cleanup(state multistep.StateBag) {
	ui := state.Get("ui").(packer.Ui)
	serverID := state.Get("server_id").(string)
	startTime := state.Get("server-creation-time").(time.Time)

	wait := (ServerDestroyWaitMin * time.Minute) - time.Since(startTime)
	if wait > 0 {
		ui.Say(fmt.Sprintf("Vultr requires you to wait %dm before destroying a server, we have %v left...", ServerDestroyWaitMin, wait))
		time.Sleep(wait)
	}

	ui.Say("Destroying server " + serverID)
	if err := s.v.DeleteServer(serverID); err != nil {
		state.Put("error", err)
	}
}
