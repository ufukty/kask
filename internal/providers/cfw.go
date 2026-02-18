package providers

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

//go:embed files/workers.txt
var configCfw []byte

func cloudflareWorkers(w io.Writer) error {
	_, err := io.Copy(w, bytes.NewReader(configCfw))
	if err != nil {
		return fmt.Errorf("copying: %w", err)
	}
	return nil
}

func CloudflareWorkers(dst string) error {
	dst = filepath.Join(dst, "_headers")
	f, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()
	return cloudflareWorkers(f)
}
