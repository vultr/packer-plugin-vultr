// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !solaris
// +build !solaris

package filelock

import "github.com/gofrs/flock"

type Flock = flock.Flock

func New(path string) *Flock {
	return flock.New(path)
}
