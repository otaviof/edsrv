package editor

import "github.com/otaviof/edsrv/pkg/edsrv/file"

type Interface interface {
	// GetCommand shows the Editor command in use.
	GetCommand() string

	// GetTmpDir shows the temporary directory in use.
	GetTmpDir() string

	// Edit edits the payload using the external editor.
	Edit([]byte) (file.Interface, error)
}
