package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"go.ufukty.com/kask/internal/disk"
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
