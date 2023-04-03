// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build darwin || freebsd || linux || netbsd || openbsd || solaris
// +build darwin freebsd linux netbsd openbsd solaris

package packer

import (
	"os"
	"path/filepath"
)

func getDefaultCacheDir() string {
	var defaultConfigFileDir string

	if xdgCacheHome := os.Getenv("XDG_CACHE_HOME"); xdgCacheHome != "" {
		return filepath.Join(xdgCacheHome, "packer")
	}

	homeDir := os.Getenv("HOME")
	defaultConfigFileDir = filepath.Join(homeDir, ".cache", "packer")

	return defaultConfigFileDir
}
