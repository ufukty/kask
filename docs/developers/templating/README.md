# Templating

## Structure of template files

## Templating content

Kask provides a series of useful, dynamic information to templates at the moment they are opened for rendering to static HTML. The struct provided to template file is called `TemplateContent`, which contains many fields:

-   Node (sitemap item) information:
    -   Title of node, extracted from the source file's name, or the `H1` tag of markdown file if available.
    -   Href for other pages to link.
    -   Children, which is a list of `Node`s.
-   Root information:
    -   The root `Node` of the website, typically `href`s to the `/` of website. This is usefull to start printing a sitemap. Just define a recursive template.
-   List of stylesheets for template to include in `<head>`. See [Hierarchical CSS Splitting](../internals/hierarchical-css-splitting.md)
-   Date in Go `time.Time` type. Useful to print the year to footer.
