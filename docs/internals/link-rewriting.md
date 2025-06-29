# Link Rewriting

Kask rewrites any link found in markdown files that doesn't start with `http://`, `https://` or `/` to match page targeting links to page URL.

Consider links in a page in the path `/subdir/subsubdir/README.md`:

| URL                          | Rewritten URL                |
| ---------------------------- | ---------------------------- |
| `..`                         | `/subdir`                    |
| `../..`                      | `/`                          |
| `../../`                     | `/`                          |
| `../../a.md`                 | `/a.html`                    |
| `../../README.md`            | `/`                          |
| `../`                        | `/subdir`                    |
| `../a.md`                    | `/subdir/a.html`             |
| `../README.md`               | `/subdir`                    |
| `.`                          | `/subdir/subsubdir`          |
| `./..`                       | `/subdir`                    |
| `./../..`                    | `/`                          |
| `./../../`                   | `/`                          |
| `./../`                      | `/subdir`                    |
| `./../a.md`                  | `/subdir/a.html`             |
| `./../README.md`             | `/subdir`                    |
| `./../../a.md`               | `/a.html`                    |
| `./../../README.md`          | `/`                          |
|                              |                              |
| `a`                          | `/subdir/subsubdir/a`        |
| `./a`                        | `/subdir/subsubdir/a`        |
| `./a.md`                     | `/subdir/subsubdir/a.html`   |
| `./a/b.md`                   | `/subdir/subsubdir/a/b.html` |
| `./a/README.md`              | `/subdir/subsubdir/a`        |
| `a.md`                       | `/subdir/subsubdir/a.html`   |
| `a/b.md`                     | `/subdir/subsubdir/a/b.html` |
| `a/README.md`                | `/subdir/subsubdir/a`        |
|                              |                              |
| `./../subsubdir/a.md`        | `/subdir/subsubdir/a.html`   |
| `./../subsubdir/a/b.md`      | `/subdir/subsubdir/a/b.html` |
| `./../subsubdir/a/README.md` | `/subdir/subsubdir/a`        |
| `../subsubdir/a.md`          | `/subdir/subsubdir/a.html`   |
| `../subsubdir/a/b.md`        | `/subdir/subsubdir/a/b.html` |
| `../subsubdir/a/README.md`   | `/subdir/subsubdir/a`        |
