// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:generate packer-sdc struct-markdown

package commonsteps

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

// An iso (CD) containing custom files can be made available for your build.
//
// By default, no extra CD will be attached. All files listed in this setting
// get placed into the root directory of the CD and the CD is attached as the
// second CD device.
//
// This config exists to work around modern operating systems that have no
// way to mount floppy disks, which was our previous go-to for adding files at
// boot time.
type CDConfig struct {
	// A list of files to place onto a CD that is attached when the VM is
	// booted. This can include either files or directories; any directories
	// will be copied onto the CD recursively, preserving directory structure
	// hierarchy. Symlinks will have the link's target copied into the directory
	// tree on the CD where the symlink was. File globbing is allowed.
	//
	// Usage example (JSON):
	//
	// ```json
	// "cd_files": ["./somedirectory/meta-data", "./somedirectory/user-data"],
	// "cd_label": "cidata",
	// ```
	//
	// Usage example (HCL):
	//
	// ```hcl
	// cd_files = ["./somedirectory/meta-data", "./somedirectory/user-data"]
	// cd_label = "cidata"
	// ```
	//
	// The above will create a CD with two files, user-data and meta-data in the
	// CD root. This specific example is how you would create a CD that can be
	// used for an Ubuntu 20.04 autoinstall.
	//
	// Since globbing is also supported,
	//
	// ```hcl
	// cd_files = ["./somedirectory/*"]
	// cd_label = "cidata"
	// ```
	//
	// Would also be an acceptable way to define the above cd. The difference
	// between providing the directory with or without the glob is whether the
	// directory itself or its contents will be at the CD root.
	//
	// Use of this option assumes that you have a command line tool installed
	// that can handle the iso creation. Packer will use one of the following
	// tools:
	//
	//   * xorriso
	//   * mkisofs
	//   * hdiutil (normally found in macOS)
	//   * oscdimg (normally found in Windows as part of the Windows ADK)
	CDFiles []string `mapstructure:"cd_files"`
	// Key/Values to add to the CD. The keys represent the paths, and the values
	// contents. It can be used alongside `cd_files`, which is useful to add large
	// files without loading them into memory. If any paths are specified by both,
	// the contents in `cd_content` will take precedence.
	//
	// Usage example (HCL):
	//
	// ```hcl
	// cd_files = ["vendor-data"]
	// cd_content = {
	//   "meta-data" = jsonencode(local.instance_data)
	//   "user-data" = templatefile("user-data", { packages = ["nginx"] })
	// }
	// cd_label = "cidata"
	// ```
	CDContent map[string]string `mapstructure:"cd_content"`
	CDLabel   string            `mapstructure:"cd_label"`
}

func (c *CDConfig) Prepare(ctx *interpolate.Context) []error {
	var errs []error
	var err error

	if c.CDFiles == nil {
		c.CDFiles = make([]string, 0)
	}

	// Create new file list based on globbing.
	var files []string
	for _, path := range c.CDFiles {
		if strings.ContainsAny(path, "*?[") {
			var globbedFiles []string
			globbedFiles, err = filepath.Glob(path)
			if len(globbedFiles) > 0 {
				files = append(files, globbedFiles...)
			}
		} else {
			_, err = os.Stat(path)
			files = append(files, path)
		}
		if err != nil {
			errs = append(errs, fmt.Errorf("Bad CD disk file '%s': %s", path, err))
		}
		c.CDFiles = files
	}

	return errs
}
