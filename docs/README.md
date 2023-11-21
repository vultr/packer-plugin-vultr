The [Vultr](https://www.Vultr.com/) Packer plugin provides a builder for building images in
Vultr.

### Installation

To install this plugin, copy and paste this code into your Packer configuration, then run [`packer init`](https://www.packer.io/docs/commands/init).

```hcl
packer {
  required_plugins {
    vultr = {
      version = ">= 2.5.0"
      source  = "github.com/vultr/vultr"
    }
  }
}
```

Alternatively, you can use `packer plugins install` to manage installation of this plugin.

```sh
$ packer plugins install github.com/vultr/vultr
```


### Components

#### Builders

- [vultr](/packer/integrations/vultr/latest/components/builder/vultr) - The builder takes a source image, runs any provisioning necessary on the image after launching it, then snapshots it into a reusable image. This reusable image can then be used as the foundation of new servers that are launched within Vultr.
