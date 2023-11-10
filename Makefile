SHELL := /bin/sh

# application name and build output directory
APP := edsrv
OUTPUT_DIR ?= bin

# application directories
CMD ?= ./cmd/$(APP)/...
PKG ?= ./pkg/$(APP)/...

# golang build and test configuration
GOHOSTOS ?= $(shell go env GOHOSTOS)
GOFLAGS ?= -v -a -race
CGO_ENABLED ?= 1
CGO_LDFLAGS ?= -s -w
GOFLAGS_TEST ?= -failfast -race -cover -v

# integration testing directory and default ginkgo flags
GINKGO_FLAGS ?= -vv --race --fail-fast
TEST_INTEGRATION ?= ./test/integration/...
ACT_WORKFLOWS ?= .github/workflows/test.yaml

# executables full path and intallation prefix
BIN ?= $(OUTPUT_DIR)/$(APP)
PREFIX ?= /usr/local/bin

# github action tag (release) version
GITHUB_REF_NAME ?= ${GITHUB_REF_NAME:-}

# macos plist to define a edit-server service, running in the background
PLIST ?= contrib/$(APP).plist
LAUNCHAGENT_DIR ?= ~/Library/LaunchAgents
LAUNCHAGENT_LABEL ?= io.github.otaviof.edsrv
LAUNCHAGENT_PLIST ?= $(LAUNCHAGENT_DIR)/$(LAUNCHAGENT_LABEL).plist

# general arguments for "run" target
ARGS ?=

.EXPORT_ALL_VARIABLES:

#
# Tools
#

# Installs GoReleaser.
tool-goreleaser: GOFLAGS =
tool-goreleaser:
	@which goreleaser >/dev/null 2>&1 || \
		go install github.com/goreleaser/goreleaser@latest >/dev/null 2>&1

# Installs golangcdi-lint.
tool-golangci-lint: GOFLAGS =
tool-golangci-lint:
	@which golangci-lint >/dev/null 2>&1 || \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest >/dev/null 2>&1

# Installs Ginkgo matching the "go.mod" version.
tool-ginkgo: GOFLAGS =
tool-ginkgo:
	@which ginkgo >/dev/null 2>&1 || \
		go install github.com/onsi/ginkgo/v2/ginkgo >/dev/null 2>&1

# Instllas GitHub CLI ("gh").
tool-gh: GOFLAGS =
tool-gh:
	@which gh >/dev/null 2>&1 || \
		go install github.com/cli/cli/v2/cmd/gh@latest >/dev/null 2>&1

#
# Build and Run
#

default: build

# Builds the application binary.
.PHONY: $(BIN)
$(BIN):
	go build -o $(BIN) $(CMD)

build: $(BIN)

# Uses goreleaser to create a snapshot build.
.PHONY: goreleaser-snapshot
goreleaser-snapshot: tool-goreleaser
	goreleaser build --clean --snapshot --single-target -o=$(BIN)

snapshot: goreleaser-snapshot

# Runs the application using ARGS to inform extra parameter and flags.
.PHONY: run
run: GOFLAGS = 
run:
	go run $(CMD) $(ARGS)

# Cleans up the output build directory completely.
.PHONY: clean
clean:
	test -d "$(OUTPUT_DIR)" && rm -rf "$(OUTPUT_DIR)" >/dev/null

#
# Test and Lint
#

# Runs the unitary tests.
.PHONY: test-unit
test-unit:
	go test $(GOFLAGS_TEST) $(CMD) $(PKG)

# Runs the integration tests.
.PHONY: test-integration
test-integration: tool-ginkgo
	ginkgo run $(GINKGO_FLAGS) $(TEST_INTEGRATION)

# Runs all the tests available.
test: test-unit test-integration

# Simulates GitHub Action workflows with "act" (https://github.com/nektos/act)
.PHONY: act
act:
	act --workflows=$(ACT_WORKFLOWS) $(ARGS)

# Uses golangci-lint to inspect the code base.
.PHONY: lint
lint: tool-golangci-lint
	golangci-lint run ./...

#
# Install
#

# Installs the application using the binary built by goreleaser, which tends to be
# sightly smaller than a regular build, the installation happens on ${PREFIX}
# directory with the same application name.
.PHONY: install
install: snapshot
	sudo install -o root -g wheel -m 0755 $(BIN) $(PREFIX)

# Uninstalls the primary application executable.
.PHONY: uninstall
uninstall:
	sudo rm -f -v ${PREFIX}/${APP} || true

#
# Service
#

# Calls the application to check the service status
.PHONY: status
status:
	$(APP) status $(ARGS)

# Installs the plist service definition for macOS.
.PHONY: macos-install-launchagent
macos-install-launchagent: install
	install -g staff -m 0755 $(PLIST) $(LAUNCHAGENT_PLIST)

.PHONY: sleep
sleep:
	@sleep 1

# Loads the launch-agent service file.
.PHONY: macos-launchctl-load
macos-launchctl-load:
	launchctl load -w $(LAUNCHAGENT_PLIST)

# Shows the macOS service status.
.PHONY: macos-launchctl-list
macos-launchctl-list:
	launchctl list $(LAUNCHAGENT_LABEL)

# Deploys the macOS service and starts it.
.PHONY: macos-deploy-service
macos-deploy-service: \
	macos-install-launchagent \
	macos-launchctl-load \
	sleep \
	macos-launchctl-list \
	status

# Unloads the launch-agent plist file.
.PHONY: macos-unload-launchagent
macos-unload-launchagent:
	launchctl unload -w $(LAUNCHAGENT_PLIST)

# Removes the launch-agent plist file.
macos-remove-launchagent:
	rm -f -v $(LAUNCHAGENT_PLIST) || true

# Removes the whole macos service, the oppositive of "deploy" target.
macos-remove-service: \
	macos-unload-launchagent \
	sleep \
	macos-remove-launchagent

# Stops the macOS service.
.PHONY: macos-launchctl-stop
macos-launchctl-stop:
	launchctl stop $(LAUNCHAGENT_LABEL)

# Starts the macOS service.
.PHONY: macos-launchctl-start
macos-launchctl-start:
	launchctl start $(LAUNCHAGENT_LABEL)

# Restarts the macOS service.
mac-restart-service: \
	macos-launchctl-stop \
	sleep \
	macos-launchctl-start \
	sleep \
	macos-launchctl-list \
	status

#
# GitHub Release
#

# Inspects GITHUB_REF_NAME variable, when the application is being released on
# GitHub this variable shows the subject revision tag, which is also used as
# version.
github-ref-name-probe:
ifeq ($(strip $(GITHUB_REF_NAME)),)
	$(error variable GITHUB_REF_NAME is not set)
endif

# Creates a new GitHub release with GITHUB_REF_NAME.
github-release-create: tool-gh
	gh release view $(GITHUB_REF_NAME) >/dev/null 2>&1 || \
		gh release create --generate-notes $(GITHUB_REF_NAME)

# Runs "goreleaser" to build the artefacts and upload them into the current
# release payload, it amends the release in progress with the application
# executables.
goreleaser-release: tool-goreleaser
goreleaser-release: CGO_ENABLED = 0
goreleaser-release: GOFLAGS = -a
goreleaser-release:
	goreleaser release --clean --fail-fast

# Releases the GITHUB_REF_NAME.
release: \
	github-ref-name-probe \
	github-release-create \
	goreleaser-release
