package builder

import (
	"fmt"

	"go.ufukty.com/kask/internal/providers"
)

func (b *builder) createDeploymentConfiguration() error {
	switch b.args.Provider {
	case ProviderDefault:
		return nil
	case ProviderCloudflareWorkers:
		err := providers.CloudflareWorkers(b.args.Dst, b.args.Verbose)
		if err != nil {
			return fmt.Errorf("cloudflare workers: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("unknown provider: %T(%q)", b.args.Provider, b.args.Provider)
	}
}
