package memfs

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/util"
	"github.com/stretchr/testify/assert"
)

func TestRootExists(t *testing.T) {
	fs := New()
	f, err := fs.Stat("/")
	assert.NoError(t, err)
	assert.True(t, f.IsDir())
}

func TestCapabilities(t *testing.T) {
	fs := New()
	_, ok := fs.(billy.Capable)
	assert.True(t, ok)

	caps := billy.Capabilities(fs)
	assert.Equal(t, billy.DefaultCapabilities&^billy.LockCapability, caps)
}

func TestNegativeOffsets(t *testing.T) {
	fs := New()
	f, err := fs.Create("negative")
	assert.NoError(t, err)

	buf := make([]byte, 100)
	_, err = f.ReadAt(buf, -100)
	assert.ErrorContains(t, err, "readat negative: negative offset")

	_, err = f.Seek(-100, io.SeekCurrent)
	assert.NoError(t, err)
	_, err = f.Write(buf)
	assert.ErrorContains(t, err, "writeat negative: negative offset")
}

func TestExclusive(t *testing.T) {
	fs := New()
	f, err := fs.OpenFile("exclusive", os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
	assert.NoError(t, err)

	fmt.Fprint(f, "mememememe")

	err = f.Close()
	assert.NoError(t, err)

	_, err = fs.OpenFile("exclusive", os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
	assert.ErrorContains(t, err, os.ErrExist.Error())
}

func TestOrder(t *testing.T) {
	var err error

	files := []string{
		"a",
		"b",
		"c",
	}
	fs := New()
	for _, f := range files {
		_, err = fs.Create(f)
		assert.NoError(t, err)
	}

	attempts := 30
	for n := 0; n < attempts; n++ {
		actual, err := fs.ReadDir("")
		assert.NoError(t, err)

		for i, f := range files {
			assert.Equal(t, actual[i].Name(), f)
		}
	}
}

func TestNotFound(t *testing.T) {
	fs := New()
	files, err := fs.ReadDir("asdf")
	assert.Len(t, files, 0)
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestTruncateAppend(t *testing.T) {
	fs := New()
	err := util.WriteFile(fs, "truncate_append", []byte("file-content"), 0666)
	assert.NoError(t, err)

	f, err := fs.OpenFile("truncate_append", os.O_WRONLY|os.O_TRUNC|os.O_APPEND, 0666)
	assert.NoError(t, err)

	n, err := f.Write([]byte("replace"))
	assert.NoError(t, err)
	assert.Equal(t, n, len("replace"))

	err = f.Close()
	assert.NoError(t, err)

	data, err := util.ReadFile(fs, "truncate_append")
	assert.NoError(t, err)
	assert.Equal(t, string(data), "replace")
}

func TestReadlink(t *testing.T) {
	tests := []struct {
		name    string
		link    string
		want    string
		wantErr *error
	}{
		{
			name:    "symlink not found",
			link:    "/404",
			wantErr: &os.ErrNotExist,
		},
		{
			name: "self-targeting symlink",
			link: "/self",
			want: "/self",
		},
		{
			name: "symlink",
			link: "/bar",
			want: "/foo",
		},
		{
			name: "symlink to windows path",
			link: "/win",
			want: "c:\\test\\123",
		},
		{
			name: "symlink to network path",
			link: "/net",
			want: "\\test\\123",
		},
	}

	// Cater for memfs not being os-agnostic.
	if runtime.GOOS == "windows" {
		tests[1].want = "\\self"
		tests[2].want = "\\foo"
		tests[3].want = "\\c:\\test\\123"
	}

	fs := New()

	// arrange fs for tests.
	assert.NoError(t, fs.Symlink("/self", "/self"))
	assert.NoError(t, fs.Symlink("/foo", "/bar"))
	assert.NoError(t, fs.Symlink("c:\\test\\123", "/win"))
	assert.NoError(t, fs.Symlink("\\test\\123", "/net"))

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := fs.Readlink(tc.link)

			if tc.wantErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			} else {
				assert.ErrorIs(t, err, *tc.wantErr)
			}
		})
	}
}

func TestSymlink2(t *testing.T) {
	tests := []struct {
		name    string
		target  string
		link    string
		want    string
		wantErr string
	}{
		{
			name:   "new symlink unexistent target",
			target: "/bar",
			link:   "/foo",
			want:   "/bar",
		},
		{
			name:   "self-targeting symlink",
			target: "/self",
			link:   "/self",
			want:   "/self",
		},
		{
			name:   "new symlink to file",
			target: "/file",
			link:   "/file-link",
			want:   "/file",
		},
		{
			name:   "new symlink to dir",
			target: "/dir",
			link:   "/dir-link",
			want:   "/dir",
		},
		{
			name:   "new symlink to win",
			target: "c:\\foor\\bar",
			link:   "/win",
			want:   "c:\\foor\\bar",
		},
		{
			name:   "new symlink to net",
			target: "\\net\\bar",
			link:   "/net",
			want:   "\\net\\bar",
		},
		{
			name:   "new symlink to net",
			target: "\\net\\bar",
			link:   "/net",
			want:   "\\net\\bar",
		},
		{
			name:    "duplicate symlink",
			target:  "/bar",
			link:    "/foo",
			wantErr: os.ErrExist.Error(),
		},
		{
			name:    "symlink over existing file",
			target:  "/foo/bar",
			link:    "/file",
			want:    "/file",
			wantErr: os.ErrExist.Error(),
		},
	}

	// Cater for memfs not being os-agnostic.
	if runtime.GOOS == "windows" {
		tests[0].want = "\\bar"
		tests[1].want = "\\self"
		tests[2].want = "\\file"
		tests[3].want = "\\dir"
		tests[4].want = "\\c:\\foor\\bar"
	}

	fs := New()

	// arrange fs for tests.
	err := fs.MkdirAll("/dir", 0o600)
	assert.NoError(t, err)
	_, err = fs.Create("/file")
	assert.NoError(t, err)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := fs.Symlink(tc.target, tc.link)

			if tc.wantErr == "" {
				got, err := fs.Readlink(tc.link)
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			} else {
				assert.ErrorContains(t, err, tc.wantErr)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	tests := []struct {
		name string
		elem []string
		want string
	}{
		{name: "empty", elem: []string{""}, want: ""},
		{name: "c:", elem: []string{"C:"}, want: "C:"},
		{name: "simple rel", elem: []string{"a", "b", "c"}, want: "a/b/c"},
		{name: "simple rel backslash", elem: []string{"\\", "a", "b", "c"}, want: "\\/a/b/c"},
		{name: "simple abs slash", elem: []string{"/", "a", "b", "c"}, want: "/a/b/c"},
		{name: "c: rel", elem: []string{"C:\\", "a", "b", "c"}, want: "C:\\/a/b/c"},
		{name: "c: abs", elem: []string{"/C:\\", "a", "b", "c"}, want: "/C:\\/a/b/c"},
		{name: "\\ rel", elem: []string{"\\\\", "a", "b", "c"}, want: "\\\\/a/b/c"},
		{name: "\\ abs", elem: []string{"/\\\\", "a", "b", "c"}, want: "/\\\\/a/b/c"},
	}

	// Cater for memfs not being os-agnostic.
	if runtime.GOOS == "windows" {
		tests[1].want = "C:."
		tests[2].want = "a\\b\\c"
		tests[3].want = "\\a\\b\\c"
		tests[4].want = "\\a\\b\\c"
		tests[5].want = "C:\\a\\b\\c"
		tests[6].want = "\\C:\\a\\b\\c"
		tests[7].want = "\\\\a\\b\\c"
		tests[8].want = "\\\\\\a\\b\\c"
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := New().Join(tc.elem...)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestSymlink(t *testing.T) {
	fs := New()
	err := fs.Symlink("test", "test")
	assert.NoError(t, err)

	f, err := fs.Open("test")
	assert.NoError(t, err)
	assert.NotNil(t, f)

	fi, err := fs.ReadDir("test")
	assert.NoError(t, err)
	assert.Nil(t, fi)
}
