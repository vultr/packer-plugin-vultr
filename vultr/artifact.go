package vultr

import "github.com/JamesClonk/vultr/lib"

// Artifact ...
type Artifact struct {
	SnapshotID string
	apiKey     string
}

// BuilderId ...
func (a *Artifact) BuilderId() string {
	return "vultr"
}

// Files ...
func (a *Artifact) Files() []string {
	return nil
}

// Id ...
func (a *Artifact) Id() string {
	return a.SnapshotID
}

// String ...
func (a *Artifact) String() string {
	return "Snapshot: " + a.SnapshotID
}

// State ...
func (a *Artifact) State(name string) interface{} {
	return nil
}

// Destroy ...
func (a *Artifact) Destroy() error {
	return lib.NewClient(a.apiKey, nil).DeleteSnapshot(a.SnapshotID)
}
