# Getting started

The documentation for docit is pretty minimal since it's a minimal tool. In this page we'll cover the basics, and we cover a bit more
of each command on their own pages.

## Commands

There's only 2 commands in docit `build` and `serve`. They will each be discussed in more detail in their own pages.
But to get started quickly all you need is a folder with markdown files. Subfolders don't really matter since it will all
be a flat structure in the end, which means you can organize your files however you want.

### Build

Say you have a folder `docs_md` full of markdown files and you want to output your site to `docs` for use with github pages.

All you have to do is:

```bash
docit -i ./docs_md -o ./docs -t "Name of my package"
```

And there you go! You should now have some html files in ./docs.

### Serve

Do you want to preview your site?

```bash
docit serve -s ./docs
```

This will serve your site in localhost:8000

Go ahead and give it a try!
