# Link Rewriting

To enable writers to link pages without minding how they will be translated into the build, Kask implicitly rewrites the found links as to make them point to the pages of files the original link points to.

Relative links have different variations. They can point to files in the parent directory, or in one of the contained directories. Consider the file `/subdir/subsubdir/README.md` has some links to other files in the content directory. This table shows how the links point to subdirectories will be rewritten:

| URL in input files | URL in output files          |
| ------------------ | ---------------------------- |
| `.`                | `/subdir/subsubdir/`         |
| `a`                | `/subdir/subsubdir/a/`       |
| `a.md`             | `/subdir/subsubdir/a.html`   |
| `a/b.md`           | `/subdir/subsubdir/a/b.html` |
| `a/README.md`      | `/subdir/subsubdir/a/`       |

This table shows how the links to parent directory will be resolved:

| URL in input files | URL in output files |
| ------------------ | ------------------- |
| `..`               | `/subdir/`          |
| `../..`            | `/`                 |
| `../../a.md`       | `/a.html`           |
| `../../README.md`  | `/`                 |
| `../a.md`          | `/subdir/a.html`    |
| `../README.md`     | `/subdir/`          |

Kask also removes the redundancies in links, where the path is written with entering and exiting segments of one directory:

| URL in input files         | URL in output files          |
| -------------------------- | ---------------------------- |
| `../subsubdir/a.md`        | `/subdir/subsubdir/a.html`   |
| `../subsubdir/a/b.md`      | `/subdir/subsubdir/a/b.html` |
| `../subsubdir/a/README.md` | `/subdir/subsubdir/a/`       |

Additionally, a leading and trailing `./` wouldn't change the results.

Also, Kask leaves external links as they are. Those start with `http://` or `https://`.
