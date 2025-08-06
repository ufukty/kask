package parser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ufukty/diagramer/pkg/sequence/parser/parse"
	"github.com/ufukty/diagramer/pkg/sequence/parser/parse/ast"
)

func Reader(src io.Reader) (*ast.Diagram, error) {
	diagram := &ast.Diagram{
		Lifelines:  make(map[string]*ast.Lifeline),
		Messages:   []*ast.Message{},
		AutoNumber: false,
	}

	scanner := bufio.NewScanner(src)
	for i := 0; scanner.Scan(); i++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "%%") {
			continue
		}

		switch {
		case line == "sequenceDiagram":
			continue

		case line == "autoNumber":
			diagram.AutoNumber = true

		case strings.HasPrefix(line, "participant") || strings.HasPrefix(line, "actor"):
			if p := parse.Lifeline(line); p != nil {
				diagram.Lifelines[p.Name] = p
			}

		case strings.Contains(line, "->>"):
			if m := parse.Message(line); m != nil {
				diagram.Messages = append(diagram.Messages, m)
			}

		default:
			fmt.Printf("skipping line #%d: %s\n", i, line)
		}
	}

	return diagram, nil
}

func File(path string) (*ast.Diagram, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening: %w", err)
	}
	defer file.Close() //nolint:errcheck
	return Reader(file)
}
