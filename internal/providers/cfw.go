package providers

import (
	_ "embed"
	"fmt"
	"io"
	"text/template"

	"go.ufukty.com/kask/internal/disk"
)

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

func CloudflareWorkers(dst disk.WriteFS, assetDirs []string, verbose bool) error {
	if verbose {
		fmt.Printf("creating %s\n", dst)
	}
	wc, err := dst.Create("_headers")
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	defer wc.Close()
	return cloudflareWorkers(wc, assetDirs)
}
