package bundler

import (
	"fmt"
	"io"
	"os"

	"github.com/google/uuid"
)

func newid() (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("uuid.NewRandom: %w", err)
	}
	return u.String(), nil
}

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
