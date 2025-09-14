# For Developers

## Folder structure

Kask creates a structure of sitemap which:

- Visitable directories ends with the directory name,
- Pages end with page name with file extension.

Visitable directories are the ones contain either of `index.tmpl` or `README.md`. Pages are always rendered with ending of `.html` regardles the source file ends with `.tmpl` or `.md`.

### Example

Let's say this the folder structure for pages and directories you supplied to Kask:

```sh
.
└── acme
    ├── career
    │   └── index.html
    ├── docs
    │   ├── README.md
    │   ├── birdseed.md
    │   ├── download.md
    │   └── magnet.md
    ├── index.html
    └── products
        └── index.html
```

Resulting sitemap will be like:

```sh
.
└── acme
    ├── career
    │   └── index.html
    ├── docs
    │   ├── index.html
    │   ├── birdseed.html
    │   ├── download.html
    │   └── magnet.html
    ├── index.html
    └── products
        └── index.html
```

Although the `.Node.Href`s will point to the directory-ending for visitable directories:

```sh
.
└── acme
    ├── career
    ├── docs
    │   ├── birdseed.html
    │   ├── download.html
    │   └── magnet.html
    └── products
```

So, a request to `/acme/career` will make the server respond with content of `/acme/career/index.html`, without redirections totally opaque to user agent. This resuls with lean sitemap with minimum number of items.

Make sure you check `.Node.Href != ""` before printing them within anchor tags inside list items in printing your sitemap. Otherwise they are not visitable.
