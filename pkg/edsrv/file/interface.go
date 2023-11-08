package file

import "log/slog"

type Interface interface {
	// LoggerWith decorates logger with File attributes.
	LoggerWith(*slog.Logger) *slog.Logger

	// Name shows the file name, full path location.
	Name() string

	// Read reads file contents.
	Read() ([]byte, error)

	// Remove removes the temporary file.
	Remove() error
}
