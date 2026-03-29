package memory

import (
	"io"
	"io/fs"
	"testing"

	"go.ufukty.com/kask/internal/assert"
)

func TestDir_Create(t *testing.T) {
	d := New()
	_, err := d.Create("lorem")
	if err != nil {
		t.Fatalf("act, unexpected error: %v", err)
	}
}

func TestDir_MkdirAll(t *testing.T) {
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

func TestDir_MkdirAll_overwriteAsFile(t *testing.T) {
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

func TestDescriptor_createWriteRead(t *testing.T) {
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

func TestDescriptor_statSize(t *testing.T) {
	d := New()
	input1 := []byte("Fusce vel posuare.")
	input2 := []byte("Vitae maximus mi bibendum ac.")

	t.Run("create", func(t *testing.T) {
		err := d.WriteFile("input.txt", input1)
		if err != nil {
			t.Fatalf("prep, WriteFile: %v", err)
		}
	})

	t.Run("compare", func(t *testing.T) {
		length := int64(len(input1))
		t.Run("with the size acquired from descriptor", func(t *testing.T) {
			fd, err := d.Open("input.txt")
			if err != nil {
				t.Fatalf("prep, Open: %v", err)
			}
			defer fd.Close()
			fi, err := fd.Stat()
			if err != nil {
				t.Fatalf("prep, Stat: %v", err)
			}
			if fi.Size() != length {
				t.Errorf("assert, length: expected %d, got %d", length, fi.Size())
			}
		})
		t.Run("with the size acquired from directory", func(t *testing.T) {
			fi, err := d.Stat("input.txt")
			if err != nil {
				t.Fatalf("prep, Stat: %v", err)
			}
			if fi.Size() != length {
				t.Errorf("assert, length: expected %d, got %d", length, fi.Size())
			}
		})
	})

	t.Run("write more", func(t *testing.T) {
		f, err := d.Open("input.txt")
		if err != nil {
			t.Fatalf("prep, Open: %v", err)
		}
		defer f.Close()
		fd, ok := f.(*descriptor)
		if !ok {
			t.Fatalf("prep, expected writable file descriptor, got %T", f)
		}
		_, err = fd.Write(input2)
		if err != nil {
			t.Fatalf("act, Write: %v", err)
		}
	})

	t.Run("compare again", func(t *testing.T) {
		newLength := int64(len(input1) + len(input2))
		t.Run("with the size acquired from descriptor", func(t *testing.T) {
			fd, err := d.Open("input.txt")
			if err != nil {
				t.Fatalf("prep, Open: %v", err)
			}
			defer fd.Close()
			fi, err := fd.Stat()
			if err != nil {
				t.Fatalf("prep, Stat: %v", err)
			}
			if fi.Size() != newLength {
				t.Errorf("assert, length: expected %d, got %d", newLength, fi.Size())
			}
		})
		t.Run("with the size acquired from directory", func(t *testing.T) {
			fi, err := d.Stat("input.txt")
			if err != nil {
				t.Fatalf("prep, Stat: %v", err)
			}
			if fi.Size() != newLength {
				t.Errorf("assert, length: expected %d, got %d", newLength, fi.Size())
			}
		})
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

func TestDir_fsWalkDir(t *testing.T) {
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

func TestDir_ReadFile(t *testing.T) {
	d := New()
	expected := []byte("some content")

	t.Run("write", func(t *testing.T) {
		if err := d.WriteFile("text.txt", expected); err != nil {
			t.Fatalf("act: %v", err)
		}
	})

	var got []byte
	t.Run("ReadFile", func(t *testing.T) {
		var err error
		if got, err = d.ReadFile("text.txt"); err != nil {
			t.Fatalf("ReadFile: %v", err)
		}
	})

	assert.Results(t, string(expected), string(got))
}

func TestDir_ReadAll(t *testing.T) {
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

func TestDescriptor_readEmptyFile(t *testing.T) {
	d := New()
	err := d.WriteFile("empty.txt", []byte{})
	if err != nil {
		t.Fatalf("prep, create: %v", err)
	}
	fd, err := d.Open("empty.txt")
	if err != nil {
		t.Fatalf("prep, open: %v", err)
	}
	got, err := io.ReadAll(fd)
	if err != nil {
		t.Fatalf("act, ReadAll: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty, got %d bytes", len(got))
	}
}
