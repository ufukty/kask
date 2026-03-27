package builder

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"go.ufukty.com/kask/internal/disk"
	"go.ufukty.com/kask/internal/disk/memory"
	"go.ufukty.com/kask/internal/paths"
	"go.ufukty.com/kask/internal/rewriter"
	"go.ufukty.com/kask/pkg/kask"
)

func check(tmp fs.StatFS, path string) bool {
	_, err := tmp.Stat(path)
	return err == nil
}

func dfs(n *kask.Node) []*kask.Node {
	cs := []*kask.Node{n}
	for _, c := range n.Children {
		cs = append(cs, dfs(c)...)
	}
	return cs
}

func hrefs(root *kask.Node) []string {
	ss := []string{}
	for _, n := range dfs(root) {
		ss = append(ss, n.Href)
	}
	return ss
}

func titles(root *kask.Node) []string {
	ss := []string{}
	for _, n := range dfs(root) {
		ss = append(ss, n.Title)
	}
	return ss
}

func breadcrumbs(root *kask.Node) []string {
	bs := []string{root.Title}
	for _, c := range root.Children {
		for _, b := range breadcrumbs(c) {
			bs = append(bs, fmt.Sprintf("%s / %s", root.Title, b))
		}
	}
	return bs
}

func buildTestSite(t *testing.T, src, domain string) (*builder, *memory.Dir) {
	dst := memory.New()
	b := newBuilder(builderArgs{
		Src:     disk.NewReal(filepath.Join("testdata", src)), // TODO: use on-memory FS
		Dst:     dst,
		Domain:  domain,
		Dev:     true,
		Verbose: false,
	})
	if err := b.Build(); err != nil {
		t.Fatalf("prep, builder.Build: %v", err)
	}
	return b, dst
}

func lines(ls ...string) string {
	return strings.Join(ls, "\n")
}

func assertfile(t *testing.T, fs fs.StatFS, path string) {
	t.Run(strings.ReplaceAll(path, "/", "\\"), func(t *testing.T) {
		if !check(fs, path) {
			t.Log(fs)
			t.Errorf("assert, file not found: %s", path)
		}
	})
}

func files(dst fs.ReadDirFS) ([]string, error) {
	ss := []string{}
	err := fs.WalkDir(dst, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			ss = append(ss, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk: %w", err)
	}
	return ss, nil
}

var anchor = regexp.MustCompile(`<a[^>]*>[^<]*</a>`)

func findAnchorTags(fs disk.ReadFS, path string) ([]string, error) {
	f, err := fs.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading: %w", err)
	}
	return anchor.FindAllString(string(f), -1), nil
}

func fixture() *builder {
	rw := rewriter.New(paths.Paths{Src: ".", Dst: ".", Url: "https://kask.ufukty.com/"})
	m := map[string]string{
		// leaves
		"a/page.tmpl":   "https://kask.ufukty.com/a/page.html",
		"a/index.tmpl":  "https://kask.ufukty.com/a/",
		"a/b/README.md": "https://kask.ufukty.com/a/b/",
		"a/b/page.md":   "https://kask.ufukty.com/a/b/page.html",

		// visitable dirs:
		".":   "https://kask.ufukty.com/",
		"a/":  "https://kask.ufukty.com/a/",
		"a/b": "https://kask.ufukty.com/a/b/",

		// assets
		".assets/font.woff2":             "https://kask.ufukty.com/.assets/font.woff2",
		"a/.assets/img.jpg":              "https://kask.ufukty.com/a/.assets/img.jpg",
		"a/.assets/img@2x.jpg":           "https://kask.ufukty.com/a/.assets/img%402x.jpg",
		"a/.assets/img@3x.jpg":           "https://kask.ufukty.com/a/.assets/img%403x.jpg",
		"a/.assets/poster.jpg":           "https://kask.ufukty.com/a/.assets/poster.jpg",
		"a/.assets/video.mp4":            "https://kask.ufukty.com/a/.assets/video.mp4",
		"a/.assets/og.jpg":               "https://kask.ufukty.com/a/.assets/og.jpg",
		"a/.assets/embedded-player.html": "https://kask.ufukty.com/a/.assets/embedded-player.html",
	}
	for s, u := range m {
		rw.Bank(s, u)
	}
	return &builder{rw: rw, incorrectLinks: map[string][]string{}}
}
