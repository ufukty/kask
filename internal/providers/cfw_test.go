package providers

import (
	"bytes"
	"strings"
	"testing"
)

func TestCloudflareWorkers(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	err := cloudflareWorkers(b)
	if err != nil {
		t.Errorf("act, unexpected error: %v", err)
	}
	s := b.String()
	if !strings.Contains(s, "Cache-Control") {
		t.Errorf("assert, expected to contain %q", "Cache-Control")
	}
}
