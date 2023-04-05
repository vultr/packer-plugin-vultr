package vultr

import (
	"context"
	"fmt"
	"log"

	registryimage "github.com/hashicorp/packer-plugin-sdk/packer/registry/image"

	"github.com/vultr/govultr/v3"
)

type Artifact struct {
	// The ID of the snapshot
	SnapshotID string

	// The Description of the snapshot
	Description string

	// The client for making changes
	client *govultr.Client

	// config definition from the builder
	config *Config

	// State data used by HCP container registry
	StateData map[string]interface{}
}

func (a *Artifact) BuilderId() string {
	return BuilderID
}

func (a *Artifact) Files() []string {
	return nil
}

func (a *Artifact) Id() string {
	return a.SnapshotID
}

func (a *Artifact) String() string {
	return fmt.Sprintf("Vultr Snapshot: %s (%s)", a.Description, a.SnapshotID)
}

func (a *Artifact) State(name string) interface{} {
	if name == registryimage.ArtifactStateURI {
		img, err := registryimage.FromArtifact(a,
			registryimage.WithID(a.SnapshotID),
			registryimage.WithProvider("Vultr"),
			registryimage.WithRegion(a.config.RegionID),
		)

		if err != nil {
			log.Printf("[DEBUG] error encountered when creating a registry image %v", err)
			return nil
		}
		return img
	}
	return a.StateData[name]
}

func (a *Artifact) Destroy() error {
	log.Printf("Destroying Vultr Snapshot: %s (%s)", a.SnapshotID, a.Description)
	err := a.client.Snapshot.Delete(context.Background(), a.SnapshotID)
	return err
}
