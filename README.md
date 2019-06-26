Vultr builder plugin for Packer
===============================

A Packer builder for creating Vultr snapshots.

## Building

```sh
$ go get github.com/ivandeex/packer-builder-vultr
```

## Configuration

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
