package memory

import (
	"io"
	"testing"

	"go.ufukty.com/kask/internal/assert"
)

func TestDescriptor_ReadAll(t *testing.T) {
	d := New()
	expected := "hello world"
	
	t.Run("write", func(t *testing.T) {
		err := d.WriteFile("a.txt", []byte(expected))
		if err != nil {
			t.Errorf("act: %v", err)
		}
	})

	var got []byte
	t.Run("read", func(t *testing.T) {
		f, err := d.Open("a.txt")
		if err != nil {
			t.Errorf("prep: %v", err)
		}
		got, err = io.ReadAll(f)
		if err != nil {
			t.Fatalf("act, ReadAll: %v", err)
		}
	})

	assert.Results(t, expected, string(got))
}
