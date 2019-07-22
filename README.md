[![GitHub latest release](https://img.shields.io/github/release/ivandeex/packer-builder-vultr.svg)](https://github.com/ivandeex/packer-builder-vultr/releases)
[![GitHub downloads](https://img.shields.io/github/downloads/ivandeex/packer-builder-vultr/total.svg)](https://github.com/ivandeex/packer-builder-vultr/releases)
[![Build Status](https://travis-ci.org/ivandeex/packer-builder-vultr.svg?branch=master)](https://travis-ci.org/ivandeex/packer-builder-vultr)
[![Go Report Card](https://goreportcard.com/badge/github.com/ivandeex/packer-builder-vultr)](https://goreportcard.com/report/github.com/ivandeex/packer-builder-vultr)
[![GoDoc](https://godoc.org/github.com/ivandeex/packer-builder-vultr?status.svg)](https://godoc.org/github.com/ivandeex/packer-builder-vultr/vultr)

# Packer Builder for Vultr

This is a [HashiCorp Packer](https://www.packer.io/) plugin for creating [Vultr](https://www.vultr.com/) snapshots. It uses the public [Vultr API](https://www.vultr.com/api/).

## Installation
* Download binaries from the [releases page](https://github.com/ivandeex/packer-builder-vultr/releases).
* [Install](https://www.packer.io/docs/extending/plugins.html#installing-plugins) the plugin, or simply put it into the same directory with JSON templates.

## Build

Install Go, then run:
```sh
$ go get github.com/ivandeex/packer-builder-vultr
```

## Examples

See Ubuntu template in the [examples folder](https://github.com/ivandeex/packer-builder-vultr/tree/master/examples/).

You must have your Vultr API token obtained from the Vultr account pages.

```json
{
  "variables": {
      "vultr_api_key": "{{ env `VULTR_API_KEY` }}"
  },
  "builders": [{
      "type": "vultr",
      "api_key": "{{ user `vultr_api_key` }}",
      "description": "My Server",
      "region_name": "Amsterdam",
      "plan_name": "1024 MB RAM,25 GB SSD,1.00 TB BW",
      "os_name": "Ubuntu 18.04 x64",
      "ssh_username": "root",
      "shutdown_command": "/sbin/shutdown -P now"
  }],
  "provisioners": [{
      "type": "file",
      "source": "message.txt",
      "destination": "/root/message.txt"
  }]
}
```

## Parameter Reference

### Build Parameters
* `type`(string) — Always `vultr` for this builder, required.
* `api_key`(string) — Your Vultr API token, required.
* `description`(string) — A name of the resulting snapshot (required).
* `state_timeout`(string) — The time to wait, as a duration string (like `5m`, `60s` or `2h15m45s`), for an instance to power-on or for Vultr to create a final snapshot) before timing out. The default state timeout is "10m" (10 minutes).
* `shutdown_command`(string) — The build will shutdown the instance after all template provisioning is complete but before commanding Vultr to take snapshot. This parameter allows to customize the shutdown command (defaults to `shutdown -P now`).

### Connection

* `ssh_port`(number) — This is only required if you start from a prebaked snapshot. Defaults to `22`.
* `ssh_username`(string) — Optional user name to login into instance after power-on, the default is `root`. If you start building from a standard OS image, then the only option is `root` (usually). You may need to customize user name if you start from a prebaked snapshot.
* `ssh_password`(string) — Not required unless you start from a prebaked snapshot.
* `ssh_private_key_file`(string) — Path to a PEM encoded private key file to use to authenticate with SSH. The `~` can be used in path and will be expanded to the home directory of current user.
* `ssh_timeout`(string) — maximum time to wait for SSH availability after Vultr has created the instance (optional string like `5m` or `60s`, see [packer documentation](http://packer.io/docs/templates/communicator.html#ssh_timeout)

You can skip both password and private key file if you start from a standard OS image, then the builder will use a random root password provided by Vultr to login into new instance.

You can define `ssh_password` and `ssh_private_key_file` together. Then SSH communicator will try them in sequence.

Also, you can use other SSH [connection parameters](http://packer.io/docs/templates/communicator.html#ssh) supported by packer.

### Instance Location (required)

* `region_name`(string) — The name of the region to launch the instance in, e.g. `Frankfurt` or `New Jersey`. Consequently, this is the region where the snapshot will be available.
* `region_code`(string) — You can use a short code instead of a full region name. 
* `region_id`(number) — You can use a numeric region ID instead of a full region name.

Only one of `region_name`, `region_code` or `region_id` is required.

Run the following command to obtain all available regions with their names, short codes and numeric IDs:
```sh
$ curl https://api.vultr.com/v1/regions/list | jq
```

### Instance Parameters
* `plan_name`(string) — The name of the base OS image to use. This is the image that will be used to launch a new instance and provision it.
* `plan_id`(number) — You can use a numeric plan ID instead of a full plan name.
* `ipv6`(boolean, optional) — Set to `true` to enable ipv6 for the instance being created. This defaults to `false`, or not enabled.
* `private_networking`(boolean, optional) — Set to `true` to enable private networking for the instance being created. This defaults to `false`, or not enabled.

Either `plan_name` or `plan_id` is required.

Run the following command to obtain all available plans with their names and numeric IDs:
```sh
$ curl "https://api.vultr.com/v1/plans/list?type=vc2" | jq
```

### Image (required)
You can use either standard OS image from the Vultr library or a pre-saved instance snapshot.

#### Standard OS Image
* `os_name`(string) — The name of the base OS image to use. This is the image that will be used to launch a new instance and provision it.
* `os_id`(number) — You can use a numeric OS image ID instead of a full plan name. Either `os_name` or `os_id` is required.

Run the following command to obtain all available OSes with their names and numeric IDs:
```sh
$ curl https://api.vultr.com/v1/os/list | jq
```

#### Pre-baked Snapshot
* `snapshot_name`(string) — The name of a pre-saved instance snapshot. See the list of snapshots on your Vultr dashboard. The builder will fail if there are multiple snapshots with the same name.
* `snapshot_id`(string) — You can use internal snapshot ID instead of the full name, if you know it. Either `snapshot_name` or `snapshot_id` is required.

### Instance Customizations (optional)
* `startup_script_name`(string) — The name of a startup script to run on the instance before running your template provisioners. Startup scripts can be a powerful tool for configuring the instance from which the image is made. The script must be manually configured in your Vultr dashboard before use. The builder will fail if there are multiple startup scripts with the same name.
* `startup_script_id`(number) — You can use a numeric startup script ID instead of the full name.
* `sshkey_name`(string) — Name of the SSH key to add to the instance on launch. The key must be manually configured in your Vultr dashboard before use. The builder will fail if there are multiple SSH keys with the same name. The `sshkey` parameter name deliberately lacks underscore to make it different from connection related parameters.
* `sshkey_id`(string) — You can use a short SSH key ID instead of the full name.

## Authors
* [**Ivan Andreev**](https://github.com/ivandeex)
