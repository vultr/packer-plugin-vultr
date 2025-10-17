# For full specification on the configuration of this file visit:
# https://github.com/hashicorp/integration-template#metadata-configuration
integration {
  name = "Vultr"
  description = "A plugin for creating Vultr snapshots."
  identifier = "packer/vultr/vultr"
  component {
    type = "builder"
    name = "Vultr"
    slug = "vultr"
  }
}
