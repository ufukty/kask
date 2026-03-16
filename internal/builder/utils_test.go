package builder

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"go.ufukty.com/kask/internal/disk"
	"go.ufukty.com/kask/internal/paths"
	"go.ufukty.com/kask/internal/rewriter"
	"go.ufukty.com/kask/pkg/kask"
)

func check(tmp, path string) bool {
	_, err := os.Stat(filepath.Join(tmp, path))
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

func files(path string) []string {
	ss := []string{}
	fs.WalkDir(os.DirFS(path), ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			ss = append(ss, path)
		}
		return nil
	})
	return ss
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
