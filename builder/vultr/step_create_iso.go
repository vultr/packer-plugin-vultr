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

type stepCreateISO struct {
	client *govultr.Client
}

// Run provides the stepCreateISO run function
func (s *stepCreateISO) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	c := state.Get("config").(*Config)
	ui := state.Get("ui").(packer.Ui)

	if c.ISOURL != "" {
		ui.Say("Creating ISO in Vultr account...")

		isoReq := &govultr.ISOReq{
			URL: c.ISOURL,
		}

		iso, _, err := s.client.ISO.Create(ctx, isoReq)
		if err != nil {
			err = errors.New("Error creating ISO: " + err.Error())
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}

		state.Put("iso", iso)

		ui.Say(fmt.Sprintf("Waiting %ds for ISO %s (%s) to complete uploading...",
			int(c.stateTimeout/time.Second), iso.FileName, iso.ID))

		err = waitForISOState("complete", iso.ID, s.client, c.stateTimeout)
		if err != nil {
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}

		if _, _, err = s.client.ISO.Get(context.Background(), iso.ID); err != nil {
			errOut := fmt.Errorf("error getting ISO: %s", err)
			state.Put("error", errOut)
			ui.Error(errOut.Error())
			return multistep.ActionHalt
		}
	}

	return multistep.ActionContinue
}

// Cleanup provides the step create ISO cleanup function
func (s *stepCreateISO) Cleanup(state multistep.StateBag) {
	iso, isoOK := state.GetOk("iso")
	if !isoOK {
		return
	}

	ui := state.Get("ui").(packer.Ui)
	isoID := iso.(*govultr.ISO).ID

	ui.Say("Destroying ISO " + isoID)
	if err := s.client.ISO.Delete(context.Background(), isoID); err != nil {
		state.Put("error", err)
	}
}
