// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package communicator

// WinRMConfig is configuration that can be returned at runtime to
// dynamically configure WinRM.
type WinRMConfig struct {
	Username string
	Password string
}
