package memfs

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeparator(t *testing.T) {
	tests := []struct {
		name   string
		option Option
		want   string
	}{
		{name: "Default", option: func(o *options) {}, want: "/"},
		{name: "WithLegacy", option: WithLegacy(), want: string(filepath.Separator)},
	}

	for _, tc := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			o := &options{}
			tc.option(o)
			got := o.separator()
			assert.Equal(t, tc.want, got)
		})
	}
}
