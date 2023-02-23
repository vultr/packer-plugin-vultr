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

type stepCreateISO struct {
	client *govultr.Client
}

func (s *stepCreateISO) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	c := state.Get("config").(*Config)
	ui := state.Get("ui").(packer.Ui)

	if len(c.ISOURL) > 0 {
		ui.Say("Creating ISO in Vultr account...")

		isoReq := &govultr.ISOReq{
			URL: c.ISOURL,
		}

		iso, err := s.client.ISO.Create(ctx, isoReq)
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

		if _, err = s.client.ISO.Get(context.Background(), iso.ID); err != nil {
			err := fmt.Errorf("error getting ISO: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
	}

	return multistep.ActionContinue
}

func (s *stepCreateISO) Cleanup(state multistep.StateBag) {
	iso, iso_ok := state.GetOk("iso")
	if !iso_ok {
		return
	}

	ui := state.Get("ui").(packer.Ui)
	isoID := iso.(*govultr.ISO).ID

	ui.Say("Destroying ISO " + isoID)
	if err := s.client.ISO.Delete(context.Background(), isoID); err != nil {
		state.Put("error", err)
	}
}
