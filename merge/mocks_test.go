package merge_test

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/Prutswonder/submerger/merge"
	"github.com/stretchr/testify/mock"
)

/* Logger mock */

type mockLogger struct {
	merge.Logger
	mock.Mock
}

func (m *mockLogger) Printf(format string, v ...interface{}) {
	m.Called(format, v)
}

func (m *mockLogger) Println(v ...interface{}) {
	m.Called(v)
}

/* FileWarker mock */

type mockFileWalker struct {
	merge.FileWalker
	mock.Mock
}

func (m *mockFileWalker) Walk(root string, walkFn filepath.WalkFunc) error {
	i := m.Called(root, walkFn)
	return i.Error(0)
}

/* Commander mock */

type mockCommander struct {
	merge.Commander
	mock.Mock
}

func (m *mockCommander) Command(name string, arg ...string) merge.Cmd {
	i := m.Called(name, arg)
	return i.Get(0).(merge.Cmd)
}

/* Cmd mock */

type mockCmd struct {
	merge.Cmd
	mock.Mock
}

func (m *mockCmd) StdoutPipe() (io.ReadCloser, error) {
	i := m.Called()
	return i.Get(0).(io.ReadCloser), i.Error(1)
}

func (m *mockCmd) SetEnvironment(env []string) {
	m.Called(env)
}

func (m *mockCmd) StderrPipe() (io.ReadCloser, error) {
	i := m.Called()
	return i.Get(0).(io.ReadCloser), i.Error(1)
}

func (m *mockCmd) Start() error {
	i := m.Called()
	return i.Error(0)
}

func (m *mockCmd) Wait() error {
	i := m.Called()
	return i.Error(0)
}

/* os.FileInfo mock */

type mockFileInfo struct {
	os.FileInfo
	mock.Mock
}

func (m *mockFileInfo) Name() string {
	// base name of the file
	i := m.Called()
	return i.String(0)
}

func (m *mockFileInfo) Size() int64 {
	// length in bytes for regular files; system-dependent for others
	i := m.Called()
	return i.Get(0).(int64)
}

func (m *mockFileInfo) Mode() os.FileMode {
	// file mode bits
	i := m.Called()
	return i.Get(0).(os.FileMode)
}

func (m *mockFileInfo) ModTime() time.Time {
	// modification time
	i := m.Called(0)
	return i.Get(0).(time.Time)
}

func (m *mockFileInfo) IsDir() bool {
	// abbreviation for Mode().IsDir()
	i := m.Called(0)
	return i.Bool(0)
}

func (m *mockFileInfo) Sys() interface{} {
	// underlying data source (can return nil)
	i := m.Called(0)
	return i.Get(0)
}
