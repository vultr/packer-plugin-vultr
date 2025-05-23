//go:generate packer-sdc mapstructure-to-hcl2 -type Config
package vultr

import (
	"errors"
	"fmt"
	"os"
	"time"

	common "github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

const (
	defaultStateTimeout = 10 * time.Minute
)

// Config provides the config struct
type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	Comm                communicator.Config `mapstructure:",squash"`
	ctx                 interpolate.Context

	APIKey string `mapstructure:"api_key"`

	Description string `mapstructure:"snapshot_description"`
	RegionID    string `mapstructure:"region_id"`
	PlanID      string `mapstructure:"plan_id"`
	OSID        int    `mapstructure:"os_id"`
	SnapshotID  string `mapstructure:"snapshot_id"`
	ISOID       string `mapstructure:"iso_id"`
	ISOURL      string `mapstructure:"iso_url"`
	AppID       int    `mapstructure:"app_id"`
	ImageID     string `mapstructure:"image_id"`

	EnableIPV6 bool     `mapstructure:"enable_ipv6"`
	EnableVPC  bool     `mapstructure:"enable_vpc"`
	EnableVPC2 bool     `mapstructure:"enable_vpc2"`
	ScriptID   string   `mapstructure:"script_id"`
	SSHKeyIDs  []string `mapstructure:"ssh_key_ids"`
	Label      string   `mapstructure:"instance_label"`
	UserData   string   `mapstructure:"userdata"`
	Hostname   string   `mapstructure:"hostname"`
	Tags       []string `mapstructure:"tags"`

	RawStateTimeout string `mapstructure:"state_timeout"`

	createTempSSHPair bool

	stateTimeout time.Duration
	interCtx     interpolate.Context
}

// Prepare provides the config prepare functionality
func (c *Config) Prepare(raws ...interface{}) error { //nolint:gocyclo
	if err := config.Decode(c, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &c.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{
				"run_command",
			},
		},
	}, raws...); err != nil {
		return err
	}

	var errs *packer.MultiError

	if c.APIKey == "" {
		// Default to environment variable for api_key, if it exists
		c.APIKey = os.Getenv("VULTR_API_KEY")
		if c.APIKey == "" {
			errs = packer.MultiErrorAppend(errs, errors.New("vultr api_key is required"))
		}
	}

	if c.Description == "" {
		def, err := interpolate.Render("packer-{{timestamp}}", nil)
		if err != nil {
			errs = packer.MultiErrorAppend(errs, fmt.Errorf("unable to render snapshot description: %s", err))
		} else {
			c.Description = def
		}
	}

	if c.Label == "" {
		def, err := interpolate.Render("packer-{{timestamp}}", nil)
		if err != nil {
			errs = packer.MultiErrorAppend(errs, fmt.Errorf("unable to render label: %s", err))
		} else {
			c.Label = def
		}
	}

	if c.RegionID == "" {
		errs = packer.MultiErrorAppend(errs, errors.New("region_id is required"))
	}

	if c.PlanID == "" {
		errs = packer.MultiErrorAppend(errs, errors.New("plan_id is required"))
	}

	imageConfig := []bool{(c.AppID != 0), (c.ISOID != ""), (c.ISOURL != ""), (c.OSID != 0), (c.SnapshotID != "")}
	imageDefined := false
	for _, isDefined := range imageConfig {
		if isDefined {
			if imageDefined {
				errs = packer.MultiErrorAppend(errs, errors.New("only set one of the following: `app_id`, `iso_id`, `iso_url`, `os_id`, `snapshot_id`"))
				break
			}
			imageDefined = true
		}
	}

	if c.SnapshotID != "" || c.ISOID != "" || c.ISOURL != "" {
		c.createTempSSHPair = false
		if c.Comm.SSHPassword == "" && c.Comm.SSHPrivateKeyFile == "" {
			errs = packer.MultiErrorAppend(errs, errors.New("define `ssh_password` or `ssh_private_key_file` for snapshot or custom OS"))
		}
	} else {
		c.createTempSSHPair = true
	}

	if c.RawStateTimeout == "" {
		c.stateTimeout = defaultStateTimeout
	} else {
		if stateTimeout, err := time.ParseDuration(c.RawStateTimeout); err == nil {
			c.stateTimeout = stateTimeout
		} else {
			errs = packer.MultiErrorAppend(errs, fmt.Errorf("unable to parse state timeout: %s", err))
		}
	}

	if es := c.Comm.Prepare(&c.interCtx); len(es) > 0 {
		errs = packer.MultiErrorAppend(errs, es...)
	}

	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}

	packer.LogSecretFilter.Set(c.APIKey)

	return nil
}
