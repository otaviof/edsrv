<p align="center">
    <a alt="CI" href="https://github.com/otaviof/edsrv/actions">
        <img src="https://github.com/otaviof/edsrv/actions/workflows/test.yaml/badge.svg">
    </a>
    <a alt="project quality report" href="https://goreportcard.com/report/github.com/otaviof/edsrv">
        <img src="https://goreportcard.com/badge/github.com/otaviof/edsrv">
    </a>
	<a alt="latest release" href="https://github.com/otaviof/edsrv/releases/latest">
		<img src="https://img.shields.io/github/v/release/otaviof/edsrv">
	</a>
</p>

`edsrv`
-------

Is a "edit server" (`edsrv`) backend for browser extensions like [TextEditAid][textEditAid], [TextAidToo][textAidToo], [withExEditor][withExEditor] and possibly more supporting the same API.

# Installation

The preferred installation method is using a pre-compiled executable available on [releases][repoReleases].

Alternatively, use the [`Makefile`](./Makefile) target `install` to build and copy the `edsrv` executable to `/usr/local/bin` (`PREFIX`):

```sh
make install
```

Another supported approach is using `go install`:

```sh
go install github.com/otaviof/edsrv/cmd/edsrv@latest
```

# Usage

The edit-server needs to run on the background in order to respond the text edit requests coming from your favorite browser extension.

To start the service run `edsrv start` with appropriate flags:

```sh
edsrv start --addr="127.0.0.1:8928" --tmp-dir="${TMPDIR}" --editor="${EDITOR}"
```

The subcommand `start` supports the following command-line flags:

| Flag        | Default          | Description                                 |
| :---------- | :--------------- | :------------------------------------------ |
| `--addr`    | `127.0.0.1:8929` | Listen address, interface and port          |
| `--tmp-dir` | `${TMPDIR}`      | Temporary directory to store edited payload |
| `--editor`  | `${EDITOR}`      | Editor to edit the payload                  |

By default `edsrv start` uses the regular temporary directory configured on your shell (`${TMPDIR}`) and editor (`${EDITOR}`).

Once the edit-server is running, you can use the `status` subcommand to confirm the servier is running, and peek runtime configuration:

```sh
edsrv status
```

## macOS Service

For macOS users, consider the [`edsrv.plist` launchd service file](./contrib/edsrv.plist) details, adapt to your needs. To deploy the launchd based service, run:

```sh
make deploy-launchd
```

# API

The edit-server API contains two only two endpoints, enabling users to edit the payload with the external editor, and check the service status.

## `POST /` 

Edits the request body payload using the external editor (`--editor`).

First, the body payload is stored on a new randomly named temporary file, under `--tmp-dir` directory. Then, the external editor (`--editor`) gets invoked blocking the request until completed. Once completed the response body carries the temporary file content and deletes it.

Thus, the `--editor` flag must be configured to wait until completed, like for instance `code -w`, `-w` implies the command line will *wait* until file is closed.

## `GET /status`

The endpoint shows the configured editor and temporary directory, i.e.:

```
$ curl -s 127.0.0.1:8928/status
editor='code -n -w', tmpDir='/tmp'
```

The same output is shown on `edsrv status` subcommand.

# Contributing

To know more details about the project automation please consider [CONTRIBUTING.md](./CONTRIBUTING.md).

[repoReleases]: https://github.com/otaviof/edsrv/releases
[textAidToo]: https://chrome.google.com/webstore/detail/text-aid-too/klbcooigafjpbiahdjccmajnaehomajc
[textEditAid]: https://chrome.google.com/webstore/detail/texteditaid/ppoadiihggafnhokfkpphojggcdigllp
[withExEditor]: https://chrome.google.com/webstore/detail/withexeditor/koghhpkkcndhhclklnnnhcpkkplfkgoi
