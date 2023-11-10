package editor

import (
	"log/slog"
	"os/exec"
	"strings"

	"github.com/otaviof/edsrv/pkg/edsrv/file"
)

// Editor represents the external editor.
type Editor struct {
	logger  *slog.Logger // shared logger instance
	command []string     // external editor command and arguments
	tmpDir  string       // path to temporary directory
}

var _ Interface = &Editor{}

// GetCommand shows the editor command.
func (e *Editor) GetCommand() string {
	return strings.Join(e.command, " ")
}

// GetTmpDir exposes the temporary directory location.
func (e *Editor) GetTmpDir() string {
	return e.tmpDir
}

// runCommandAndWait runs the editor command and waits for the result.
func (e *Editor) runCommandAndWait(f *file.File) error {
	script := e.command
	script = append(script, f.Name())

	logger := e.logger.With("script", script)
	logger.Info("running editor command and waiting...")

	cmd := exec.Command(script[0], script[1:]...) //nolint:gosec
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("error reading script combined output", "err", err.Error())
		return err
	}

	logger.Debug("editor command result", "output", string(output))
	return nil
}

// Edit edits the informed payload on a temporary file, using the external editor.
func (e *Editor) Edit(payload []byte) (file.Interface, error) {
	e.logger.Debug("creating temporary file for payload")
	f, err := file.NewFile(e.tmpDir, payload)
	if err != nil {
		return nil, err
	}
	f.LoggerWith(e.logger).Debug("temporary file created")
	if err = e.runCommandAndWait(f); err != nil {
		return nil, err
	}
	return f, nil
}

// NewEditor instantiates a new editor with the desired command and temporary
// directory.
func NewEditor(logger *slog.Logger, command, tmpDir string) *Editor {
	return &Editor{
		logger:  logger,
		command: strings.Split(command, " "),
		tmpDir:  tmpDir,
	}
}
