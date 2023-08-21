package osfs

type Option func(*options)

type options struct {
	Type
}

// WithBoundOS returns the option of using a Bound filesystem OS.
func WithBoundOS() Option {
	return func(o *options) {
		o.Type = BoundOSFS
	}
}

type Type int

const (
	ChrootOSFS Type = iota
	BoundOSFS
)
