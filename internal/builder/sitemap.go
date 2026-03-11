package builder

import (
	"encoding/xml"
	"fmt"
	"io"

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
	_, err := dst.Write([]byte(xml.Header))
	if err != nil {
		return fmt.Errorf("writing xml header: %w", err)
	}
	e := xml.NewEncoder(dst)
	e.Indent("", "  ")
	urlset := UrlSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}
	err = e.Encode(urlset)
	if err != nil {
		return fmt.Errorf("encoding: %w", err)
	}
	return nil
}
