package vultr

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
)

func TestBuilderAcc(t *testing.T) {
	if skip := testAccPreCheck(t); skip == true {
		return
	}

	basicTestCase := &acctest.PluginTestCase{
		Name: "test-vultr-builder-basic",
		Setup: func() error {
			return nil
		},
		Teardown: func() error {
			return nil
		},
		Template: testBuilderAccBasic,
		Type:     "vultr",
		Check: func(c *exec.Cmd, logfile string) error {
			if c.ProcessState != nil {
				if c.ProcessState.ExitCode() != 0 {
					return fmt.Errorf("Bad exit code. Logfile %s", logfile)
				}
			}

			logs, err := os.Open(logfile)
			if err != nil {
				return fmt.Errorf("Unable to find %s", logfile)
			}
			defer logs.Close()

			logsBytes, err := ioutil.ReadAll(logs)
			if err != nil {
				return fmt.Errorf("Unable to read %s", logfile)
			}
			logsString := string(logsBytes)

			buildGeneratedDataLog := "vultr: Destroying server"
			if matched, _ := regexp.MatchString(buildGeneratedDataLog+".*", logsString); !matched {
				t.Fatalf("Logs do not contain expected log value %q", logsString)
			}
			return nil
		},
	}
	acctest.TestPlugin(t, basicTestCase)

	// isoTestCase := basicTestCase
	// isoTestCase.Name = "test-vultr-builder-iso"
	// isoTestCase.Template = testBuilderAccISO
	// acctest.TestPlugin(t, isoTestCase)
}

func testAccPreCheck(t *testing.T) bool {
	if os.Getenv(acctest.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf(
			"Acceptance tests skipped unless env '%s' set",
			acctest.TestEnvVar))
		return true
	}

	if v := os.Getenv("VULTR_API_KEY"); v == "" {
		t.Fatal("VULTR_API_KEY must be set for acceptance tests")
		return true
	}
	return false
}

const testBuilderAccBasic = `
{
	"builders": [{
		"type": "vultr",
		"snapshot_description": "packer-test-snapshot",
		"region_id": "ewr",
		"plan_id": "vc2-1c-1gb",
		"os_id": 477,
		"ssh_timeout": "10m",
		"ssh_username": "root",
		"state_timeout": "60m"
	}]
}
`

// Requires a specially modified ISO that has been crafted to configure SSH access + relevant configuration options
const testBuilderAccISO = `
{
	"builders": [{
		"type": "vultr",
		"snapshot_description": "packer-test-snapshot-from-iso",
		"region_id": "ewr",
		"plan_id": "vc2-1c-1gb",
		"ssh_timeout": "10m",
		"ssh_username": "root",
		"iso_url": "https://example.com/customized.iso",
		"state_timeout": "60m"
	}]
}
`
