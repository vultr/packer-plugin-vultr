package vultr

import (
	"context"
	"time"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

// shutdown delays
const (
	ShutdownDelaySec = 10
	ShutdownWaitSec  = 10
)

type stepShutdown struct{}

func (s *stepShutdown) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	c := state.Get("config").(Config)
	ui := state.Get("ui").(packer.Ui)

	ui.Say("Performing graceful shutdown...")
	time.Sleep(ShutdownDelaySec * time.Second)

	comm := state.Get("communicator").(packer.Communicator)

	cmd := &packer.RemoteCmd{
		Command: c.ShutdownCommand,
	}
	if cmd.Command == "" {
		cmd.Command = "shutdown -P now"
	}

	if err := comm.Start(ctx, cmd); err != nil {
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	cmd.Wait()

	if cmd.ExitStatus() == packer.CmdDisconnect {
		ui.Say("Server is successfully shut down")
		time.Sleep(ShutdownDelaySec * time.Second)
	} else {
		ui.Say("Sleeping to ensure that server is shut down...")
	}
	return multistep.ActionContinue
}

func (s *stepShutdown) Cleanup(state multistep.StateBag) {
}
