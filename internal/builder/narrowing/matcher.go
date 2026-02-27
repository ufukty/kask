package narrowing

import (
	"fmt"
	"regexp"
)

type Range struct {
	Start, End int
}

func (r Range) String() string {
	return fmt.Sprintf("[%d:%d]", r.Start, r.End)
}

// returns each occurrence as [Range] pairs
func findAllMatches(re *regexp.Regexp, b []byte) []Range {
	ms := re.FindAllIndex(b, -1)
	rs := []Range{}
	for _, m := range ms {
		if len(m) >= 2 {
			rs = append(rs, Range{Start: m[0], End: m[1]})
		}
	}
	return rs
}

// flattens each occurrence's all submatches
func findAllSubmatches(re *regexp.Regexp, b []byte) []Range {
	ms := re.FindAllSubmatchIndex(b, -1)
	rs := []Range{}
	for _, m := range ms {
		for i := 2; i < len(m); i += 2 { // exclude the first
			rs = append(rs, Range{Start: m[i], End: m[i+1]})
		}
	}
	return rs
}

func findAll(re *regexp.Regexp, b []byte) []Range {
	if re.NumSubexp() > 0 {
		return findAllSubmatches(re, b)
	} else {
		return findAllMatches(re, b)
	}
}

// Matchers is to search narrower scopes of a text at each pattern.
// It spans multiple searches for the next ring on each capture group.
type Matchers []*regexp.Regexp

func (chain Matchers) FindAll(b []byte) []Range {
	if len(chain) == 0 {
		return []Range{}
	}
	prev := []Range{{0, len(b)}}
	for _, ring := range chain {
		next := []Range{}
		for _, scope := range prev {
			for _, m := range findAll(ring, b[scope.Start:scope.End]) {
				next = append(next, Range{scope.Start + m.Start, scope.Start + m.End})
			}
		}
		prev = next
	}
	return prev
}

func MustCompile(patterns ...string) Matchers {
	s := Matchers{}
	for _, p := range patterns {
		s = append(s, regexp.MustCompile(p))
	}
	return s
}
