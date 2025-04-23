package fsys

import (
	"context"
	"fmt"
	"os"
	"sync"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

type Directory interface {
	fs.Node
	fs.Handle
	fs.NodeStringLookuper
	fs.HandleReadDirAller
	fs.NodeCreater
	fs.NodeMkdirer
}

type DirNode struct {
	name     string
	children map[string]fs.Node
	mu       sync.Mutex
}

func NewDir(name string) *DirNode {
	return &DirNode{name: name, children: make(map[string]fs.Node)}
}

var _ Directory = (*DirNode)(nil)

func (d *DirNode) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Mode = os.ModeDir | 0755
	return nil
}

func (d *DirNode) Lookup(ctx context.Context, name string) (fs.Node, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	child, ok := d.children[name]
	if !ok {
		return nil, fuse.ENOENT
	}
	fmt.Println("Lookup:", name)
	return child, nil
}

func (d *DirNode) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	var entries []fuse.Dirent
	for name, node := range d.children {
		entry := fuse.Dirent{Name: name}
		switch node.(type) {
		case *DirNode:
			entry.Type = fuse.DT_Dir
		case *FileNode:
			entry.Type = fuse.DT_File
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (d *DirNode) Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (fs.Node, fs.Handle, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	f := NewFile(req.Name)
	d.children[req.Name] = f
	fmt.Println("Create file:", req.Name)
	return f, f, nil
}

func (d *DirNode) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fs.Node, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	dir := NewDir(req.Name)
	d.children[req.Name] = dir
	fmt.Println("Create dir:", req.Name)
	return dir, nil
}
