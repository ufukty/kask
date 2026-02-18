package providers

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func writer(dst string, verbose bool) (io.WriteCloser, error) {
	if verbose {
		fmt.Printf("creating %s\n", dst)
	}
	f, err := os.Create(dst)
	if err != nil {
		return nil, fmt.Errorf("create file: %w", err)
	}
	return f, nil
}

//go:embed files/workers.txt
var configCloudflareWorkers []byte

func cloudflareWorkers(w io.Writer) error {
	_, err := io.Copy(w, bytes.NewReader(configCloudflareWorkers))
	if err != nil {
		return fmt.Errorf("copying: %w", err)
	}
	return nil
}

func CloudflareWorkers(dst string, verbose bool) error {
	wc, err := writer(filepath.Join(dst, "_headers"), verbose)
	if err != nil {
		return fmt.Errorf("creating writer: %w", err)
	}
	defer wc.Close()
	return cloudflareWorkers(wc)
}
