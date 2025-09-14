# Recursive Templating

Go templates can call themselves inside with a part of data they are invoked with. This functionality enables us to use templates to print structures like breadcrumbs and sitemap.

## Printing Breadcrumbs

**Breadcrumbs** are basically a **traceback** (you might be familiar to this term from DFS) from one of the "leaves" of the sitemap all the way through to the root following each parent in the **ancestry**.

**To print** a breadcrumb section you can start with `TemplateContent.Node` and follow `Node.Parent` values until you hit a `Node` value its `.Parent` is `nil`, which would only happen if that `Node` value is for the root of site (eg. `example.com`)

As the "depth" of any ancestry is not same at each direction, printing a breadcumb needs recursion, as called **recursive templating** in templating context. Recursion, —based on the **termination condition**— can grow to arbitrary levels of depth in any direction of a **tree** based on a boolean expression on the currently visited node keeps evaluating to true.

Here is the full snippet of `"breadcrumb-item"`, that the most of what you need to print a full breadcrumb list:

```go-html-template
{{define "breadcrumb-item"}}
  {{with .Parent}}
    {{template "breadcrumb-item"}}
    <span class="separator"></span>
  {{end}}

  {{if ne .Href ""}}
    <a href="{{.Href}}">{{.Title}}</a>
  {{else}}
    <span>{{.Title}}</span>
  {{end}}
{{end}}
```

There are couple points here might come strange to beginners:

1. 1st line states this snippet is a template **definition**. Without a `"page"` or `"markdown-page"` named template calling this template, it would not be printed on the page.
1. 2nd line acts like a **termination condition** that stops the recursion on the root. In general, termination condition is what stops the recursive function to create an **infinite recursion**. Here, the condition is only protecting us from getting a templating error otherwise would occur inside the next recursion, where `"breacrumb-item"` is invoked with a `nil` value.
1. 3rd line starts the recursion before printing the current `.Node.Title`. This is important to print the breadcrumb items in usual ordering, from root on the left to the current page on the right.
1. 4th line prints a separator, meant to be placed between each parent-child pair. Notice this line is specifically placed inside the `{{with .Parent}}...{{end}}` block. This is only one way of printing `n-1` number of separators in an `n` items ancestry.
1. 7th to 11th lines are just for printing the current `Node` the template is invoked with. The condition block is here to decide which HTML tag to use to represent the current item in document. If the current `Node` has an `.Href` value other than `""`, then it is **visitable** which is best represented with an anchor `<a>` tag to let visitor to click on it. Notice for the other possibility, a `<span>` tag is used without `.Href` value printed.

From your `"page"` or `"markdown-page"`, call the `"breadcrumb-item"` with `.Node`:

```go-html-template
{{define "markdown-page"}}
  <html>
    <head>
    </head>
    <body>
      <nav id="breadcrumbs">
        {{template "breadcrumb-item" .Node}}
      </nav>
    </body>
  </html>
{{end}}
```

> One `.tmpl` file can house more than one template definitions. Kask only looks for the "page" or "markdown-page" to use as an entry point for rendering a webpage. Other templates are still processed regardless of what file they are defined and supplied to the Go template engine.

## Printing Sitemap

The only thing to do for **printing a sitemap** is following regular DFS order and printing the `.Title` (and `.Href` if populated) of each `Node`. Unlike printing breadcrumbs, printing of sitemap starts from the `TemplateContent.Root`.

```go-html-template
{{define "sitemap-item"}}
  {{if ne .Href ""}}
    <a href="{{.Href}}">{{.Title}}</a>
  {{else}}
    <span>{{.Title}}</span>
  {{end}}

  {{with .Children}}
    <ul>
      {{range .}}
        <li>{{template "sitemap-item" .}}</li>
      {{end}}
    </ul>
  {{end}}
{{end}}
```

Notice that this time, we make the recursive call in 11th line after we print the current `Node` in lines [2-6]. This is because sitemaps usually state the parent then its children.

Don't hesitate to tweak locations of `<ul>` and `<li>` tags. There are many ways of printing a nested list structure with DFS.

Trigger the placement of sitemap by adding a `"sitemap-item"` call with the `TemplateContent.Root` passed in; just as in below:

```go-html-template
{{define "markdown-page"}}
  <html>
    <head>
    </head>
    <body>
      <nav id="sitemap">
        {{template "sitemap-item" .Root}}
      </nav>
    </body>
  </html>
{{end}}
```

## Summary

🥳 🎉 👏 If this is the first time you hear **recursion** and understand how it works here, congrats. It is not an easy concept to understand quickly. Note that there are many issues occur when recursion is applied wrongly. To fix those issues you can start a [Q&A in our discussions](https://github.com/ufukty/kask/discussions/categories/q-a).
