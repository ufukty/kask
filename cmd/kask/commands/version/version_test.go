package version

import (
	"fmt"
	"testing"
)

func TestDigBuildInfo(t *testing.T) {
	v, err := DigBuildInfo()
	if err != nil {
		t.Errorf("act, unexpected error: %v", err)
	}
	if v == "" {
		t.Errorf("assert, unexpected empty value")
	}
	fmt.Println(v)
}
