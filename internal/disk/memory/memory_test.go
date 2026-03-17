package memory

import (
	"testing"

	"go.ufukty.com/kask/internal/assert"
)

func TestDir_mkdirAll(t *testing.T) {
	d := Dir{}
	err := d.MkdirAll("lorem/ipsum/dolor/sit/amet")
	if err != nil {
		t.Errorf("act, unexpected error: %v", err)
	}
	expected := []string{
		".",
		"lorem",
		"lorem/ipsum",
		"lorem/ipsum/dolor",
		"lorem/ipsum/dolor/sit",
		"lorem/ipsum/dolor/sit/amet",
	}
	assert.EachResult(t, expected, find(d))
}
