Contributing to `edsrv`
-----------------------

This document describes the highlevel automation to interact with the project's code-base. The required tools are automatically installed with `go install` when not yet installed on the environment.

## Build

The executables are built by `make`'s default target, as the example below. The executable is placed under `./bin` directory configured on [`Makefile`'s](./Makefile) `OUTPUT_DIR` variable.

To build the application using your local `go` installation, run:

```sh
make
```

For `snapshot` target it uses [`goreleaser`][goreleaser] to built the application:

```sh
make snapshot
```

## Lint

Static code analysis is done with [`golangci-lint`][golangciLint] and rules defined [here](./.golangci.yaml).

```sh
make lint
```

## CI

Continuous integration (CI) [tests are declared here](.github/workflows/test.yaml), where the following `Makefile` targets are invoked in order to test different aspects of the project.

In short, to run all tests:

```sh
make test
```

### Unit Testing

```sh
make test-unit
```

### Integration Testing

```sh
make test-integration
```

[goreleaser]: https://goreleaser.com/
[golangciLint]: https://github.com/golangci/golangci-lint
