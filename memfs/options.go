package memfs

type options struct {
}

func newOptions() *options {
	return &options{}
}

type Option func(*options)
