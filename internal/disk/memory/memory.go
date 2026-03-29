package memory

import (
	"io/fs"

	"go.ufukty.com/kask/internal/disk"
)

type (
	File []byte
	Dir  map[string]any // use [New] to construct
)


var (
	// write
	_ disk.WriteFS = (*Dir)(nil)
	// read
	_ fs.FS         = (*Dir)(nil)
	_ fs.ReadFileFS = (*Dir)(nil)
	_ fs.ReadDirFS  = (*Dir)(nil)
	_ fs.StatFS     = (*Dir)(nil)
	_ disk.ReadFS   = (*Dir)(nil)
)
