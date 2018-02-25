package merge_test

import (
	"path/filepath"

	"github.com/Prutswonder/submerger/merge"
	"github.com/stretchr/testify/mock"
)

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

type mockFileWalker struct {
	merge.FileWalker
	mock.Mock
}

func (m *mockFileWalker) Walk(root string, walkFn filepath.WalkFunc) error {
	i := m.Called(root, walkFn)
	return i.Error(0)
}
