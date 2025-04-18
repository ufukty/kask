package directory

import "path/filepath"

func (d *Dir) Path() string {
	return filepath.Join(d.SiteRoot, d.InSitePath)
}
