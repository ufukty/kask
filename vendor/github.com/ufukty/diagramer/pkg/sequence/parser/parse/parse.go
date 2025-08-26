package parse

import (
	"regexp"

	"github.com/ufukty/diagramer/pkg/sequence/parser/parse/ast"
)

var (
	regexLifeline = regexp.MustCompile(`(participant|actor)\s+(\w+)(?:\s+as\s+(.+))?`)
	regexMessage  = regexp.MustCompile(`([^\s-]+)\s*(?:->>[-+]?)\s*([^\s:]*)(?::\s*(.+))?`)
)

func Lifeline(line string) *ast.Lifeline {
	match := regexLifeline.FindStringSubmatch(line)
	if len(match) < 3 {
		return nil
	}
	p := &ast.Lifeline{
		Type:  match[1],
		Alias: "",
		Name:  match[2],
	}
	if len(match) > 3 {
		p.Alias = match[3]
	}
	return p
}

func Message(line string) *ast.Message {
	match := regexMessage.FindStringSubmatch(line)
	if len(match) < 4 {
		return nil
	}
	m := &ast.Message{
		From:    match[1],
		To:      match[2],
		Content: "",
	}
	if len(match) > 3 {
		m.Content = match[3]
	}
	return m
}
