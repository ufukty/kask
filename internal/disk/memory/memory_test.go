package memory

import (
	"io"
	"io/fs"
	"slices"
	"testing"
	"testing/fstest"

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
			t.Fatalf("act, unexpected error: %v", err)
		}
		var ok bool
		if fd, ok = w.(*descriptor); !ok {
			t.Fatal("assert, expected descriptor")
		}
	})

	t.Run("read", func(t *testing.T) {
		got := make([]byte, len(expected))
		if _, err := fd.Read(got); err != nil {
			t.Fatalf("act, unexpected error: %v", err)
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
			t.Fatalf("unexpected error: %v", err)
		}
		assert.Results(t, expected, string(got))
	})

	t.Run("close", func(t *testing.T) {
		if err := fd.Close(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("write after close", func(t *testing.T) {
		_, err := fd.Write([]byte("Don't stop me."))
		if err == nil {
			t.Fatalf("unexpected success: %v", err)
		}
	})
}

func TestDescriptor_statSize(t *testing.T) {
	d := New()
	input1 := []byte("Fusce vel posuare.")
	input2 := []byte("Vitae maximus mi bibendum ac.")

	var fd *descriptor
	t.Run("create and open", func(t *testing.T) {
		wc, err := d.Create("input.txt")
		if err != nil {
			t.Fatalf("act, Open: %v", err)
		}
		_, err = wc.Write(input1)
		if err != nil {
			t.Fatalf("act, Write: %v", err)
		}
		var ok bool
		fd, ok = wc.(*descriptor)
		if !ok {
			t.Fatalf("post, expected writable file descriptor, got %T", wc)
		}
	})

	t.Run("compare the file size", func(t *testing.T) {
		length := int64(len(input1))
		t.Run("from descriptor", func(t *testing.T) {
			fi, err := fd.Stat()
			if err != nil {
				t.Fatalf("prep, Stat: %v", err)
			}
			if fi.Size() != length {
				t.Errorf("assert, length: expected %d, got %d", length, fi.Size())
			}
		})
		t.Run("from directory", func(t *testing.T) {
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
		if _, err := fd.Write(input2); err != nil {
			t.Fatalf("act, Write: %v", err)
		}
	})

	t.Run("compare AGAIN the file size", func(t *testing.T) {
		newLength := int64(len(input1) + len(input2))
		t.Run("from descriptor", func(t *testing.T) {
			fi, err := fd.Stat()
			if err != nil {
				t.Fatalf("prep, Stat: %v", err)
			}
			if fi.Size() != newLength {
				t.Errorf("assert, length: expected %d, got %d", newLength, fi.Size())
			}
		})
		t.Run("from directory", func(t *testing.T) {
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
			t.Fatalf("act, stat: %v", err)
		}
		if d.Name() != "c" {
			t.Errorf("assert, name: expected %q, got %q", "c", d.Name())
		}
	})

	t.Run("absolute", func(t *testing.T) {
		d, err := d.Stat("/a/b/c")
		if err != nil {
			t.Fatalf("act, stat: %v", err)
		}
		if d.Name() != "c" {
			t.Errorf("assert, name: expected %q, got %q", "c", d.Name())
		}
	})
}

func TestDir_fsWalkDir(t *testing.T) {
	d := New()

	t.Run("mkdir", func(t *testing.T) {
		err := d.MkdirAll("lorem/ipsum/dolor/sit/amet")
		if err != nil {
			t.Fatalf("act, unexpected error: %v", err)
		}
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
			t.Fatalf("act, WalkDir: %v", err)
		}
	})

	t.Run("compare", func(t *testing.T) {
		expected := []string{
			".",
			"lorem",
			"lorem/ipsum",
			"lorem/ipsum/dolor",
			"lorem/ipsum/dolor/sit",
			"lorem/ipsum/dolor/sit/amet",
		}
		assert.EachResult(t, expected, got)
	})
}

func TestDir_ReadDir(t *testing.T) {
	var (
		d        = New()
		dirs     = []string{"lorem", "ipsum", "dolor", "sit", "amet"}
		files    = []string{"consectetur", "adisipcing", "elit"}
		expected = slices.Concat(dirs, files)
	)
	slices.Sort(expected)

	t.Run("make files and dirs", func(t *testing.T) {
		for _, name := range dirs {
			if err := d.MkdirAll(name); err != nil {
				t.Fatalf("act, MkdirAll %q: %v", name, err)
			}
		}
		for _, name := range files {
			if _, err := d.Create(name); err != nil {
				t.Fatalf("act, Create %q: %v", name, err)
			}
		}
	})

	t.Run("through Dir.ReadDir", func(t *testing.T) {
		got := []string{}
		t.Run("list", func(t *testing.T) {
			es, err := d.ReadDir(".")
			if err != nil {
				t.Fatalf("act, ReadDir: %v", err)
			}
			for _, e := range es {
				got = append(got, e.Name())
			}
		})
		t.Run("compare", func(t *testing.T) {
			assert.EachResult(t, expected, got)
		})
		t.Run("order", func(t *testing.T) {
			assert.Order(t, expected, got)
		})
	})
	t.Run("through Dir.Open+descriptor.ReadDir(-1)", func(t *testing.T) {
		got := []string{}
		t.Run("list", func(t *testing.T) {
			ds, err := d.Open(".")
			if err != nil {
				t.Fatalf("prep, Open: %v", err)
			}
			dsr, ok := ds.(*descriptor)
			if !ok {
				t.Fatalf("prep, expected Open to return a %q got %T", "descriptor", ds)
			}
			es, err := dsr.ReadDir(-1)
			if err != nil {
				t.Fatalf("act, ReadDir: %v", err)
			}
			for _, e := range es {
				got = append(got, e.Name())
			}
		})
		t.Run("compare", func(t *testing.T) {
			assert.EachResult(t, expected, got)
		})
		t.Run("order", func(t *testing.T) {
			assert.Order(t, expected, got)
		})
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
			t.Fatalf("act: %v", err)
		}
	})

	var got []byte
	t.Run("read", func(t *testing.T) {
		f, err := d.Open("a.txt")
		if err != nil {
			t.Fatalf("prep: %v", err)
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

func TestDescriptor_doubleClose(t *testing.T) {
	var wc io.WriteCloser
	t.Run("create", func(t *testing.T) {
		var err error
		if wc, err = New().Create("input.txt"); err != nil {
			t.Fatalf("act: %v", err)
		}
	})

	t.Run("close1", func(t *testing.T) {
		if err := wc.Close(); err != nil {
			t.Fatalf("act: %v", err)
		}
	})

	t.Run("close2", func(t *testing.T) {
		if err := wc.Close(); err != nil {
			t.Errorf("act: %v", err)
		}
	})
}

// see the package doc for partiallity
func TestDir_partialFSConformance(t *testing.T) {
	d := New()
	if err := d.MkdirAll("a/b"); err != nil {
		t.Fatalf("prep, MkdirAll: %v", err)
	}
	if err := d.WriteFile("a/b/hello.txt", []byte("world")); err != nil {
		t.Fatalf("prep, WriteFile: %v", err)
	}
	if err := d.WriteFile("top.txt", []byte("hi")); err != nil {
		t.Fatalf("prep, WriteFile2: %v", err)
	}
	if err := fstest.TestFS(d, "a/b/hello.txt", "top.txt"); err != nil {
		t.Fatalf("act: %v", err)
	}
}
