package vultr

import (
	"context"
	"time"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"golang.org/x/crypto/ssh"
)

type stepShutdown struct{}

func (s *stepShutdown) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	c := state.Get("config").(Config)
	ui := state.Get("ui").(packer.Ui)
	server := state.Get("server").(lib.Server)

	ui.Say("Preparing the server for a graceful shutdown...")
	config, err := sshConfig(state)
	if err != nil {
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	client, err := ssh.Dial("tcp", server.MainIP+":22", config)
	if err != nil {
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	session, err := client.NewSession()
	if err != nil {
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	if c.ShutdownCommand == "" {
		c.ShutdownCommand = "shutdown -P now"
	}
	ui.Say("Shutting down server...")
	err = session.Run(c.ShutdownCommand)
	if err != nil {
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Say("Sleeping to ensure that server is shut down...")
	time.Sleep(3 * time.Second)

	return multistep.ActionContinue
}

func (s *stepShutdown) Cleanup(state multistep.StateBag) {
}
