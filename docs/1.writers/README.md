# For Writers

## Folders

Kask generates the sitemap (the hierarchy of web pages based) on the folder structure of the "content directory". Thus, to arrange pages within sections of website, all the writer need to do is moving files between folders inside the content directory. Take a look at this example:

```
.
├── Home.md
├── About
│   ├── Mission-Vision.md
│   ├── Team
│   │   ├── Leadership
│   │   │   ├── CEO.md
│   │   │   └── CTO.md
│   │   └── Staff.md
│   └── History.md
├── Services
│   ├── Web-Development
│   │   ├── Frontend.md
│   │   └── Backend.md
│   ├── Mobile-Development.md
│   └── Consulting.md
├── Products
│   ├── SaaS-Platform
│   │   ├── Features.md
│   │   ├── Pricing
│   │   │   ├── Free.md
│   │   │   ├── Pro.md
│   │   │   └── Enterprise.md
│   └── Integrations.md
├── Blog
│   ├── Tutorials.md
│   ├── Industry-Insights
│   │   ├── AI-ML
│   │   │   ├── Case-Studies.md
│   │   │   └── Research-Articles.md
│   │   └── Web-Trends.md
│   └── News-Updates.md
├── Resources
│   ├── eBooks.md
│   ├── Whitepapers.md
│   └── Templates.md
├── Contact
│   ├── Locations
│   │   ├── New-York.md
│   │   ├── London.md
│   │   └── Tokyo.md
│   └── Support.md
├── FAQ.md
└── Careers
    └── Open-Roles
        ├── Engineering
        │   ├── Backend-Engineer.md
        │   └── Frontend-Engineer.md
        ├── Marketing.md
        └── Sales.md
```

### Ordering

The sitemap orders folder items alphabetically. To apply custom ordering to pages and folders, just prefix the filenames with numbers. One note; the numbers will also be visible in the page address, but not in the page title.

### Hidden developer files

Content directory may contain some developer-exclusive Kask files inside the content directory at any level of subfolders. Telling which files are writer files is easy. Markdown files are the writer files (those end with `.md`). In writer-developer collaborated projects writers are not supposed to understand the functionality of the content of developer files, writers are just expected to not modify, move or delete developer files without the developer. Also developer files might be designed for the files in specific folder; so moving all writer files to a new folder is cheating.

## Pages

Locate the folder in your content directory you want to create the new page under. Right click to an empty space and click "New File". Exact phrasing might be different depending on the system. Name the file something short and meaningful to the writer; it won't be used to decide on the title of page anyways. Filename should end with a `.md` which marks the file as a Markdown file.
