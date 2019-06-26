package vultr

import (
	"context"
	"errors"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/communicator"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

// vultr constants
const (
	CustomOSID   = 159
	SnapshotOSID = 164
)

// Builder ...
type Builder struct {
	config Config
	runner multistep.Runner

	v      *lib.Client
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}
}

// Prepare ...
func (b *Builder) Prepare(raws ...interface{}) (warnings []string, err error) {
	opts := &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &b.config.interCtx,
	}
	if err = config.Decode(&b.config, opts, raws...); err != nil {
		return warnings, err
	}

	b.ctx, b.cancel = context.WithCancel(context.Background())
	b.done = make(chan struct{})

	if b.config.APIKey == "" {
		return warnings, errors.New("configuration value `api_key` not defined")
	}

	b.v = lib.NewClient(b.config.APIKey, nil)

	if b.config.RegionID == 0 {
		if b.config.RegionCode != "" {
			regions, err := b.v.GetRegions()
			if err != nil {
				return warnings, err
			}
			for _, region := range regions {
				if region.Code == b.config.RegionCode {
					b.config.RegionID = region.ID
					break
				}
			}
			if b.config.RegionID == 0 {
				return warnings, errors.New("invalid region code: " + b.config.RegionCode)
			}
		} else if b.config.RegionName != "" {
			regions, err := b.v.GetRegions()
			if err != nil {
				return warnings, err
			}
			for _, region := range regions {
				if region.Name == b.config.RegionName {
					b.config.RegionID = region.ID
					break
				}
			}
			if b.config.RegionID == 0 {
				return warnings, errors.New("invalid region name: " + b.config.RegionCode)
			}
		} else {
			return warnings, errors.New("one of `region_id` or `region_code` must be defined")
		}
	}

	if b.config.PlanID == 0 {
		if b.config.PlanName != "" {
			plans, err := b.v.GetPlans()
			if err != nil {
				return warnings, err
			}
			for _, plan := range plans {
				if plan.Name == b.config.PlanName {
					b.config.PlanID = plan.ID
					break
				}
			}
			if b.config.PlanID == 0 {
				return warnings, errors.New("invalid plan name: " + b.config.PlanName)
			}
		} else {
			return warnings, errors.New("configuration value `plan_id` not defined")
		}
	}

	if b.config.SnapshotID != "" {
		b.config.OSID = SnapshotOSID
	} else if b.config.OSID == 0 {
		if b.config.OSName != "" {
			oss, err := b.v.GetOS()
			if err != nil {
				return warnings, err
			}
			for _, os := range oss {
				if os.Name == b.config.OSName {
					b.config.OSID = os.ID
					break
				}
			}
			if b.config.OSID == 0 {
				return warnings, errors.New("invalid os name: " + b.config.OSName)
			}
		} else {
			return warnings, errors.New("configuration value `os_id` not defined")
		}
	}

	if (b.config.OSID == SnapshotOSID || b.config.OSID == CustomOSID) && b.config.Comm.SSHPassword == "" {
		return nil, errors.New("no SSH password defined for snapshot or custom OS")
	}

	if b.config.Description == "" {
		return warnings, errors.New("configuration value `description` is not defined")
	}

	if es := b.config.Comm.Prepare(&b.config.interCtx); len(es) > 0 {
		return warnings, multierror.Append(err, es...)
	}

	return warnings, nil
}

// Run ...
func (b *Builder) Run(ctx context.Context, ui packer.Ui, hook packer.Hook) (ret packer.Artifact, err error) {
	defer close(b.done)

	state := new(multistep.BasicStateBag)
	state.Put("config", b.config)
	state.Put("ctx", b.ctx)
	state.Put("hook", hook)
	state.Put("ui", ui)

	steps := []multistep.Step{
		&stepCreate{b.v},
		&communicator.StepConnect{
			Config:    &b.config.Comm,
			Host:      commHost,
			SSHConfig: sshConfig,
		},
		&common.StepProvision{},
		&stepShutdown{},
		&stepSnapshot{b.v},
	}

	b.runner = &multistep.BasicRunner{Steps: steps}
	b.runner.Run(ctx, state)

	if rawErr, ok := state.GetOk("error"); ok {
		return nil, rawErr.(error)
	}

	artifact := Artifact{
		SnapshotID: state.Get("snapshot").(lib.Snapshot).ID,
		apiKey:     b.config.APIKey,
	}
	return &artifact, nil
}

// Cancel ...
func (b *Builder) Cancel() {
	b.cancel()
	<-b.done
}
