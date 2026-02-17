package providers

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
)

//go:embed files/workers.txt
var configCfw []byte

func CloudflareWorkers(w io.Writer) error {
	_, err := io.Copy(w, bytes.NewReader(configCfw))
	if err != nil {
		return fmt.Errorf("copying: %w", err)
	}
	return nil
}
