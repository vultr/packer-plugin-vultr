# Change Log

## [v1.0.2](https://github.com/vultr/packer-builder-vultr/compare/v1.0.1..v1.0.2) (2019-10-17)
### Enhancements
- Update govultr + packer to latest releases [#18](https://github.com/vultr/packer-builder-vultr/pull/18)
- Updating Travis supported go versions [#17](https://github.com/vultr/packer-builder-vultr/pull/17)

## [v1.0.1](https://github.com/vultr/packer-builder-vultr/compare/v1.0.0..v1.0.1) (2019-09-11)
### Bug
- Manually shutdown instead of Halt API in shutdown step [#15](https://github.com/vultr/packer-builder-vultr/pull/15)

## [v1.0.0](https://github.com/vultr/packer-builder-vultr/compare/v0.2.1..v1.0.0) (2019-08-12)
### Enhancements
- The packer Vultr plugin has been refactored.
- Unit and Acceptance tests included 
- New documentation page [here](https://github.com/vultr/packer-builder-vultr/blob/master/website/source/docs/builders/vultr.html.md)
- Added the following new config options
  - `instance_label`
  - `userdata`
  - `hostname`
  - `tag`
 
### Breaking changes
- `region_code` renamed to `region_id`
- `region_id` is now an `int`
- `startup_script_id` renamed to `script_id`
- `script_id` is now an `string`
- `sshkey_id` renamed to `ssh_keys_id`
- `sshkey_id` is now an `[]string`
- `ipv6` renamed to `enable_ipv6`
- `private_networking` renamed to `enable_private_network`
- `description` renamed to `snapshot_description`
- removal of: 
  - `region_name`
  - `plan_name`
  - `os_name`
  - `snapshot_name`
  
## [v0.2.1](https://github.com/vultr/packer-builder-vultr/releases/tag/v0.2.1) (2019-06-30)
### Features
- Initial release