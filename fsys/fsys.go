package fsys

import (
	"bazil.org/fuse/fs"
)

type FS struct {
	root Directory
}

func NewFS(root Directory) *FS {
	return &FS{root: root}
}

func (f *FS) Root() (fs.Node, error) {
	return f.root, nil
}
