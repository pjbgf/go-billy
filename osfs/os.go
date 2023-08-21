//go:build !js
// +build !js

// Package osfs provides a billy filesystem for the OS.
package osfs

import (
	"io/fs"
	"os"
	"sync"

	"github.com/go-git/go-billy/v5"
)

const (
	defaultDirectoryMode = 0o755
	defaultCreateMode    = 0o666
)

// New returns a new OS filesystem.
func New(baseDir string, opts ...Option) billy.Filesystem {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	return o.Type.New(baseDir)
}

// file is a wrapper for an os.File which adds support for file locking.
type file struct {
	*os.File
	m sync.Mutex
}

func (t Type) New(baseDir string) billy.Filesystem {
	if t == BoundOSFS {
		return newBoundOS(baseDir)
	}

	return newChrootOS(baseDir)
}

func readDir(dir string) ([]os.FileInfo, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	infos := make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		fi, err := entry.Info()
		if err != nil {
			return nil, err
		}
		infos = append(infos, fi)
	}
	return infos, nil
}

func tempFile(dir, prefix string) (billy.File, error) {
	f, err := os.CreateTemp(dir, prefix)
	if err != nil {
		return nil, err
	}
	return &file{File: f}, nil
}
