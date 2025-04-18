# Kask _(casc'ade)_

Kask builds static websites out of an hierarchy of `.md`, `.html`, `.tmpl` and `.css` files in a way optimized for the fastest page load.

Kask doesn't provide much flexibility and doesn't support any sort of dynamic content to achive simplest usage.

Kask can compile Go templates into pages. Which is useful feature to define common parts of a website like header and footer once and use every page. Go knowledge is not required to create and use templates. See [Templating](#Templating) for details.

Kask converts Markdown files to HTML and puts each Markdown file in to a page template to create final pages. All Markdown rendering done at building phase, hence the fastest page load is achived.

Kask doesn't need the person who create the content of website to be technical. Content creator just needs to provide folders of Markdown files and developer only needs to design one template page in simple HTML and CSS. See [Markdown page template](#Markdown-Page-Template) for details.

## Features

### Optimized to ease populating website content

- Generate Markdown files in your file system to create pages in the website
- Sitemap follows folder hierarchy.
- Standard Markdown features are supported: code blocks, tables, images etc.

### Flexibility

- Either provide one HTML page for the all markdown pages or write any page as HTML document
- Bring your own styles (CSS files)
- Write JS scripts

### Minimal development

- Go templates
- Sitemap reflects the folder hierarchy.
- Reuse of HTML files for components like header, footer etc.
- Provides sitemap data structure to be used in templating.
- Maximum code reuse without sacrificing per-section layout and styling.

### Performance

- Minimum page load with balanced (per-section) css file bundling.
- Builds contain only static files which are better for browser performance.

## Usage

### Kask project layout

Starts as the same with a vanilla HTML-CSS-JS project:

```
. acme co
├── career
│   └── index.html
├── docs
│   ├── birdseed.md
│   ├── magnet.md
│   └── meta.yml
├── index.html
└── products
    └── index.html
```

The only difference is the `.kask` and `.assets` folders that can be placed in any level of directories.

### `.kask` folder

Any folder within the directory contains the website content can also contain a folder named `.kask` which can have below files:

```
.kask
├── propagate
│   ├── *.css
│   └── *.tmpl
├── *.css
├── *.tmpl
├── page.html
└── meta.yml
```

| File/Pattern | Description                                                                                                                                    |
| ------------ | ---------------------------------------------------------------------------------------------------------------------------------------------- |
| `*.css`      | all css files in one `.kask` folder gets bundled into one file. the bundle file gets linked by all pages in the containing directory.          |
| `*.tmpl`     | all tmpl files in one `.kask` folder gets put in the templating process of pages in the containing directory.                                  |
| `page.html`  | full html document that is used for placing compiled each markdown page into. it references `{{ markdown }}` function at the appropriate spot. |
| `meta.yml`   | contains misc. annotation related to containing directory.                                                                                     |
| `.prop/`     | `*.css` and `*.tmpl` files in `.prop` directory are treated like they are copied in every subdir of containing directory of `.kask`            |

As long as they are wanted to be applied site-wide; placing css, tmpl and page.html files in .kask/propagate is better for maximizing code-reuse.

### `.assets` folder

Store images, videos, documents that are linked by html pages or markdown files. Kask will copy the content of this folder into the related folder in the build, without requiring link adjustments.

### Templating

#### Markdown Page Template

Content of a basic Markdown template `acme/products/.kask/page.html`:

```html
<html>
  <head> </head>
  <body>
    <main>{{.MarkdownContent}}</main>
    <aside>{{.MarkdownTOC}}</aside>
  </body>
</html>
```

Kask will compile every Markdown file in the `acme/products` folder to HTML and put each to one copy of `acme/products/.kask/page.html`.

Kask also provides the outline of the Markdown file in `.MarkdownTOC` which reflects the hierarchy of headings in the documents and `#` links to them. It is an HTML code that contains a nested unordered lists that contains `a` tags, inside a `nav` tag.

#### Templating data

- `.Stylesheets` (`[]string)`)  
  Contains the list of `*.css` files, that is needed to be placed in the `<head>` of web page. Since Kask bundles multiple css files into one file per folder and supports propagation there might be more than 1 stylesheet for each page. Such usage is appropriate for most of the time:

  ```go-templates
  {{range .Stylesheets}}
    <link rel="stylesheet" href="{{.}}">{{end}}
  ```

- `.Node` (`*Node`)

- `.WebSiteRoot` (`*Node`)

- `.MarkdownContent` (`string`)

- `.MarkdownTOC` (`string`)

- `.Time` (`time.Time`)

- `.Dir` (`[]string`)

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>{{.Title}}</title>
    <meta name="description" content="" />
    <link rel="stylesheet" href="" />

    {{.Stylesheets}}
  </head>
  <body>
    {{.Sitemap}}
  </body>
</html>
```

### Build a Kask project

```
kask build -in 'path' -out 'path'
```

### Development server

```
kask serve -in 'path' -p 8080
```

Don't use this server in production as it doesn't provide any security measure and scale a web server should do.
