package builder

import (
	"fmt"
	"path/filepath"

	"go.ufukty.com/kask/internal/providers"
)

func serializeAssetsDirs(d *dir2) []string {
	ss := []string{}
	if d.original.Assets {
		ss = append(ss, filepath.Join(d.paths.Dst, ".assets"))
	}
	for _, c := range d.subdirs {
		if cc := serializeAssetsDirs(c); len(cc) > 0 {
			ss = append(ss, cc...)
		}
	}
	return ss
}

func (b *builder) createDeploymentConfiguration(root *dir2) error {
	switch b.args.Provider {
	case ProviderDefault:
		return nil
	case ProviderCloudflareWorkers:
		err := providers.CloudflareWorkers(b.args.Dst, serializeAssetsDirs(root), b.args.Verbose)
		if err != nil {
			return fmt.Errorf("cloudflare workers: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("unknown provider: %T(%q)", b.args.Provider, b.args.Provider)
	}
}
