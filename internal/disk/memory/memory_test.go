package memory

import (
	"io"
	"io/fs"
	"testing"

	"go.ufukty.com/kask/internal/assert"
)

func TestDir_create(t *testing.T) {
	d := New()
	_, err := d.Create("lorem")
	if err != nil {
		t.Fatalf("act, unexpected error: %v", err)
	}
}

func TestDir_mkdirAll(t *testing.T) {
	d := New()

	t.Run("relative", func(t *testing.T) {
		err := d.MkdirAll("lorem/ipsum/dolor/sit/amet")
		if err != nil {
			t.Fatalf("act, unexpected error: %v", err)
		}
	})

	t.Run("absolute", func(t *testing.T) {
		err := d.MkdirAll("/consectetur/adipiscing")
		if err != nil {
			t.Fatalf("act, unexpected error: %v", err)
		}
	})

	expected := []string{
		".",
		"consectetur",
		"consectetur/adipiscing",
		"lorem",
		"lorem/ipsum",
		"lorem/ipsum/dolor",
		"lorem/ipsum/dolor/sit",
		"lorem/ipsum/dolor/sit/amet",
	}
	assert.EachResult(t, expected, find(d))
}

func TestDir_mkdirAll_overwriteAsFile(t *testing.T) {
	d := New()

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

func TestFileDescriptor_createWriteRead(t *testing.T) {
	d := New()

	var w io.WriteCloser
	t.Run("create", func(t *testing.T) {
		var err error
		if w, err = d.Create("lorem"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	expected := "Consectetur adipiscing elit."

	t.Run("write", func(t *testing.T) {
		if _, err := w.Write([]byte(expected)); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	var fd *descriptor
	t.Run("find file", func(t *testing.T) {
		w, err := d.Open("lorem")
		if err != nil {
			t.Errorf("act, unexpected error: %v", err)
		}
		var ok bool
		if fd, ok = w.(*descriptor); !ok {
			t.Error("assert, expected descriptor")
		}
	})

	t.Run("read", func(t *testing.T) {
		got := make([]byte, len(expected))
		if _, err := fd.Read(got); err != nil {
			t.Errorf("act, unexpected error: %v", err)
		}
		assert.Results(t, expected, string(got))
	})

	expected = "Nam vulputate lectus ligula."

	t.Run("write again", func(t *testing.T) {
		if _, err := w.Write([]byte(expected)); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("read again", func(t *testing.T) {
		got := make([]byte, len(expected))
		if _, err := fd.Read(got); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		assert.Results(t, expected, string(got))
	})

	t.Run("close", func(t *testing.T) {
		if err := fd.Close(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("write after close", func(t *testing.T) {
		_, err := fd.Write([]byte("Don't stop me."))
		if err == nil {
			t.Errorf("unexpected success: %v", err)
		}
	})
}

func TestDir_Stat(t *testing.T) {
	d := New()

	err := d.MkdirAll("a/b/c/d")
	if err != nil {
		t.Fatalf("prep, mkdir all: %v", err)
	}

	t.Run("relative", func(t *testing.T) {
		d, err := d.Stat("a/b/c")
		if err != nil {
			t.Errorf("act, stat: %v", err)
		}
		if d.Name() != "c" {
			t.Errorf("assert, name: expected %q, got %q", "c", d.Name())
		}
	})

	t.Run("absolute", func(t *testing.T) {
		d, err := d.Stat("/a/b/c")
		if err != nil {
			t.Errorf("act, stat: %v", err)
		}
		if d.Name() != "c" {
			t.Errorf("assert, name: expected %q, got %q", "c", d.Name())
		}
	})
}

func TestDir_fsUtilsInterop(t *testing.T) {
	d := New()

	t.Run("mkdir", func(t *testing.T) {
		t.Run("relative", func(t *testing.T) {
			err := d.MkdirAll("lorem/ipsum/dolor/sit/amet")
			if err != nil {
				t.Fatalf("act, unexpected error: %v", err)
			}
		})
		t.Run("absolute", func(t *testing.T) {
			err := d.MkdirAll("/consectetur/adipiscing")
			if err != nil {
				t.Fatalf("act, unexpected error: %v", err)
			}
		})
	})

	got := []string{}
	t.Run("walk dir", func(t *testing.T) {
		err := fs.WalkDir(d, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			got = append(got, path)
			return nil
		})
		if err != nil {
			t.Errorf("act, WalkDir: %v", err)
		}
	})

	t.Run("compare", func(t *testing.T) {
		expected := []string{
			".",
			"consectetur",
			"consectetur/adipiscing",
			"lorem",
			"lorem/ipsum",
			"lorem/ipsum/dolor",
			"lorem/ipsum/dolor/sit",
			"lorem/ipsum/dolor/sit/amet",
		}
		assert.EachResult(t, expected, got)
	})
}
