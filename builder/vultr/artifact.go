package vultr

import (
	"context"
	"fmt"
	"log"

	registryimage "github.com/hashicorp/packer-plugin-sdk/packer/registry/image"

	"github.com/vultr/govultr/v3"
)

// Artifact provides the artifact struct
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

// BuilderId provides the builder ID
func (a *Artifact) BuilderId() string { //nolint:stylecheck
	return BuilderID
}

// Files provides nil
func (a *Artifact) Files() []string {
	return nil
}

// Id provides the snapshot ID
func (a *Artifact) Id() string { //nolint:stylecheck
	return a.SnapshotID
}

// String provides the snapshot description and ID in a string
func (a *Artifact) String() string {
	return fmt.Sprintf("Vultr Snapshot: %s (%s)", a.Description, a.SnapshotID)
}

// State provides the artifact state
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

// Destroy destroys the artifact snapshot
func (a *Artifact) Destroy() error {
	log.Printf("Destroying Vultr Snapshot: %s (%s)", a.SnapshotID, a.Description)
	err := a.client.Snapshot.Delete(context.Background(), a.SnapshotID)
	return err
}
