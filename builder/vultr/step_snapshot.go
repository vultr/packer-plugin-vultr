package vultr

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/vultr/govultr/v3"
)

type stepCreateSnapshot struct {
	client *govultr.Client
}

// Run provides the step create snapshot run functionality
func (s *stepCreateSnapshot) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	c := state.Get("config").(*Config)
	ui := state.Get("ui").(packer.Ui)
	instance := state.Get("server").(*govultr.Instance)

	config := &oauth2.Config{}
	ts := config.TokenSource(ctx, &oauth2.Token{AccessToken: c.APIKey})
	s.client = govultr.NewClient(oauth2.NewClient(ctx, ts))

	snapshotReq := &govultr.SnapshotReq{
		InstanceID:  instance.ID,
		Description: c.Description,
	}
	snapshot, _, err := s.client.Snapshot.Create(ctx, snapshotReq)
	if err != nil {
		errOut := fmt.Errorf("error creating snapshot: %s", err)
		state.Put("error", errOut)
		ui.Error(errOut.Error())
		return multistep.ActionHalt
	}

	ui.Say(fmt.Sprintf("Waiting %ds for snapshot %s to complete...",
		int(c.stateTimeout/time.Second), snapshot.ID))

	err = waitForSnapshotState("complete", snapshot.ID, s.client, c.stateTimeout)
	if err != nil {
		errOut := fmt.Errorf("error waiting for snapshot: %s", err)
		state.Put("error", errOut)
		ui.Error(errOut.Error())
		return multistep.ActionHalt
	}

	state.Put("snapshot", snapshot)
	return multistep.ActionContinue
}

// Cleanup provides the step create snapshot cleanup functionality
func (s *stepCreateSnapshot) Cleanup(state multistep.StateBag) {
}
