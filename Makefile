BINARY = packer-plugin-vultr
PLUGIN_DIR = ~/.packer.d/plugins
GOBIN = $(shell go env GOPATH)/bin
HASHICORP_PACKER_PLUGIN_SDK_VERSION?=$(shell go list -m github.com/hashicorp/packer-plugin-sdk | cut -d " " -f2)

.PHONY: default
default: build test install

.PHONY: build
build:
	go install

.PHONY: install
install: build
	mkdir -p $(PLUGIN_DIR)
	cp -f $(GOBIN)/$(BINARY) $(PLUGIN_DIR)/$(BINARY)

.PHONY: test
test:
	go test -v . ./builder/vultr -timeout 60m

.PHONY: lint
lint:
	golint . ./builder/vultr
	golangci-lint run --skip-dirs=test,vendor --fast ./...

.PHONY: clean
clean:
	rm -rf dist $(BINARY)

build-binary:
	go install
	@go build -o ${BINARY}
	mkdir -p $(PLUGIN_DIR)
	cp -f $(GOBIN)/$(BINARY) $(PLUGIN_DIR)/$(BINARY)

install-packer-sdc: ## Install packer sofware development command
	@go install github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc@${HASHICORP_PACKER_PLUGIN_SDK_VERSION}

ci-release-docs: install-packer-sdc
	@packer-sdc renderdocs -src docs -partials docs-partials/ -dst docs/
	@/bin/sh -c "[ -d docs ] && zip -r docs.zip docs/"

plugin-check: install-packer-sdc build-binary
	@packer-sdc plugin-check ${BINARY}

build-docs: install-packer-sdc
	@if [ -d ".docs" ]; then rm -r ".docs"; fi
	@packer-sdc renderdocs -src "docs" -partials docs-partials/ -dst ".docs/"
	@./.web-docs/scripts/compile-to-webdocs.sh "." ".docs" ".web-docs" "vultr"
	@rm -r ".docs"
