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

### Titles

Kask render the `"title"` template for each Html based page to acquire the user given title for those, when available. To enable this behavior; define a second template named `"title"` as below inside the files of each desired page:

!---

```go-html-template
{{define "title"}}Life is life, na na nana na.{{end}}
```

!---
Figure: `songs.tmpl`

If you desire, you might reuse the `"title"` template inside your `"page"` template:

!---

```go-html-template
{{define "title"}}Life is life, na na nana na.{{end}}

{{define "page"}}
<html>
<head>
<title>{{template "title"}}</title>
</head>
</html>
{{end}}
```

!---
Figure: `songs.tmpl`

For those pages that doesn't contain a `"title"` named template Kask will derive a title using the filename.

## Reusing templates

It makes sense to define the page layout for multiple pages of a section from one place; even when some pages of the section needs customizations. In those circumstances you can define the page template inside the `.kask` or even `.kask/propagate` folders. Which directs Kask to use those as the default `"page"` template when the subfolders, or files of the folder don't override it.

!---

```go-html-template
{{define "page"}}
<html>
  <head>
    <title>{{template "title" .}}</title>
  </head>
  <body>{{template "page-content" .}}</body>
</html>
{{end}}

{{define "markdown-page"}}
<html>
  <head></head>
  <body>
    <main>{{trustedHtml .Markdown.Content}}</main
    <aside>{{trustedHtml .Markdown.Toc}}</aside>
  </body>
</html>
{{end}}
```

!---
Figure: `.kask/propagate/page.tmpl`

!---

```go-html-template
{{define "title"}}Lorem ipsum dolor sit amet.{{end}}

{{define "page-content"}}
{{.Date}}
{{end}}
```

!---
Figure: `contact.tmpl`

Also, moving the `"page"` template into `.kask/propagate` will make it available for pages inside subfolders.

## Escaping

For any text looks like a code piece; the underlying templating engine, thus Kask, will transform it in order to avoid harmful content end up running on the browser of visitors.

To place specific and trusted code on the page as is; you are presented with a series of options:

| Function name     | Use case                                   |
| ----------------- | ------------------------------------------ |
| `trustedCss`      | Stylesheet, CSS rule production and values |
| `trustedHtml`     | HTML document fragments                    |
| `trustedHtmlAttr` | HTML attributes                            |
| `trustedJs`       | EcmaScript5 expressions                    |
| `trustedJsStr`    | Escaped JS                                 |
| `trustedSrcSet`   | Image srcset value                         |
| `trustedUrl`      | Trusted URLs                               |

Use proper function to bypass code escaping. Just as in the example in Markdown section:

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

See the markdown contents are passed through the `trustedHtml` function. This will allow preserving the rich format of the original document like headings, paragraphs, tables and codefences.

Those functions only enable the use of templating engine provided types [`CSS`](https://pkg.go.dev/html/template#CSS), [`HTML`](https://pkg.go.dev/html/template#HTML), [`HTMLAttr`](https://pkg.go.dev/html/template#HTMLAttr), [`JS`](https://pkg.go.dev/html/template#JS), [`JSStr`](https://pkg.go.dev/html/template#JSStr), [`Srcset`](https://pkg.go.dev/html/template#Srcset) and [`URL`](https://pkg.go.dev/html/template#URL) inside your templates. Visit individual links for details. Mind caps.

> Use the `trusted` utilities only when you trust the document to be free of any harmful content.
