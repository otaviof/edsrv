package file

import (
	"log/slog"
	"os"
)

// File represents the temporary file to be edited by an external editor.
type File struct {
	name string // full path to temporary file
	size int    // original size bytes
}

// LoggerWith decorates logger with current File instance attributes.
func (f *File) LoggerWith(logger *slog.Logger) *slog.Logger {
	return logger.With("file", f.name, "size", f.size)
}

// Name shows file name, full location to temporary file.
func (f *File) Name() string {
	return f.name
}

// Read reads temporary file content.
func (f *File) Read() ([]byte, error) {
	return os.ReadFile(f.name)
}

// Remove removes the temporary file.
func (f *File) Remove() error {
	return os.Remove(f.name)
}

// NewFile instantiate a new temporary file on the informed directory, and using
// the informed payload for its contents.
func NewFile(tmpDir string, payload []byte) (*File, error) {
	f, err := os.CreateTemp(tmpDir, "edsrv-*")
	if err != nil {
		return nil, err
	}
	if _, err = f.Write(payload); err != nil {
		return nil, err
	}
	if err = f.Close(); err != nil {
		return nil, err
	}
	return &File{name: f.Name(), size: len(payload)}, nil
}
