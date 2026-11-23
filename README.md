# Packer Builder for Vultr

[![Automatic Releaser](https://github.com/vultr/packer-plugin-vultr/actions/workflows/release.yml/badge.svg?branch=master)](https://github.com/vultr/packer-plugin-vultr/actions/workflows/release.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/vultr/packer-plugin-vultr.svg)](https://pkg.go.dev/github.com/vultr/packer-plugin-vultr)
[![GitHub latest release](https://img.shields.io/github/release/vultr/packer-plugin-vultr.svg)](https://github.com/vultr/packer-plugin-vultr/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/vultr/packer-plugin-vultr)](https://goreportcard.com/report/github.com/vultr/packer-plugin-vultr)
[![GitHub downloads](https://img.shields.io/github/downloads/vultr/packer-plugin-vultr/total.svg)](https://github.com/vultr/packer-plugin-vultr/releases)


This is a [HashiCorp Packer](https://www.packer.io/) plugin for creating [Vultr](https://www.vultr.com/) snapshots.

## Requirements
* [Packer](https://www.packer.io/intro/getting-started/install.html)
* [Go 1.27.4+](https://go.dev/doc/install) (building from source only)

## Build & Installation

### Packer init
Starting from version 1.7, Packer supports a new packer init command allowing automatic installation of Packer plugins. Read the [Packer documentation](https://www.packer.io/docs/commands/init) for more information

```hcl
packer {
  required_plugins {
    vultr = {
      source  = "github.com/vultr/vultr"
      version = "~> 3.0"
    }
  }
}
```

Run `packer init .` in the directory containing your template.

### Upgrading to v3

* VPC 2.0 is no longer supported by GoVultr. Remove `enable_vpc2 = true` and use
  `enable_vpc` with a supported VPC configuration instead.
* The Packer SDK no longer accepts DSA SSH keys. Use Ed25519, ECDSA, or RSA keys.

### Install from source:

To install the plugin from source you will have to clone this Github repository locally. Once cloned run `make build` which will put the plugin in your GO's bin folder.
Then please read the documentation on [Installing plugins](https://www.packer.io/docs/plugins#installing-plugins).

### Install from release:

* Download binaries from the [releases page](https://github.com/vultr/packer-plugin-vultr/releases).
* To install the plugin please read the documentation on [Installing plugins](https://www.packer.io/docs/plugins#installing-plugins)

## Using the plugin
See the Vultr Provider [documentation](docs/builders/vultr.mdx) to get started using the Vultr provider.

## Contributing
Feel free to send pull requests our way! Please see the [contributing](CONTRIBUTING.md) guidelines.

## Authors
* [**Ivan Andreev**](https://github.com/ivandeex)
* [**Fady Farid**](https://github.com/afady)
* [**David Dymko**](https://github.com/ddymko)
