package bundle

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func appendfile(dst io.Writer, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
	}
	defer in.Close()
	_, err = io.Copy(dst, in)
	if err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}
	_, err = fmt.Fprintf(dst, "\n\n")
	if err != nil {
		return fmt.Errorf("fmt.Fprintf: %w", err)
	}
	return nil
}

func Files(files []string) (string, error) {
	dst := bytes.NewBuffer([]byte{})
	for _, path := range files {
		if err := appendfile(dst, path); err != nil {
			return "", fmt.Errorf("appending %s: %s", path, err)
		}
	}
	return dst.String(), nil
}
