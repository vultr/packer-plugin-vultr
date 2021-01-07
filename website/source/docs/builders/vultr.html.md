---
description: |
    The vultr Packer builder is able to create new images for use with
    Vultr. The builder takes a source image, runs any provisioning necessary
    on the image after launching it, then snapshots it into a reusable image. This
    reusable image can be then used as the foundation of new servers that are
    launched within Vultr.
layout: docs
page_title: 'Vultr - Builders'
sidebar_current: 'docs-builders-vultr'
---

# Vultr Builder

Type: `vultr`

The `vultr` Packer builder is able to create new images for use with
[Vultr](https://www.vultr.com). The builder takes a source image,
runs any provisioning necessary on the image after launching it, then snapshots
it into a reusable image. This reusable image can be then used as the
foundation of new servers that are launched within Vultr.

The builder does *not* manage images. Once it creates an image, it is up to you
to use it or delete it.

**NOTE**: Packer-Builder-Vultr [v2+](https://github.com/vultr/packer-builder-vultr/blob/master/CHANGELOG.md#v200-2020-11-23) uses [API V2](https://www.vultr.com/api/v2)

## Configuration Reference

There are many configuration options available for the builder. They are
segmented below into two categories: required and optional parameters. Within
each category, the available configuration keys are alphabetized.

In addition to the options listed here, a
[communicator](https://www.packer.io/docs/communicators) can be configured for this
builder.

### Required:

-   `api_key` (string) - The Vultr API Key to access your account.

-   `os_id` (int) - The id of the os to use. This will be the OS that will be used to launch a new instance and provision it. See [List Operating Systems](https://www.vultr.com/api/v2/#operation/list-os).

-   `region_id` (string) - The id of the region to launch the instance in. See [List Regions](https://www.vultr.com/api/v2/#operation/list-regions).
    
-   `plan_id` (string) - The id of the plan you wish to use. See [List Plans](https://www.vultr.com/api/v2/#tag/plans).

### Optional:

-   `snapshot_description` (string) - Description of the snapshot.

-   `snapshot_id` (string) -   If you've selected the 'snapshot' (OS 164) operating system, this should be the ID of the snapshot. See [Snapshot](https://www.vultr.com/api/v2/#operation/list-snapshots).

-   `iso_id` (string) - If you've selected the 'custom' (OS 159) operating system, this is the ID of a specific ISO to mount during the deployment. See [ISO](https://www.vultr.com/api/v2/#operation/list-isos).

-   `app_id` (int) - If launching an application (OSID 186), this is the ID to launch. See [App](https://www.vultr.com/api/v2/#operation/list-applications).

-   `enable_ipv6` (bool) - IPv6 subnet will be assigned to the machine.

-   `enable_private_network` (bool) - Enables private networking support to the new server.

-   `script_id` (string) - If you've not selected a 'custom' (OS 159) operating system, this can be the `id` of a startup script to execute on boot. See [Startup Script](https://www.vultr.com/api/v2/#operation/list-startup-scripts).

-   `ssh_key_ids` (array of string) - List of SSH keys to apply to this server on install. Separate keys with commas. See [SSH Key](https://www.vultr.com/api/v2/#operation/list-ssh-keys).

-   `instance_label` (string) - This is a text label that will be shown in the control panel.

-   `userdata` (string) - Base64 encoded user-data.

-   `hostname` (string) - Hostname to assign to this server.

-   `tag` (string) - The tag to assign to this server.

-   `state_timeout` (string) - A duration to wait for the instance to boot, or a snapshot to be taken. Must be a string in [golang Duration-parsable format](https://golang.org/pkg/time/#ParseDuration), like "10m" or "30s". 

## Basic Example

Here is a Vultr builder example. The vultr_api_key should be replaced with an actual Vultr API Key

``` json
{
 "variables": {
        "vultr_api_key": "{{ env `VULTR_API_KEY` }}"
    },
    "builders": [{
        "type": "vultr",
        "api_key": "{{ user `vultr_api_key` }}",
        "snapshot_description": "Packer-test-with updates",
        "region_id": "ewr",
        "plan_id": "vc2-1c-1gb",
        "os_id": 127,
        "state_timeout": "15m"
    }]
}
```
