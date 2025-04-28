# Writing Page Template to Render Markdown Pages

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

-   `.Content`: This is the HTML equiavalent of the content of `*.md` file that is the `"markdown-page"` is currently called for. You are expected to write the content of this field to an appropriate place in your template.
-   `.Toc`: This is a tree of nodes, representing "Table of Contents". Note that, since a markdown page is okay to have more than one `H1` title, the TOC starts with a dummy root, representing a H0 with no renderable title. You only need to iterate on its children to access H1s. Printing a Toc require recursive templating.
