package merge

import "path/filepath"

type (
	// FileWalker is a wrapper for filepath.Walk()
	FileWalker interface {
		Walk(root string, walkFn filepath.WalkFunc) error
	}

	fileWalkerImpl struct{}
)

// NewFileWalker instantiates a new FileWalker.
func NewFileWalker() FileWalker {
	return &fileWalkerImpl{}
}

func (f fileWalkerImpl) Walk(root string, walkFn filepath.WalkFunc) error {
	return filepath.Walk(root, walkFn)
}
