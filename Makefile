BINARY = packer-plugin-vultr
GOBIN = $(shell go env GOPATH)/bin
GOLANGCI_LINT_VERSION = v2.15.0
HASHICORP_PACKER_PLUGIN_SDK_VERSION?=$(shell go list -m github.com/hashicorp/packer-plugin-sdk | cut -d " " -f2)

.PHONY: default
default: build test install

.PHONY: build
build:
	go install

.PHONY: install
install: build-binary
	packer plugins install --path ./$(BINARY) github.com/vultr/vultr

.PHONY: test
test:
	go test -v ./... -timeout 60m

.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run ./...

.PHONY: clean
clean:
	rm -rf dist $(BINARY) docs.zip

.PHONY: build-binary
build-binary:
	@go build -o ${BINARY}

install-packer-sdc: ## Install packer sofware development command
	@go install github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc@${HASHICORP_PACKER_PLUGIN_SDK_VERSION}

ci-release-docs: install-packer-sdc
	@$(GOBIN)/packer-sdc renderdocs -src docs -partials docs-partials/ -dst docs/
	@/bin/sh -c "[ -d docs ] && zip -r docs.zip docs/"

plugin-check: install-packer-sdc build-binary
	@$(GOBIN)/packer-sdc plugin-check ${BINARY}
