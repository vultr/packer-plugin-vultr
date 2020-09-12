package vultr

import (
	"os"
	"testing"

	builderT "github.com/hashicorp/packer/helper/builder/testing"
)

func TestBuilderAcc_basic(t *testing.T) {
	builderT.Test(t, builderT.TestCase{
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
        "os_id": 127,
        "ssh_username": "root",
		"state_timeout": "20m"
	}]
}
`
