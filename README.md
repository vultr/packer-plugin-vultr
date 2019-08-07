# Packer Builder for Vultr

[![Build Status](https://travis-ci.org/vultr/packer-builder-vultr.svg?branch=master)](https://travis-ci.org/vultr/packer-builder-vultr)
[![GoDoc](https://godoc.org/github.com/vultr/packer-builder-vultr?status.svg)](https://godoc.org/github.com/vultr/packer-builder-vultr/vultr)
[![GitHub latest release](https://img.shields.io/github/release/vultr/packer-builder-vultr.svg)](https://github.com/vultr/packer-builder-vultr/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/vultr/packer-builder-vultr)](https://goreportcard.com/report/github.com/vultr/packer-builder-vultr)
[![GitHub downloads](https://img.shields.io/github/downloads/vultr/packer-builder-vultr/total.svg)](https://github.com/vultr/packer-builder-vultr/releases)


This is a [HashiCorp Packer](https://www.packer.io/) plugin for creating [Vultr](https://www.vultr.com/) snapshots.

## Requirements
* [Packer](https://www.packer.io/intro/getting-started/install.html)
* [Go 1.12+](https://golang.org/doc/install)

## Build & Installation

### Install from source:

Clone repository to `$GOPATH/src/github.com/vultr/packer-builder-vultr`

```sh
$ mkdir -p $GOPATH/src/github.com/vultr; cd $GOPATH/src/github.com/vultr
$ git clone git@github.com:vultr/packer-builder-vultr.git
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/vultr/packer-builder-vultr
$ make build
```

Link the build to Packer

```sh
$ln -s $GOPATH/bin/packer-builder-vultr ~/.packer.d/plugins/packer-builder-vultr 
```

### Install from release:

* Download binaries from the [releases page](https://github.com/vultr/packer-builder-vultr/releases).
* [Install](https://www.packer.io/docs/extending/plugins.html#installing-plugins) the plugin, or simply put it into the same directory with JSON templates.
* Move the downloaded binary to `~/.packer.d/plugins/`

## Using the plugin
See the Vultr Provider [documentation](website/source/docs/builders/vultr.html.md) to get started using the Vultr provider.

## Contributing
Feel free to send pull requests our way! Please see the [contributing](CONTRIBUTING.md) guidelines.

## Authors
* [**Ivan Andreev**](https://github.com/ivandeex)
* [**Fady Farid**](https://github.com/afady)
* [**David Dymko**](https://github.com/ddymko)
