# Page templates

## Markdown pages

Kask will use any template with name of `"markdown-page"` to pass HTML translation of markdown files. You can create a `*.tmpl` file to define this template on any of the `.kask/propagate` folders either at the level of your `*.md` file or above, or in the `.kask` folder at the level of it.

```sh
.
└── acme
    ├── career
    │   └── index.html
    ├── docs
    │   ├── .kask
    │   │   ├── ...
    │   │   ├── page.tmpl # the .tmpl file defines "markdown-page"
    │   │   └── ...
    │   ├── README.md   # will use acme/docs/.kask/page.tmpl
    │   ├── birdseed.md # will use acme/docs/.kask/page.tmpl
    │   ├── download.md # will use acme/docs/.kask/page.tmpl
    │   └── magnet.md   # will use acme/docs/.kask/page.tmpl
    ├── index.html
    └── products
        └── index.html
```

Markdown page templates --just like any other `.kask/*.tmpl` file-- are overridable. The more specific template file will be applied instead of the one defined in above directories.

## Structure of Markdown Page Templates

Markdown page templates are provided additional `.Markdown` field on top of [usual contents](../README.md#templating-content) of `TemplateContent`. Markdown field contains two fields:

- `.Content`: This is the HTML equiavalent of the content of `*.md` file that is the `"markdown-page"` is currently called for. You are expected to write the content of this field to an appropriate place in your template.
- `.Toc`: This is a tree of nodes, representing "Table of Contents". Note that, since a markdown page is okay to have more than one `H1` title, the TOC starts with a dummy root, representing a H0 with no renderable title. You only need to iterate on its children to access H1s. Printing a Toc require recursive templating.

```go-html-template
{{define "markdown-page"}}
<html>
  <body>
    <main>{{trustedHtml .Markdown.Content}}</main>
    <aside>{{trustedHtml .Markdown.Toc}}</aside>
  </body>
</html>
{{end}}
```

See the [Escaping](#escaping) section for `trustedHtml`.

## HTML pages

Kask loads all shared template files stored inside the `.kask` folder of containing folder and its parent folders for rendering all `.tmpl` ending files in a content directory.

!---

```go-html-template
{{define "page"}}
<html>
  <head></head>
  <body>{{.Date}}</body>
</html>
{{end}}
```

!---
Figure: `contact.tmpl`
Unlike Markdown based pages; template based pages of the same folder can have different layouts and styling. Just define the `"page"` template at the each page file, where customization is desired.

### Shared page template

Rendering each page takes creating a copy of the shared template stack and adding the contents of template file in rendering step. Thus template pages are not strictly expected to have `{{define "page"}}` block; instead, one directory can have it as a shared template file and page templates can only define the content exclusive to each page.

!---

```go-html-template
{{define "page"}}
<html>
  <head></head>
  <body>{{template "page-content" .}}</body>
</html>
{{end}}
```

!---
Figure: `.kask/page.tmpl`

!---

```go-html-template
{{define "page-content"}}
{{.Date}}
{{end}}
```

!---
Figure: `contact.tmpl`

Also, moving the `"page"` template into `.kask/propagate` will make it available for pages inside subfolders.
