# Hierarchical CSS Splitting

There are two fundamentals on justifying hierarchical CSS splitting:

-   Visitor expects fastest page load even when they enter to website from a page in depths of the whole.
-   Pages belong to a website need styles:
    -   exclusive to that page,
    -   exclusive to a subsection of the website and shared amongst other pages belong to the same subsection and,
    -   shared amongst all pages belong to website.

Hierarchical CSS splitting in Kask "bundles" each `/**/.kask/propagate` folder contents at its own. Resulting the webpages in deeper levels to have bigger number of stylesheets linked starting from the most "generic" to most "specific". The most spefic one is the bundle made out of css files in the `.kask` folder at the same level with the page.

```sh
.
├── .kask
│   ├── propagate
│   │   └── *.css
│   └── *.css
└── docs
    ├── .kask
    │   ├── propagate
    │   │   └── *.css
    │   └── *.css
    ├── tutorials
    │       └── *.md # [.kask/propagate/*.css, docs/.kask/propagate/*.css]
    └── *.md # [.kask/propagate/*.css, docs/.kask/propagate/*.css, docs/.kask/*.css]
```

## Attaching bundles to page

Templates named `"page"` and `"markdown-page"` are expected to include all bundles passed to them inside `<head>` tag:

```go-html-template
{{define "page"}}
<html>
    <head>
        {{range .Stylesheets}}
        <link rel="stylesheet" href="{{.}}" />
        {{end}}
    </head>
</html>
{{end}}
```
