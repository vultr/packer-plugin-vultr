package vultr

import (
	"os"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
)

func TestBuilderAcc_basic(t *testing.T) {
	acctest.Test(t, acctest.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Builder:  &Builder{},
		Template: testBuilderAccBasic,
	})
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("VULTR_API_KEY"); v == "" {
		t.Fatal("VULTR_API_KEY must be set for acceptance tests")
	}
}

const testBuilderAccBasic = `
{
	"builders": [{
		"type": "test",
		"snapshot_description": "packer-test-snapshot",
        "region_id": "ewr",
        "plan_id": "vc2-1c-1gb",
        "os_id": 167,
        "ssh_username": "root",
		"state_timeout": "20m"
	}]
}
`
