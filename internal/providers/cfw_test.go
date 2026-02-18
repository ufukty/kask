package providers

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestCloudflareWorkers(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	err := cloudflareWorkers(b, []string{".assets", "birds/.assets"})
	if err != nil {
		t.Errorf("act, unexpected error: %v", err)
	}
	s := b.String()
	if !strings.Contains(s, "Cache-Control") {
		t.Errorf("assert, expected to contain %q", "Cache-Control")
	}
}

func ExampleCloudflareWorkers() {
	err := cloudflareWorkers(os.Stdout, []string{".assets", "birds/.assets"})
	if err != nil {
		panic(fmt.Errorf("cloudflareWorkers: %w", err))
	}
	// Output:
	// /.assets/*
	//   Cache-Control: public, max-age=14400, must-revalidate
	//
	// /birds/.assets/*
	//   Cache-Control: public, max-age=14400, must-revalidate
	//
}
