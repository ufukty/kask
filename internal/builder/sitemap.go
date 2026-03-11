package builder

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"go.ufukty.com/kask/pkg/kask"
)

type Url struct {
	Loc string `xml:"loc"`
}

type UrlSet struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	URLs    []Url    `xml:"url"`
}

func subtree(n *kask.Node) []*kask.Node {
	ns := []*kask.Node{n}
	for _, c := range n.Children {
		cs := subtree(c)
		if len(cs) > 0 {
			ns = append(ns, cs...)
		}
	}
	return ns
}

func writeSitemap(dst io.Writer, root *kask.Node) error {
	urls := []Url{}
	for _, n := range subtree(root) {
		if n.Href != "" {
			urls = append(urls, Url{Loc: n.Href})
		}
	}
	if _, err := dst.Write([]byte(xml.Header)); err != nil {
		return fmt.Errorf("header: %w", err)
	}
	e := xml.NewEncoder(dst)
	e.Indent("", "  ")
	urlset := UrlSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}
	if err := e.Encode(urlset); err != nil {
		return fmt.Errorf("body: %w", err)
	}
	return nil
}

func createSitemap(dst string, root *kask.Node) error {
	s, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("opening: %w", err)
	}
	defer s.Close()
	if err = writeSitemap(s, root); err != nil {
		return fmt.Errorf("writing: %w", err)
	}
	return nil
}

func writeRobots(dst io.Writer, sitemap string) error {
	if _, err := fmt.Fprintf(dst, "sitemap: %s\n", sitemap); err != nil {
		return fmt.Errorf("sitemap: %w", err)
	}
	return nil
}

func createRobotsFile(dst, sitemap string) error {
	r, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("opening: %w", err)
	}
	defer r.Close()
	if err = writeRobots(r, sitemap); err != nil {
		return fmt.Errorf("contents: %w", err)
	}
	return nil
}

func (b *builder) sitemap() error {
	sitemap := filepath.Join(b.args.Dst, "sitemap.xml")
	if err := createSitemap(sitemap, b.root3); err != nil {
		return fmt.Errorf("sitemap: %w", err)
	}

	robots := filepath.Join(b.args.Dst, "robots.txt")
	if err := createRobotsFile(robots, sitemap); err != nil {
		return fmt.Errorf("robots: %w", err)
	}

	return nil
}
