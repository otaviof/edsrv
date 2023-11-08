package editor

import "github.com/otaviof/edsrv/pkg/edsrv/file"

type FakeEditor struct {
	payload []byte
}

var _ Interface = &FakeEditor{}

func (*FakeEditor) GetCommand() string {
	return "fake-editor"
}

func (*FakeEditor) GetTmpDir() string {
	return "none"
}

func (e *FakeEditor) Edit([]byte) (file.Interface, error) {
	return file.NewFakeFile("fake", e.payload), nil
}

func NewFakeEditor(payload []byte) *FakeEditor {
	return &FakeEditor{
		payload: payload,
	}
}
