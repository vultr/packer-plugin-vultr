package vultr

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/multistep/commonsteps"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/vultr/govultr/v2"
)

// BuilderID is the unique ID for the builder
const BuilderID = "packer.vultr"

// Builder ...
type Builder struct {
	config Config
	runner multistep.Runner
}

func (b *Builder) ConfigSpec() hcldec.ObjectSpec { return b.config.FlatMapstructure().HCL2Spec() }

func (b *Builder) Prepare(raws ...interface{}) ([]string, []string, error) {
	warnings, errs := b.config.Prepare(raws...)
	if errs != nil {
		return nil, warnings, errs
	}
	return nil, nil, nil
}

func (b *Builder) Run(ctx context.Context, ui packer.Ui, hook packer.Hook) (ret packer.Artifact, err error) {
	ui.Say("Running Vultr builder...")

	client := newVultrClient(b.config.APIKey)

	state := new(multistep.BasicStateBag)
	state.Put("config", &b.config)
	state.Put("hook", hook)
	state.Put("ui", ui)

	steps := []multistep.Step{
		&stepCreateSSHKey{
			client:       client,
			Debug:        b.config.PackerDebug,
			DebugKeyPath: fmt.Sprintf("vultr_%s.pem", b.config.PackerBuildName),
		},
		&stepCreateServer{client},
		&communicator.StepConnect{
			Config:    &b.config.Comm,
			Host:      communicator.CommHost(b.config.Comm.Host(), "server_ip"),
			SSHConfig: b.config.Comm.SSHConfigFunc(),
		},
		&commonsteps.StepProvision{},
		&commonsteps.StepCleanupTempKeys{
			Comm: &b.config.Comm,
		},
		&stepShutdown{client},
		&stepCreateSnapshot{client},
	}

	b.runner = commonsteps.NewRunner(steps, b.config.PackerConfig, ui)
	b.runner.Run(ctx, state)

	if rawErr, ok := state.GetOk("error"); ok {
		return nil, rawErr.(error)
	}

	if _, ok := state.GetOk(multistep.StateCancelled); ok {
		return nil, errors.New("build was cancelled")
	}

	if _, ok := state.GetOk(multistep.StateHalted); ok {
		return nil, errors.New("build was halted")
	}

	if _, ok := state.GetOk("snapshot"); !ok {
		return nil, errors.New("cannot find snapshot in state")
	}

	snapshot := state.Get("snapshot").(*govultr.Snapshot)
	artifact := &Artifact{
		SnapshotID:  snapshot.ID,
		Description: snapshot.Description,
		client:      client,
	}

	return artifact, nil
}
