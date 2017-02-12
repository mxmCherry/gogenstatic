// Command gogenstatic generates "static" package,
// which embeds all the files from public root dir ("--src" param, defaults to "./public")
// and provides HTTP handler to serve them.
//
// Example usage:
//   gogenstatic --src=path/to/public/dir --dst=path/to/place/static/package
//
// It is intended to be used with go generate: https://blog.golang.org/generate
//
// Assuming, you have this project structure:
//   project/
//     public/ # contains some HTML/JS/CSS files
//     whatever.go
//
// Then you can add this comment:
//   //go generate gogenstatic
//
// to any of your *.go files (whatever.go) and run:
//   go generate ./...
//
// in your project dir,
// and it'll generate `static` subpackage,
// with all the files from `./public/` embedded:
//   project/
//     public/
//     static/ # this will be generated from public/
//     whatever.go # contains comment //go generate gogenstatic
//
// Then import generated subpackage and use it's handler:
//   import "path/to/your/project/static"
//   http.Handle("/mountpoint/", http.StripPrefix("/mountpoint/", static.Handler()))
//
// Or, if you want no runtime hash lookups, you can use handler for a single file:
//   http.Handle("/js/file.js", static.HandlerFor("path/to/file.js"))
package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/mxmCherry/gogenstatic/internal/gogenstatic"
)

func run() error {
	const dirPerm = 0755
	const filePerm = 0644

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	src := flag.String("src", filepath.Join(wd, "public"), "public root dir")
	dst := flag.String("dst", filepath.Join(wd, "static"), "destination dir to write static.go file")
	flag.Parse()

	if err = os.MkdirAll(*dst, dirPerm); err != nil {
		return err
	}

	outFile, err := os.OpenFile(
		filepath.Join(*dst, "static.go"),
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		filePerm,
	)
	if err != nil {
		return err
	}
	if err = outFile.Truncate(0); err != nil {
		return err
	}

	return gogenstatic.Generate(outFile, *src)
}

func main() {
	if err := run(); err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}
}
