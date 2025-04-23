package fsys

import (
	"context"
	"sync"
	"time"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

type FileNode struct {
	name    string
	content []byte
	mu      sync.Mutex
}

type File interface {
	fs.Node
	fs.Handle
	fs.HandleReader
	fs.HandleWriter
}

func NewFile(name string) *FileNode {
	return &FileNode{name: name}
}

var _ File = (*FileNode)(nil)

func (f *FileNode) Attr(ctx context.Context, a *fuse.Attr) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	a.Mode = 0644
	a.Size = uint64(len(f.content))
	a.Mtime = time.Now()
	return nil
}

func (f *FileNode) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	start := int(req.Offset)
	if start >= len(f.content) {
		resp.Data = []byte{}
		return nil
	}
	end := start + req.Size
	if end > len(f.content) {
		end = len(f.content)
	}
	resp.Data = f.content[start:end]
	return nil
}

func (f *FileNode) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	end := int(req.Offset) + len(req.Data)
	if end > len(f.content) {
		tmp := make([]byte, end)
		copy(tmp, f.content)
		f.content = tmp
	}
	copy(f.content[req.Offset:], req.Data)
	resp.Size = len(req.Data)
	return nil
}
