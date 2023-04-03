// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// build solaris

package filelock

// Flock is a noop on solaris for now.
// TODO(azr): PR github.com/gofrs/flock for this.
type Flock = Noop

func New(string) *Flock {
	return &Flock{}
}
