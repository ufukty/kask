package memory

import (
	"io"
	"testing"

	"go.ufukty.com/kask/internal/assert"
)

func TestDir_create(t *testing.T) {
	d := Dir{}
	_, err := d.Create("lorem")
	if err != nil {
		t.Fatalf("act, unexpected error: %v", err)
	}
}

func TestDir_mkdirAll(t *testing.T) {
	d := &Dir{}
	err := d.MkdirAll("lorem/ipsum/dolor/sit/amet")
	if err != nil {
		t.Fatalf("act, unexpected error: %v", err)
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

func TestDir_mkdirAll_overwriteAsFile(t *testing.T) {
	d := &Dir{}

	t.Run("create as file", func(t *testing.T) {
		if _, err := d.Create("lorem"); err != nil {
			t.Fatalf("prep, unexpected error: %v", err)
		}
	})

	t.Run("overwrite as dir", func(t *testing.T) {
		if err := d.MkdirAll("lorem"); err == nil {
			t.Fatalf("act, unexpected success.")
		}
	})
}

func TestFile_createWriteRead(t *testing.T) {
	expected := "Consectetur adipiscing elit."
	d := &Dir{}

	var w io.WriteCloser
	t.Run("create", func(t *testing.T) {
		var err error
		if w, err = d.Create("lorem"); err != nil {
			t.Fatalf("prep, unexpected error: %v", err)
		}
	})

	t.Run("write", func(t *testing.T) {
		if _, err := w.Write([]byte(expected)); err != nil {
			t.Fatalf("act, unexpected error: %v", err)
		}
	})

	var f *File
	t.Run("find file", func(t *testing.T) {
		n, ok := (*d)["lorem"]
		if !ok {
			t.Fatalf("prep, node doesn't exist")
		}
		f, ok = n.(*File)
		if !ok {
			t.Fatalf("prep, node is not file")
		}
	})

	t.Run("read", func(t *testing.T) {
		got := string(*f)
		assert.Results(t, expected, got)
	})

	t.Run("write again", func(t *testing.T) {
		if _, err := w.Write([]byte(expected)); err != nil {
			t.Fatalf("act, unexpected error: %v", err)
		}
	})

	t.Run("read again", func(t *testing.T) {
		got := string(*f)
		assert.Results(t, expected+expected, got)
	})
}
