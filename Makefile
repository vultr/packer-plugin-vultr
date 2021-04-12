BINARY = packer-plugin-vultr
PLUGIN_DIR = ~/.packer.d/plugins
GOBIN = $(shell go env GOPATH)/bin

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
	go test -v . ./builder/vultr

.PHONY: lint
lint:
	golint . ./builder/vultr
	golangci-lint run --skip-dirs=test,vendor --fast ./...

.PHONY: clean
clean:
	rm -rf dist $(BINARY)
