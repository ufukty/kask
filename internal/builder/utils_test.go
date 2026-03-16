package builder

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"go.ufukty.com/kask/internal/disk"
	"go.ufukty.com/kask/internal/paths"
	"go.ufukty.com/kask/internal/rewriter"
	"go.ufukty.com/kask/pkg/kask"
)

func tescape(s string) string {
	return strings.ReplaceAll(s, "/", "\\")
}

func testname(parent, child string) string {
	return fmt.Sprintf("%s=>%s", tescape(parent), tescape(child))
}

func check(tmp, path string) bool {
	_, err := os.Stat(filepath.Join(tmp, path))
	return err == nil
}

func dfsWithAncestry(n []*kask.Node, f func([]*kask.Node)) {
	f(n)
	for _, c := range n[len(n)-1].Children {
		dfsWithAncestry(append(slices.Clone(n), c), f)
	}
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

func lines(ls ...string) string {
	return strings.Join(ls, "\n")
}

func readFile(path string) string {
	c, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("os.ReadFile: %w", err))
	}
	return string(c)
}

func assertfile(t *testing.T, tmp, path string) {
	t.Run(strings.ReplaceAll(path, "/", "\\"), func(t *testing.T) {
		if !check(tmp, path) {
			t.Log(tmp)
			t.Errorf("assert, file not found: %s", path)
		}
	})
}

func printFiles(path string) {
	fs.WalkDir(os.DirFS(path), ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			fmt.Println(path)
		}
		return nil
	})
}

var anchor = regexp.MustCompile(`<a[^>]*>[^<]*</a>`)

func findAnchorTags(path string) []string {
	return anchor.FindAllString(readFile(path), -1)
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
