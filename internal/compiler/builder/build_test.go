package builder

import (
	"fmt"
	"os"
	"testing"

	"github.com/ufukty/kask/internal/compiler/builder/directory"
)

func TestBuild(t *testing.T) {
	dir, err := directory.Inspect("testdata/acme")
	if err != nil {
		t.Fatal(fmt.Errorf("prep, directory.Inspect: %w", err))
	}

	tmp, err := os.MkdirTemp(os.TempDir(), "kask-test-build-*")
	if err != nil {
		t.Fatal(fmt.Errorf("os.MkdirTemp: %w", err))
	}

	fmt.Println("temp folder:", tmp)

	err = Build(tmp, dir, &UserSettings{Domain: "http://localhost:8080"})
	if err != nil {
		t.Fatal(fmt.Errorf("act, Build: %w", err))
	}
}
