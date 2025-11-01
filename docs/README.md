# Kask Docs

<img src=".assets/card-og.png" style="width:min(100%, 640px);border-radius:8px">

Kask is a site compiler where content is managed within folders mainly in Markdown files. The page and category organization follows the folder hierarchy. Simply moving a Markdown file into a folder is enough to move a page between sections in the website.

## Writer-developer collaboration

In collaboration, the writers are expected to "hand over" the latest copy of "content directory" to the developer to perform compilation of website with Kask and deploying it to servers. If the writer is comfortable, they can also use GitHub or similar platform to "sync" the content directory between each other.

Because of the Kask rules, opening the content directory in a Markdown editor won't show the developer files, as either they are hidden files, or stored inside hidden folders (eg. `.kask`).

## For writers

- Work on the website content without the clutter of developer files, directly in their file system using folders and markdown files.
- Use the whatever Markdown editor they are most comfortable with.

## For web developers

- Easily create websites with both Markdown based pages and custom HTML pages.
- Arrange page sections and shared page components within individual files with Go templates.
- Benefit from out-of-the-box best client side performance with CSS splitting and server-side rendered Markdown content.

## GitHub

Repository is on the [GitHub](https://github.com/ufukty/kask).

## License

See [LICENSE](https://github.com/ufukty/kask/blob/main/LICENSE).
