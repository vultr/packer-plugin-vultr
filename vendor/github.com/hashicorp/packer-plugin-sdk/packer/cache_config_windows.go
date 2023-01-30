// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build windows
// +build windows

package packer

const (
	defaultConfigFile = "packer_cache"
)

func getDefaultCacheDir() string {
	return defaultConfigFile
}
