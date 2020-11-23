BINARY = packer-builder-vultr
PLUGIN_DIR = ~/.packer.d/plugins
GOBIN = $(shell go env GOPATH)/bin
GORELEASER_VER = 0.147.2
GOLANGCI_LINT_VER = 1.17.1

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
	go test -v . ./vultr

.PHONY: lint
lint:
	golint . ./vultr
	golangci-lint run --skip-dirs=test,vendor --fast ./...

.PHONY: clean
clean:
	rm -rf dist $(BINARY)

.PHONY: setup-tools
setup-tools:
	# we want that `go get` install utilities, but in the module mode its
	# behaviour is different; actually, `go get` would rather modify the
	# local `go.mod`, so let's disable modules here.
	GO111MODULE=off go get -u golang.org/x/lint/golint
	GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports
	# goreleaser and golangci-lint take pretty long to build
	# as an optimization, let's just download the binaries
	curl -sL "https://github.com/goreleaser/goreleaser/releases/download/v$(GORELEASER_VER)/goreleaser_Linux_x86_64.tar.gz" | tar -xzf - -C $(GOBIN) goreleaser
	curl -sL "https://github.com/golangci/golangci-lint/releases/download/v$(GOLANGCI_LINT_VER)/golangci-lint-$(GOLANGCI_LINT_VER)-linux-amd64.tar.gz" | tar -xzf - -C $(GOBIN) --strip-components=1 "golangci-lint-$(GOLANGCI_LINT_VER)-linux-amd64/golangci-lint"
