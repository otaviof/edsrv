package file

import "log/slog"

type FakeFile struct {
	name    string
	payload []byte
}

var _ Interface = &FakeFile{}

func (*FakeFile) LoggerWith(logger *slog.Logger) *slog.Logger {
	return logger
}

func (f *FakeFile) Name() string {
	return f.name
}

func (f *FakeFile) Read() ([]byte, error) {
	return f.payload, nil
}

func (*FakeFile) Remove() error {
	return nil
}

func NewFakeFile(name string, payload []byte) *FakeFile {
	return &FakeFile{
		name:    name,
		payload: payload,
	}
}
