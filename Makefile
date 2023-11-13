SHELL := /bin/sh

# application name and build output directory
APP := edsrv
OUTPUT_DIR ?= bin

# application directories
CMD ?= ./cmd/$(APP)/...
PKG ?= ./pkg/$(APP)/...

# golang build and test configuration
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

# github action current ref name, and credentials
GITHUB_REF_NAME ?= ${GITHUB_REF_NAME:-}
GITHUB_TOKEN ?= ${GITHUB_TOKEN:-}

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
	goreleaser build --clean --snapshot --single-target -o=$(BIN) $(ARGS)

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

# Installs the plist service definition for macOS.
.PHONY: install-launchagent
install-launchagent: install
	install -g staff -m 0755 $(PLIST) $(LAUNCHAGENT_PLIST)

# Loads the launch-agent service file.
.PHONY: launchctl-load
launchctl-load:
	launchctl load -w $(LAUNCHAGENT_PLIST)

# Auxiliary target to run "sleep" command.
.PHONY: sleep
sleep:
	@sleep 1

# Shows the macOS service status.
.PHONY: launchctl-list
launchctl-list:
	launchctl list $(LAUNCHAGENT_LABEL)

# Calls the application to check the service status
.PHONY: status
status:
	$(APP) status $(ARGS)

# Deploys the macOS service and starts it.
.PHONY: deploy-launchd
deploy-launchd: \
	install-launchagent \
	launchctl-load \
	sleep \
	launchctl-list \
	status

# Unloads the launch-agent plist file.
.PHONY: launchctl-unload
launchctl-unload:
	launchctl unload -w $(LAUNCHAGENT_PLIST)

# Removes the launch-agent plist file.
remove-launchagent:
	rm -f -v $(LAUNCHAGENT_PLIST) || true

# Removes the whole macos service, the oppositive of "deploy" target.
.PHONY: remove-launchd
remove-launchd: \
	launchctl-unload \
	remove-launchagent

# Stops the macOS service.
.PHONY: launchctl-stop
launchctl-stop:
	launchctl stop $(LAUNCHAGENT_LABEL)

# Starts the macOS service.
.PHONY: launchctl-start
launchctl-start:
	launchctl start $(LAUNCHAGENT_LABEL)

# Restarts the macOS service.
restart-launchd: \
	launchctl-stop \
	sleep \
	launchctl-start \
	sleep \
	launchctl-list \
	status

#
# GitHub Release
#

# Asserts the required environment variables are set and the target release
# version starts with "v".
github-preflight:
ifeq ($(strip $(GITHUB_REF_NAME)),)
	$(error variable GITHUB_REF_NAME is not set)
endif
ifeq ($(shell echo ${GITHUB_REF_NAME} |grep -v -E '^v'),)
	@echo GITHUB_REF_NAME=\"${GITHUB_REF_NAME}\"
else
	$(error invalid GITHUB_REF_NAME, it must start with "v")
endif
ifeq ($(strip $(GITHUB_TOKEN)),)
	$(error variable GITHUB_TOKEN is not set)
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
	goreleaser release --clean --fail-fast $(ARGS)

# Releases the GITHUB_REF_NAME.
github-release: \
	github-preflight \
	github-release-create \
	goreleaser-release
