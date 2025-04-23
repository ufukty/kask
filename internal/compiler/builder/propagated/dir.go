package propagated

import (
	"fmt"
	"html/template"
	"maps"
	"path/filepath"
	"slices"

	"github.com/ufukty/kask/internal/compiler/builder/directory"
	"github.com/ufukty/kask/internal/compiler/builder/propagated/bundle"
)

type Dir struct {
	Name    string
	Assets  string
	Subdirs []*Dir
	Pages   struct {
		Markdown []string
		Html     []string
	}
	Tmpl        *template.Template
	Stylesheets []string // paths
}

type args struct {
	Dev bool
}

type context struct {
	Path        string
	Template    *template.Template
	Stylesheets []string // paths
	Page        string   // path
}

type artifacts struct {
	Stylesheets map[string]string // [path: content]
}

func (a *artifacts) merge(b *artifacts) {
	maps.Copy(a.Stylesheets, b.Stylesheets)
}

func d(dir *directory.Dir, ctx context, args args) (*Dir, *artifacts, error) {
	tmpl, err := ctx.Template.Clone()
	if err != nil {
		return nil, nil, fmt.Errorf("cloning propagated templates: %w", err)
	}

	dir2 := &Dir{
		Name:        dir.Name,
		Assets:      dir.Assets,
		Subdirs:     []*Dir{},
		Pages:       dir.Pages,
		Tmpl:        tmpl,
		Stylesheets: slices.Clone(ctx.Stylesheets),
	}

	artf := &artifacts{
		Stylesheets: map[string]string{},
	}

	if dir.Kask != nil && dir.Kask.Propagate != nil {
		if len(dir.Kask.Propagate.Css) > 0 {
			dst := filepath.Join(ctx.Path, "styles.propagate.css")
			css, err := bundle.Files(dir.Kask.Propagate.Css)
			if err != nil {
				return nil, nil, fmt.Errorf("bundling .kask/propagate/*.css: %w", err)
			}
			artf.Stylesheets[dst] = css
			dir2.Stylesheets = append(dir2.Stylesheets, dst)
		}
		if len(dir.Kask.Propagate.Tmpl) > 0 {
			propTmpl, err := dir2.Tmpl.ParseFiles(dir.Kask.Propagate.Tmpl...)
			if err != nil {
				return nil, nil, fmt.Errorf("loading .kask/propagate/*.tmpl: %w", err)
			}
			ctx.Template = propTmpl
			tmpl, err = propTmpl.Clone()
			if err != nil {
				return nil, nil, fmt.Errorf("cloning for non-propagated templates at same folder: %w", err)
			}
		}
	}

	if dir.Kask != nil {
		if len(dir.Kask.Css) > 0 {
			dst := filepath.Join(ctx.Path, "styles.css")
			css, err := bundle.Files(dir.Kask.Css)
			if err != nil {
				return nil, nil, fmt.Errorf("bundling .kask/*.css: %w", err)
			}
			artf.Stylesheets[dst] = css
			dir2.Stylesheets = append(dir2.Stylesheets, dst)
		}
		if len(dir.Kask.Tmpl) > 0 {
			dir2.Tmpl, err = tmpl.ParseFiles(dir.Kask.Tmpl...)
			if err != nil {
				return nil, nil, fmt.Errorf("loading .kask/*.tmpl: %w", err)
			}
		}
	}

	for _, child := range dir.Subdirs {
		subdir, subartf, err := d(child, ctx, args)
		if err != nil {
			return nil, nil, fmt.Errorf("%s: %w", child.Name, err)
		}
		dir2.Subdirs = append(dir2.Subdirs, subdir)
		artf.merge(subartf)
	}

	return dir2, artf, nil
}

func Directory(root *directory.Dir, dev bool) (*Dir, *artifacts, error) {
	return d(root, context{
		Path:        "",
		Template:    template.New("page"),
		Stylesheets: []string{},
		Page:        "",
	}, args{
		Dev: dev,
	})
}
