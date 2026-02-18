package markdown

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"go.ufukty.com/kask/internal/builder/rewriter"
	"go.ufukty.com/kask/pkg/kask"
)

func matcher(tokens ...string) *regexp.Regexp {
	return regexp.MustCompile(strings.Join(tokens, `\s*`))
}

func TestToHtml_content(t *testing.T) {
	r := rewriter.New()
	r.Bank(".assets/img.jpg", "/.assets/img.jpg")
	r.Bank("sibling.md", "/sibling.html")
	p, err := ToHtml("testdata", "page.md", r)
	if err != nil {
		t.Fatal(fmt.Errorf("act, ToHtml: %w", err))
	}
	content := strings.ReplaceAll(p.Content, "\n", " ")

	t.Run("paragraph", func(t *testing.T) {
		pattern := matcher("<p>A paragraph.</p>")
		if !pattern.MatchString(content) {
			t.Error("could not find")
		}
	})

	t.Run("h1 title", func(t *testing.T) {
		pattern := matcher(`<h1 id="a-title-for-a-markdown-doc">A title for a Markdown doc</h1>`)
		if !pattern.MatchString(content) {
			t.Error("could not find")
		}
	})

	t.Run("img", func(t *testing.T) {
		pattern := matcher(
			"<p>",
			`<img src=".assets/img.jpg" alt="an image" />`,
			"</p>",
		)
		if !pattern.MatchString(content) {
			t.Error("could not find")
		}
	})

	t.Run("table", func(t *testing.T) {
		pattern := matcher(
			"<table>",
			"<thead>", "<tr>", "<th>Header</th>", "</tr>", "</thead>",
			"<tbody>", "<tr>", "<td>Cell <code>1/1</code></td>", "</tr>", "</tbody>",
			"</table>",
		)
		if !pattern.MatchString(content) {
			t.Error("could not find")
		}
	})

	t.Run("codefence", func(t *testing.T) {
		pattern := matcher(
			`<pre tabindex="0" class="chroma"><code><span class="line"><span class="cl"><span class="kr">class</span> <span class="nx">lorem</span> <span class="p">{}</span>`,
			`</span></span></code></pre>`,
		)
		if !pattern.MatchString(content) {
			t.Error("could not find")
		}
	})

	t.Run("unordered list", func(t *testing.T) {
		pattern := matcher(
			`<h2 id="an-unordered-list">An unordered list</h2>`,
			"<ul>",
			"<li>lorem", "<ul>", "<li>ipsum</li>", "</ul>", `<pre tabindex="0" class="chroma"><code><span class="line"><span class="cl"><span class="nb">cd</span> bin`, `</span></span></code></pre>`, "</li>",
			"<li>dolor", "<ul>", "<li>sit</li>", "</ul>", "</li>",
			"</ul>",
		)
		if !pattern.MatchString(content) {
			t.Error("could not find")
		}
	})

	t.Run("ordered list", func(t *testing.T) {
		pattern := matcher(
			`<h3 id="an-ordered-list">An ordered list</h3>`,
			"<ol>",
			"<li>", "<p>another day</p>", "</li>",
			"<li>", "<p>another slay</p>", "</li>",
			"<li>", `<p><a href="/sibling.html">and a link</a></p>`, "</li>",
			"</ol>",
		)
		if !pattern.MatchString(content) {
			t.Error("could not find")
		}
	})
}

func printToc(n *kask.MarkdownTocNode) {
	fmt.Printf("%s%s (%s)\n", strings.Repeat("  ", n.Level), n.Title, n.ID)
	for _, c := range n.Children {
		printToc(c)
	}
}

func ExampleToHtml_toc() {
	r := rewriter.New()
	r.Bank(".assets/img.jpg", "/.assets/img.jpg")
	r.Bank("sibling.md", "/sibling.html")
	p, err := ToHtml("testdata", "page.md", r)
	if err != nil {
		panic(fmt.Errorf("act, ToHtml: %w", err))
	}
	printToc(p.Toc)
	// Output:
	// root ()
	//   A title for a Markdown doc (a-title-for-a-markdown-doc)
	//     An unordered list (an-unordered-list)
	//       An ordered list (an-ordered-list)
}
