SHELL := /bin/sh

# application name and build output directory
APP := edsrv
OUTPUT_DIR ?= bin

# application directories
CMD ?= ./cmd/$(APP)/...
PKG ?= ./pkg/$(APP)/...

# golang build and test configuration
GOARCH ?= $(shell go env GOARCH)
GOHOSTOS ?= $(shell go env GOHOSTOS)
GOFLAGS ?= -race -a
CGO_LDFLAGS ?= -s -w
GOFLAGS_TEST ?= -failfast -race -cover -v

# integration testing directory and default ginkgo flags
TEST_INTEGRATION ?= ./test/integration/...
GINKGO_FLAGS ?= -vv --race --fail-fast

# executables full path and intallation prefix
BIN ?= $(OUTPUT_DIR)/$(APP)
PREFIX ?= /usr/local/bin

# macos plist to define a edit-server service, running in the background
PLIST ?= contrib/$(APP).plist
LAUNCHAGENT_DIR ?= ~/Library/LaunchAgents
LAUNCHAGENT_LABEL ?= io.github.otaviof.edsrv
LAUNCHAGENT_PLIST ?= $(LAUNCHAGENT_DIR)/$(LAUNCHAGENT_LABEL).plist

# general arguments for "run" target
ARGS ?=

.EXPORT_ALL_VARIABLES:

default: build

# Builds the application binary.
.PHONY: $(BIN)
$(BIN):
	go build -o $(BIN) $(CMD)

build: $(BIN)

# Runs the application using ARGS to inform extra parameter and flags.
.PHONY: run
run: GOFLAGS = 
run:
	go run $(CMD) $(ARGS)

# Cleans up the output build directory completely.
.PHONY: clean
clean:
	test -d "$(OUTPUT_DIR)" && rm -rf "$(OUTPUT_DIR)" >/dev/null

# Runs the unitary tests.
.PHONY: test-unit
test-unit:
	go test $(GOFLAGS_TEST) $(CMD) $(PKG)

# Runs the integration tests.
.PHONY: test-integration
test-integration:
	ginkgo run $(GINKGO_FLAGS) $(TEST_INTEGRATION)

# Runs all the tests available.
test: test-unit test-integration

# Installs Ginkgo command-line matching the "go.mod" version.
install-ginkgo:
	go install -v github.com/onsi/ginkgo/v2/ginkgo

# Installs the tools needed for releasing, linting and testing, if they are not
# installed already.
.PHONY: install-tools
install-tools: install-ginkgo
	which goreleaser >/dev/null 2>&1 || \
		go install github.com/goreleaser/goreleaser@latest
	which golangci-lint >/dev/null 2>&1 || \
		go install -v github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Uses golangci-lint to inspect the code base.
.PHONY: lint
lint:
	golangci-lint run ./...

# Uses goreleaser to create a snapshot build.
.PHONY: snapshot
snapshot:
	goreleaser build \
		--clean \
		--snapshot \
		--single-target \
		--output=$(BIN) >/dev/null 2>&1

# Installs the application using the binary built by goreleaser, which tends to be
# sightly smaller than a regular build, the installation happens on ${PREFIX}
# directory with the same application name.
.PHONY: install
install: snapshot
	sudo install -o root -g wheel -m 0755 $(BIN) $(PREFIX)

# Installs the plist service definition for macOS.
.PHONY: install-launchagent
install-launchagent: install
	install -g staff -m 0755 $(PLIST) $(LAUNCHAGENT_PLIST)

.PHONY: sleep
sleep:
	@sleep 1

# Loads the launch-agent service file.
.PHONY: load-macos
load-macos:
	launchctl load -w $(LAUNCHAGENT_PLIST)

# Shows the macOS service status.
.PHONY: status-macos
status-macos:
	launchctl list $(LAUNCHAGENT_LABEL)
	$(APP) status

# Installs, launch and list the status of the macOS service.
.PHONY: install-macos
install-macos: \
	install-launchagent \
	load-macos \
	sleep \
	status-macos

# Unloads and removes the macOS launch-agent service.
.PHONY: uninstall-macos 
uninstall-macos:
	launchctl unload -w $(LAUNCHAGENT_PLIST)
	rm -f -v $(LAUNCHAGENT_PLIST) 

# Stops the macOS service.
.PHONY: stop-macos
stop-macos:
	launchctl stop $(LAUNCHAGENT_LABEL)

# Starts the macOS service.
.PHONY: start-macos
start-macos:
	launchctl start $(LAUNCHAGENT_LABEL)

# Restarts the macos service.
restart-macos: \
	stop-macos \
	sleep \
	start-macos \
	sleep \
	status-macos
