package narrowing

import (
	"regexp"
	"testing"
)

func TestFindAllMatches_noSubExpressions(t *testing.T) {
	input := []byte(`<p><a href="anchor-target"></a></p>`)
	expected := Range{3, 31}
	got := findAll(regexp.MustCompile(`<a[^>]*>[^<]*</a>`), input)
	if len(got) != 1 {
		t.Fatalf("assert, length: expected %d, got %d", 1, len(got))
	}
	if expected != got[0] {
		t.Errorf("assert, content: expected %q, got %q", expected, got)
	}
}

func TestFindAllMatches_subExpressions(t *testing.T) {
	input := []byte(`img-source-2x 2x, img-source-3x 3x, img-source-wide 1000w`)
	expected := []Range{
		{0, 13},  // for "img-source-2x"
		{18, 31}, // for "img-source-3x"
		{36, 51}, // for "img-source-wide"
	}
	got := findAll(regexp.MustCompile(`([^\s]+)\s+\d+(?:\.\d+)?[wx]`), input)
	if len(got) != len(expected) {
		t.Fatalf("assert, length: expected %d, got %d", len(expected), len(got))
	}
	for i := range len(got) {
		if expected[i] != got[i] {
			t.Errorf("assert, item: expected %s, got %s", expected[i], got[i])
		}
	}
}

// this is a snapshot from the builder package, in case if it diverge later
func TestMatchers_FindAll(t *testing.T) {
	example := []byte(`<a href="anchor-target">Lorem ipsum dolor sit amet.</a><img src="img-source" srcset="img-source-set-2x 2x, img-source-set-3x 3x, img-source-set-wide 1000w">`)
	type tc struct {
		matcher  Matchers
		expected []Range
	}
	tcs := map[string]tc{
		"a-href":     {MustCompile(`<a[^>]*>[^<]*</a>`, `href="([^"]*)"`), []Range{{9, 22}}},
		"img-src":    {MustCompile(`<img[^>]*/?>`, `src="([^"]*)"`), []Range{{65, 75}}},
		"img-srcset": {MustCompile(`<img[^>]*/?>`, `srcset="\s*([^"]*)\s*"`, `([^\s]+)\s+\d+(?:\.\d+)?[wx]`), []Range{{85, 102}, {107, 124}, {129, 148}}},
	}
	for tn, tc := range tcs {
		t.Run(tn, func(t *testing.T) {
			got := tc.matcher.FindAll(example)
			if len(tc.expected) != len(got) {
				t.Fatalf("assert, length: expected %d, got %d", len(tc.expected), len(got))
			}
			for i := range len(got) {
				if tc.expected[i] != got[i] {
					t.Errorf("assert, item %d: expected %s, got %s", i, tc.expected[i], got[i])
				}
			}
		})
	}
}
