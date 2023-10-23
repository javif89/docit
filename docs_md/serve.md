# The serve command

The serve command will start up a dev server so you can preview your site.

## Basic

```bash
docit serve -s ./docs
```

Assuming ./docs is where your html files are being built to, this will serve your site from there.

## Build and watch

The docit serve command can take the `--build` and `--watch` flags.

If you use `--build` you should pass the `-i` and `-o` flags same as with the build command. That way it knows where to build your site from and where to output it to. The result will be output to the path passed in the `-s` option.

You can also pass the `-t` flag to set the site title.

`--watch` will watch the `-i` path for changes and automatically rebuild your site on updates. No hot reload yet so you'll
have to refresh the page yourself but it's coming soon.
