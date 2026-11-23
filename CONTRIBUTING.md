# Contributing to `packer-plugin-vultr`

Feedback, bug reports, and pull requests are welcome.

## Getting started

Clone this repository and install the Go version specified in `go.mod` (or newer).

```sh
make build
make test
make lint
```

Format Go changes with `gofmt`. After changing the configuration schema, run
`make install-packer-sdc` and `PATH="$(go env GOPATH)/bin:$PATH" go generate ./...`.
Run `make plugin-check` to check Packer plugin compatibility. To install a local
build, install Packer and run `make install`.

Acceptance tests require `PACKER_ACC=1` and `VULTR_API_KEY` and create billable
Vultr resources.

## Versioning

This project follows [SemVer](https://semver.org/): breaking changes increment
the major version, features the minor version, and fixes the patch version.

## Releases

Maintainers prepare release changes in a pull request:

1. Update `Version` in `main.go`, `CHANGELOG.md`, and the documented plugin version.
2. Run tests, lint, plugin checks, and a GoReleaser snapshot build.
3. After the changes pass review and CI, push a release commit to `master` with
   the exact subject `Release vX.Y.Z #major`, `#minor`, or `#patch` as appropriate.
4. The release workflow creates that exact tag, builds the platform archives and
   documentation, signs the checksums, and creates a draft GitHub release.
5. Verify the artifacts and version metadata, add release notes, and publish the draft.
