package vultr

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
)

func TestBuilderAcc_basic(t *testing.T) {
	if skip := testAccPreCheck(t); skip == true {
		return
	}

	acctest.TestPlugin(t, &acctest.PluginTestCase{
		Name:     "test-vultr-builder-basic",
		Template: testBuilderAccBasic,
	})
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
        "os_id": 167,
        "ssh_username": "root",
		"state_timeout": "60m"
	}]
}
`
