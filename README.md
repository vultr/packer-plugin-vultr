# Packer Builder for Vultr

[![Build Status](https://travis-ci.org/vultr/packer-builder-vultr.svg?branch=master)](https://travis-ci.org/vultr/packer-builder-vultr)
[![GoDoc](https://godoc.org/github.com/vultr/packer-builder-vultr?status.svg)](https://godoc.org/github.com/vultr/packer-builder-vultr/vultr)
[![GitHub latest release](https://img.shields.io/github/release/vultr/packer-builder-vultr.svg)](https://github.com/vultr/packer-builder-vultr/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/vultr/packer-builder-vultr)](https://goreportcard.com/report/github.com/vultr/packer-builder-vultr)
[![GitHub downloads](https://img.shields.io/github/downloads/vultr/packer-builder-vultr/total.svg)](https://github.com/vultr/packer-builder-vultr/releases)


This is a [HashiCorp Packer](https://www.packer.io/) plugin for creating [Vultr](https://www.vultr.com/) snapshots.

## Requirements
* [Packer](https://www.packer.io/intro/getting-started/install.html)
* [Go 1.16+](https://golang.org/doc/install)

## Build & Installation

### Packer init
Starting from version 1.7, Packer supports a new packer init command allowing automatic installation of Packer plugins. Read the [Packer documentation](https://www.packer.io/docs/commands/init) for more information

### Install from source:

To install the plugin from source you will have to clone this Github repository locally. Once cloned run `make build` which will put the plugin in your GO's bin folder.
Then please read the documentation on [Installing plugins](https://www.packer.io/docs/plugins#installing-plugins).

### Install from release:

* Download binaries from the [releases page](https://github.com/vultr/packer-builder-vultr/releases).
* To install the plugin please read the documentation on [Installing plugins](https://www.packer.io/docs/plugins#installing-plugins)

## Using the plugin
See the Vultr Provider [documentation](docs/builders/vultr.mdx) to get started using the Vultr provider.

## Contributing
Feel free to send pull requests our way! Please see the [contributing](CONTRIBUTING.md) guidelines.

## Authors
* [**Ivan Andreev**](https://github.com/ivandeex)
* [**Fady Farid**](https://github.com/afady)
* [**David Dymko**](https://github.com/ddymko)
