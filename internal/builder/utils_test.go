package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"go.ufukty.com/kask/internal/disk"
	"go.ufukty.com/kask/internal/paths"
	"go.ufukty.com/kask/internal/rewriter"
	"go.ufukty.com/kask/pkg/kask"
)

func check(tmp, path string) bool {
	_, err := os.Stat(filepath.Join(tmp, path))
	return err == nil
}

func dfs(n []*kask.Node, f func([]*kask.Node)) {
	f(n)
	for _, c := range n[len(n)-1].Children {
		dfs(append(slices.Clone(n), c), f)
	}
}

func each(ns []*kask.Node, f func(*kask.Node) string) []string {
	ss := make([]string, 0, len(ns))
	for _, n := range ns {
		ss = append(ss, f(n))
	}
	return ss
}

func buildTestSite(path, domain string) (*builder, string) {
	tmp, err := os.MkdirTemp(os.TempDir(), "kask-test-build-*")
	if err != nil {
		panic(fmt.Errorf("buildTestSite: os.MkdirTemp: %w", err))
	}
	b := newBuilder(builderArgs{
		Src:     disk.NewReal(path), // TODO: use on-memory FS
		Dst:     disk.NewReal(tmp),  // TODO: use on-memory FS
		Domain:  domain,
		Dev:     true,
		Verbose: false,
	})
	err = b.Build()
	if err != nil {
		panic(fmt.Errorf("buildTestSite: b.Build: %w", err))
	}
	return b, tmp
}

func readFile(path string) string {
	c, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("os.ReadFile: %w", err))
	}
	return string(c)
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
