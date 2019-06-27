package vultr

import (
	"context"
	"fmt"
	"time"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

type stepSnapshot struct {
	v *lib.Client
}

func (s *stepSnapshot) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	c := state.Get("config").(Config)
	ui := state.Get("ui").(packer.Ui)
	server := state.Get("server").(lib.Server)

	snapshot, err := s.v.CreateSnapshot(server.ID, c.Description)
	if err != nil {
		err := fmt.Errorf("Error creating snapshot: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Say(fmt.Sprintf("Waiting %ds for snapshot %s to complete...",
		int(c.stateTimeout/time.Second), snapshot.ID))

	err = waitForSnapshotState("complete", snapshot.ID, s.v, c.stateTimeout)
	if err != nil {
		err := fmt.Errorf("Error waiting for snapshot: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	state.Put("snapshot", snapshot)
	return multistep.ActionContinue
}

func (s *stepSnapshot) Cleanup(state multistep.StateBag) {
}
