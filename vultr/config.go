package vultr

import (
	"time"

	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/communicator"
	"github.com/hashicorp/packer/template/interpolate"
)

// Config ...
type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	Comm                communicator.Config `mapstructure:",squash"`
	APIKey              string              `mapstructure:"api_key"`
	Description         string              `mapstructure:"description"`
	RegionID            int                 `mapstructure:"region_id"`
	RegionName          string              `mapstructure:"region_name"`
	RegionCode          string              `mapstructure:"region_code"`
	PlanID              int                 `mapstructure:"plan_id"`
	PlanName            string              `mapstructure:"plan_name"`
	OSID                int                 `mapstructure:"os_id"`
	OSName              string              `mapstructure:"os_name"`
	IPv6                bool                `mapstructure:"ipv6"`
	PrivateNetworking   bool                `mapstructure:"private_networking"`
	ScriptID            int                 `mapstructure:"startup_script_id"`
	ScriptName          string              `mapstructure:"startup_script_name"`
	SSHKeyID            string              `mapstructure:"sshkey_id"`
	SSHKeyName          string              `mapstructure:"sshkey_name"`
	SnapshotID          string              `mapstructure:"snapshot_id"`
	SnapshotName        string              `mapstructure:"snapshot_name"`
	ShutdownCommand     string              `mapstructure:"shutdown_command"`
	RawStateTimeout     string              `mapstructure:"state_timeout"`
	stateTimeout        time.Duration
	interCtx            interpolate.Context
}
