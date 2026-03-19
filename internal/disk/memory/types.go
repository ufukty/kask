package memory

import (
	"io"

	"go.ufukty.com/kask/internal/disk"
)

type (
	File []byte
	Dir  map[string]any
)

// write
var (
	_ io.WriteCloser = (*File)(nil)
	_ disk.WriteFS   = (*Dir)(nil)
)
