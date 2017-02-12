# gogenstatic

Command gogenstatic generates "static" package, which embeds all the files from public root dir ("--src" param, defaults to "./public") and provides HTTP handler to serve them.

# Usage

```bash
  gogenstatic --src=path/to/public/dir --dst=path/to/place/static/package
```

It is intended to be used with [go generate](https://blog.golang.org/generate).

Assuming, you have this project structure:

```
  project/
    public/ # contains some HTML/JS/CSS files
    whatever.go
```

Then you can add this comment:

```go
  //go generate gogenstatic
```

to any of your `*.go` files (like `whatever.go`) and run:

```bash
  go generate ./...
```

in your project dir, and it'll generate `static` subpackage, with all the files from `./public/` embedded:

```
  project/
    public/
    static/ # this will be generated from public/
    whatever.go # contains comment //go generate gogenstatic
```

Then import generated subpackage and use it's handler:

```go
  import "path/to/your/project/static"
  http.Handle("/mountpoint/", http.StripPrefix("/mountpoint/", static.Handler()))
```

Or, if you want no runtime hash lookups, you can use handler for a single file:

```go
  http.Handle("/js/file.js", static.HandlerFor("path/to/file.js"))
```
