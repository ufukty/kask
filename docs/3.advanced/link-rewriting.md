# Link Rewriting

Kask rewrites any link found in markdown files that doesn't start with `http://`, `https://` or `/` to match page targeting links to page URL.

Consider links in a page in the path `/subdir/subsubdir/README.md`:

| URL in input files         | URL in output files          |
| -------------------------- | ---------------------------- |
| `a`                        | `/subdir/subsubdir/a/`       |
| `a.md`                     | `/subdir/subsubdir/a.html`   |
| `a/b.md`                   | `/subdir/subsubdir/a/b.html` |
| `a/README.md`              | `/subdir/subsubdir/a/`       |
| `../subsubdir/a.md`        | `/subdir/subsubdir/a.html`   |
| `../subsubdir/a/b.md`      | `/subdir/subsubdir/a/b.html` |
| `../subsubdir/a/README.md` | `/subdir/subsubdir/a/`       |
| `.`                        | `/subdir/subsubdir/`         |
| `..`                       | `/subdir/`                   |
| `../..`                    | `/`                          |
| `../../a.md`               | `/a.html`                    |
| `../../README.md`          | `/`                          |
| `../a.md`                  | `/subdir/a.html`             |
| `../README.md`             | `/subdir/`                   |

Additionally, a leading and trailing `./` wouldn't change the results.
