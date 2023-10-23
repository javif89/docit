# The build command

## Options

```bash
-i --input | Where your markdown files are. Defaults to the current folder.
```

```bash
-o --output | Folder to output your static site to. Defaults to ./build
```

```bash
-t --title | Your site's title. This will show up on the navigation bar and the html title.
```

## Example

If you had a package called "mypkg" and your markdown files were in `docs_md` you could run:

```bash
docit build -i ./docs_md -o ./docs -t mypkg
```
