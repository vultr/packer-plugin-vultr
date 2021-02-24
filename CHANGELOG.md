# Change Log

## [v2.2.1](https://github.com/vultr/packer-builder-vultr/compare/v2.2.0..v2.2.1) (2021-02-24)
### Bug fixes
- Updated HCL config to support API v2 type attributes [#116](https://github.com/vultr/packer-builder-vultr/pull/116)

## [v2.2.0](https://github.com/vultr/packer-builder-vultr/compare/v2.1.0..v2.2.0) (2021-02-17)
### Enhancements
- Updated plugin to support Packer v1.7.0. This release does not support `packer init` this will be released in `v2.3.0` [#110](https://github.com/vultr/packer-builder-vultr/pull/110)

## [v2.1.0](https://github.com/vultr/packer-builder-vultr/compare/v2.0.1..v2.1.0) (2021-02-16)
### Bug
- Remove the requirement to pass in `os_id` with other installation types. This is no longer required in API v2. [#105](https://github.com/vultr/packer-builder-vultr/pull/105)

### Dependencies
- GoVultr v2.3.2 -> v2.4.0 [#108](https://github.com/vultr/packer-builder-vultr/pull/108)

## [v2.0.1](https://github.com/vultr/packer-builder-vultr/compare/v2.0.0..v2.0.1) (2021-01-07)
### Dependencies
- hcl 2.7.1 -> 2.8.1 [#98](https://github.com/vultr/packer-builder-vultr/pull/98)
- zclconf/go-cty 1.7.0 -> 1.7.1 [#96](https://github.com/vultr/packer-builder-vultr/pull/96)
- GoVultr v2.0.0 -> v2.3.2 [#98](https://github.com/vultr/packer-builder-vultr/pull/98)


## [v2.0.0](https://github.com/vultr/packer-builder-vultr/compare/v1.0.15..v2.0.0) (2020-11-23)
### Changes
GoVultr was updated to v2.0.0 which is running on API v2.

### Dependencies 
- hcl 2.7.0 -> 2.7.1 [#82](https://github.com/vultr/packer-builder-vultr/pull/82)
- packer 1.6.4 -> 1.6.5 [#82](https://github.com/vultr/packer-builder-vultr/pull/82)
- zclconf/go-cty 1.6.1 -> 1.7.0 [#82](https://github.com/vultr/packer-builder-vultr/pull/82)
- GoVultr v0.5.0 -> v2.0.0 [#81](https://github.com/vultr/packer-builder-vultr/pull/81)

## [v1.0.15](https://github.com/vultr/packer-builder-vultr/compare/v1.0.14..v1.0.15) (2020-11-11)
### Bug Fix
- missing defer on f.close in ssh debug [#78](https://github.com/vultr/packer-builder-vultr/pull/78)

## [v1.0.14](https://github.com/vultr/packer-builder-vultr/compare/v1.0.13..v1.0.14) (2020-10-19)
### Dependencies
- hcl 2.6.0 -> 2.7.0 [#72](https://github.com/vultr/packer-builder-vultr/pull/72)

## [v1.0.13](https://github.com/vultr/packer-builder-vultr/compare/v1.0.12..v1.0.13) (2020-10-06)
### Dependencies
- packer 1.6.2 -> 1.6.4 [#70](https://github.com/vultr/packer-builder-vultr/pull/70)

## [v1.0.12](https://github.com/vultr/packer-builder-vultr/compare/v1.0.11..v1.0.12) (2020-09-25)
### Changes
- Vultr-packer-builder now uses the core ssh communicator instead of the one that was written [66](https://github.com/vultr/packer-builder-vultr/pull/66)

### Dependencies
- zclconf/go-cty 1.6.0 -> 1.6.1 [#61](https://github.com/vultr/packer-builder-vultr/pull/61)

## [v1.0.11](https://github.com/vultr/packer-builder-vultr/compare/v1.0.10..v1.0.11) (2020-09-01)
### Dependencies
- packer 1.6.1 -> 1.6.2 [#59](https://github.com/vultr/packer-builder-vultr/pull/59)
- zclconf/go-cty 1.5.1 -> 1.6.0 [#58](https://github.com/vultr/packer-builder-vultr/pull/58)
- GoVultr v0.4.2 -> v0.5.0 [#56](https://github.com/vultr/packer-builder-vultr/pull/56)

## [v1.0.10](https://github.com/vultr/packer-builder-vultr/compare/v1.0.9..v1.0.10) (2020-07-31)
### Dependencies 
- packer 1.5.6 -> 1.6.1 [#54](https://github.com/vultr/packer-builder-vultr/pull/54)
- zclconf/go-cty 1.4.0 -> 1.5.1 [#53](https://github.com/vultr/packer-builder-vultr/pull/53) [#52](https://github.com/vultr/packer-builder-vultr/pull/52)


## [v1.0.9](https://github.com/vultr/packer-builder-vultr/compare/v1.0.8..v1.0.9) (2020-06-09)
### Dependencies 
- hcl/v2 2.5.1 -> 2.6.0 [#48](https://github.com/vultr/packer-builder-vultr/pull/48)
- GoVultr v0.4.1 -> v0.4.2 [#47](https://github.com/vultr/packer-builder-vultr/pull/47)

## [v1.0.8](https://github.com/vultr/packer-builder-vultr/compare/v1.0.7..v1.0.8) (2020-05-19)
### Dependencies 
- hcl/v2 2.5.0 -> 2.5.1 (Fixes panic) [#45](https://github.com/vultr/packer-builder-vultr/pull/45)

## [v1.0.7](https://github.com/vultr/packer-builder-vultr/compare/v1.0.6..v1.0.7) (2020-05-11)
### Dependencies 
- GoVultr v0.3.3 -> v0.4.1 [#43](https://github.com/vultr/packer-builder-vultr/pull/43)
- hcl/v2 2.4.0 -> 2.5.0 [#42](https://github.com/vultr/packer-builder-vultr/pull/42)
- packer 1.5.5 -> 1.5.6 [#41](https://github.com/vultr/packer-builder-vultr/pull/41)

## [v1.0.6](https://github.com/vultr/packer-builder-vultr/compare/v1.0.5..v1.0.6) (2020-04-16)
### Dependencies 
- GoVultr v0.3.2 -> v0.3.3 [#38](https://github.com/vultr/packer-builder-vultr/pull/38)
- hcl/v2 2.3.0 -> 2.4.0 [#37](https://github.com/vultr/packer-builder-vultr/pull/37)

## [v1.0.5](https://github.com/vultr/packer-builder-vultr/compare/v1.0.4..v1.0.5) (2020-03-31)
### Dependencies 
- GoVultr v0.3.0 -> v0.3.2 [#31](https://github.com/vultr/packer-builder-vultr/pull/31) [#34](https://github.com/vultr/packer-builder-vultr/pull/34)
- go-cty 1.2.1 -> 1.3.1 [#30](https://github.com/vultr/packer-builder-vultr/pull/30)
- packer 1.5.4 -> 1.5.5 [33](https://github.com/vultr/packer-builder-vultr/pull/33)

## [v1.0.4](https://github.com/vultr/packer-builder-vultr/compare/v1.0.3..v1.0.4) (2020-03-04)
### Enhancements
- Updated dependencies to newer versions [#27](https://github.com/vultr/packer-builder-vultr/pull/27)

## [v1.0.3](https://github.com/vultr/packer-builder-vultr/compare/v1.0.2..v1.0.3) (2020-02-14)
### Enhancements
- Updated Packer to 1.5.2 and GoVultr to v0.2.0 [#23](https://github.com/vultr/packer-builder-vultr/pull/23)

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
