# Vultr Plugins

The [Vultr](https://www.Vultr.com/) Packer plugin provides a builder for building images in
Vultr.

## Installation

### Using pre-built releases

#### Using the `packer init` command

Starting from version 1.7, Packer supports a new `packer init` command allowing
automatic installation of Packer plugins. Read the
[Packer documentation](https://www.packer.io/docs/commands/init) for more information.

To install this plugin, copy and paste this code into your Packer configuration .
Then, run [`packer init`](https://www.packer.io/docs/commands/init).

```hcl
packer {
  required_plugins {
    vultr = {
      version = ">= 2.6.3"
      source  = "github.com/vultr/vultr"
    }
  }
}
```

#### Manual installation

You can find pre-built binary releases of the plugin [here](https://github.com/vultr/packer-plugin-vultr/releases).
Once you have downloaded the latest archive corresponding to your target OS,
uncompress it to retrieve the plugin binary file corresponding to your platform.
To install the plugin, please follow the Packer documentation on
[installing a plugin](https://www.packer.io/docs/extending/plugins/#installing-plugins).


#### From Source

If you prefer to build the plugin from its source code, clone the GitHub
repository locally and run the command `go build` from the root
directory. Upon successful compilation, a `packer-plugin-vultr` plugin
binary file can be found in the root directory.
To install the compiled plugin, please follow the official Packer documentation
on [installing a plugin](https://www.packer.io/docs/extending/plugins/#installing-plugins).


## Plugin Contents

The Vultr plugin is intended as a starting point for creating Packer plugins, containing:

### Builders

- [builder](/docs/builders/vultr.mdx) - The builder takes a source image, runs any provisioning necessary on the image after launching it, then snapshots it into a reusable image. This reusable image can then be used as the foundation of new servers that are launched within Vultr.
