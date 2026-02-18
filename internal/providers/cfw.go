package providers

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"
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

type cloudflareWorkersConfiguration struct {
	AssetDirs []string
}

//go:embed templates/workers.tmpl
var templateCloudflareWorkers string

func cloudflareWorkers(w io.Writer, assetDirs []string) error {
	t, err := template.New("config").Parse(templateCloudflareWorkers)
	if err != nil {
		return fmt.Errorf("parsing the embedded configuration template: %w", err)
	}
	err = t.Execute(w, cloudflareWorkersConfiguration{
		AssetDirs: assetDirs,
	})
	if err != nil {
		return fmt.Errorf("executing the configuration template: %w", err)
	}
	return nil
}

func CloudflareWorkers(dst string, assetDirs []string, verbose bool) error {
	wc, err := writer(filepath.Join(dst, "_headers"), verbose)
	if err != nil {
		return fmt.Errorf("creating writer: %w", err)
	}
	defer wc.Close()
	return cloudflareWorkers(wc, assetDirs)
}
