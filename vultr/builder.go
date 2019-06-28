package vultr

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/communicator"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

// Special OS IDs
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
	c := &b.config
	opts := &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &b.config.interCtx,
	}
	if err = config.Decode(c, opts, raws...); err != nil {
		return warnings, err
	}

	b.ctx, b.cancel = context.WithCancel(context.Background())
	b.done = make(chan struct{})

	if c.Description == "" {
		return warnings, errors.New("configuration value `description` is not defined")
	}

	if c.APIKey == "" {
		return warnings, errors.New("configuration value `api_key` not defined")
	}

	b.v = lib.NewClient(c.APIKey, nil)

	if c.RegionID == 0 && c.RegionName == "" && c.RegionCode == "" {
		return warnings, errors.New("please define one of: `region_id`, `region_name`, `region_code`")
	}
	if c.RegionName != "" && c.RegionCode != "" {
		return warnings, errors.New("you can define either `region_name` or `region_code`, not both")
	}
	if c.RegionName != "" {
		if c.RegionID != 0 {
			return warnings, errors.New("you can define either `region_name` or `region_id`, not both")
		}
		regions, err := b.v.GetRegions()
		if err != nil {
			return warnings, err
		}
		for _, r := range regions {
			if r.Name == c.RegionName {
				c.RegionID = r.ID
				break
			}
		}
		if c.RegionID == 0 {
			return warnings, fmt.Errorf("cannot find region with name %q", c.RegionName)
		}
	}
	if c.RegionCode != "" {
		if c.RegionID != 0 {
			return warnings, errors.New("you can define either `region_code` or `region_id`, not both")
		}
		regions, err := b.v.GetRegions()
		if err != nil {
			return warnings, err
		}
		for _, r := range regions {
			if r.Code == c.RegionCode {
				c.RegionID = r.ID
				break
			}
		}
		if c.RegionID == 0 {
			return warnings, fmt.Errorf("cannot find region with code %q", c.RegionCode)
		}
	}

	if c.PlanID == 0 && c.PlanName == "" {
		return warnings, errors.New("please define one of: `plan_id`, `plan_name`")
	}
	if c.PlanName != "" {
		if c.PlanID != 0 {
			return warnings, errors.New("you can define either `plan_name` or `plan_id`, not both")
		}
		plans, err := b.v.GetPlans()
		if err != nil {
			return warnings, err
		}
		for _, p := range plans {
			if p.Name == c.PlanName {
				c.PlanID = p.ID
				break
			}
		}
		if c.PlanID == 0 {
			return warnings, fmt.Errorf("cannot find plan with name %q", c.PlanName)
		}
	}

	if c.SnapshotID == "" && c.SnapshotName == "" && c.OSID == 0 && c.OSName == "" {
		return warnings, errors.New("please define one of: `os_id`, `os_name`, `snapshot_id`, `snapshot_name`")
	}

	if c.OSName != "" {
		if c.OSID != 0 {
			return warnings, errors.New("you can define either `os_id` or `os_name`, not both")
		}
		oses, err := b.v.GetOS()
		if err != nil {
			return warnings, err
		}
		for _, os := range oses {
			if os.Name == c.OSName {
				c.OSID = os.ID
				break
			}
		}
		if c.OSID == 0 {
			return warnings, fmt.Errorf("OS name %q is invalid", c.OSName)
		}
	}

	if c.SnapshotName != "" {
		if c.SnapshotID != "" {
			return warnings, errors.New("you can define either `snapshot_id` or `snapshot_name`, not both")
		}
		snapshots, err := b.v.GetSnapshots()
		if err != nil {
			return warnings, err
		}
		for _, s := range snapshots {
			if s.Status != "complete" {
				continue // snapshot is not ready, skip
			}
			if s.Description == c.SnapshotName {
				if c.SnapshotID != "" {
					return warnings, fmt.Errorf("snapshot name %q is ambiguous", c.SnapshotName)
				}
				c.SnapshotID = s.ID
			}
		}
		if c.SnapshotID == "" {
			return warnings, fmt.Errorf("cannot find snapshot with name %q", c.SnapshotName)
		}
	}

	if c.OSID != 0 && c.SnapshotID != "" {
		return warnings, errors.New("you can define either OS or Snapshot, not both")
	}
	if c.SnapshotID != "" {
		c.OSID = SnapshotOSID
	}

	if (c.OSID == SnapshotOSID || c.OSID == CustomOSID) && c.Comm.SSHPassword == "" && c.Comm.SSHPrivateKeyFile == "" {
		return warnings, errors.New("either `ssh_password` or `ssh_private_key_file` must be defined for snapshot or custom OS")
	}

	if c.RawStateTimeout == "" {
		c.RawStateTimeout = "10m"
	}
	if c.stateTimeout, err = time.ParseDuration(c.RawStateTimeout); err != nil {
		return warnings, errors.New("invalid state timeout: " + c.RawStateTimeout)
	}

	if es := c.Comm.Prepare(&c.interCtx); len(es) > 0 {
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
