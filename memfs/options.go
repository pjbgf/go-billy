package memfs

import (
	"math"
	"os"
	"path/filepath"
)

const (
	defaultDirMode     = 0o755
	defaultFileMode    = 0o666
	defaultSymlinkMode = 0o777

	// memfs FS does not necessarily require hard limits, however by applying
	// sensible limits that are aligned with what osfs will enforce, that results
	// on similar user experience across implementations.
	//
	// https://www.gnu.org/software/libc/manual/html_node/File-Minimums.html
	// https://learn.microsoft.com/en-us/windows/win32/fileio/maximum-file-path-limitation
	// https://manpages.opensuse.org/Tumbleweed/man-pages-posix/limits.h.0p.en.html
	defaultMaxSymlink = 4096
	defaultMaxName    = 255
	defaultMaxPath    = 4096
)

type options struct {
	legacy bool

	dirMode     os.FileMode
	fileMode    os.FileMode
	symlinkMode os.FileMode
}

func newOptions() *options {
	return &options{
		dirMode:     defaultDirMode,
		fileMode:    defaultFileMode,
		symlinkMode: defaultSymlinkMode,
	}
}

type Option func(*options)

// WithLegacy enables legacy mode, making the Memory fs to
// operate at a mixed mode where some operations are OS agnostic
// and others aren't.
func WithLegacy() Option {
	return func(o *options) {
		o.legacy = true
	}
}

func (o *options) separator() string {
	if o.legacy {
		return string(filepath.Separator)
	}
	return "/"
}

func (o *options) maxSymlink() int {
	if o.legacy {
		return math.MaxInt
	}
	return defaultMaxSymlink
}

func (o *options) maxPath() int {
	if o.legacy {
		return math.MaxInt
	}
	return defaultMaxPath
}
