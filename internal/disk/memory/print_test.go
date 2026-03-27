package memory

import (
	"strings"
	"testing"
)

func TestHighlight(t *testing.T) {
	var (
		input    = "lorem/ipsum/dolor/sit/amet"
		expected = "lorem/ipsum/" + red + bold + "dolor" + reset + "/sit/amet"
	)
	got := highlight(strings.Split(input, "/"), 2)
	if expected != got {
		t.Errorf("assert:\nexpected: %s\ngot     : %s", expected, got)
	}
}
