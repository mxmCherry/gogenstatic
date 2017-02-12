// Package static provides HTTP handler for serving embedded static files.
//
// WARNING!!! It is an automatically generated file, don't edit it manually!!!
package static

import (
	"bytes"
	"net/http"
	"path"
	"sync"
	"time"
)

// modTime holds file modification (generation) time.
var modTime = time.Unix(0000000000, 0)

// files hold static file index by normalized file path.
var files map[string]file

// readerPool is a pool for *bytes.Reader.
var readerPool = sync.Pool{}

// file holds static file data.
type file struct {
	contents []byte
}

// ServeHTTP serves static file over HTTP.
func (f file) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var buf *bytes.Reader
	if v := readerPool.Get(); v == nil {
		buf = bytes.NewReader(f.contents)
	} else {
		buf = v.(*bytes.Reader)
		buf.Reset(f.contents)
	}
	http.ServeContent(w, r, path.Base(r.URL.Path), modTime, buf)
}

// handler holds common handler, that serves all embedded files.
var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// normalize file path:
	p := path.Join("/", r.URL.Path)

	file, ok := files[p]
	if !ok {
		http.NotFound(w, r)
		return
	}

	var buf *bytes.Reader
	if v := readerPool.Get(); v == nil {
		buf = bytes.NewReader(file.contents)
	} else {
		buf = v.(*bytes.Reader)
		buf.Reset(file.contents)
	}

	http.ServeContent(w, r, path.Base(p), modTime, buf)
})

// Handler returns handler for all embedded files.
func Handler() http.Handler {
	return handler
}

// HandlerFor returns static file handler for specified normalized path.
func HandlerFor(filePath string) http.Handler {
	f, ok := files[path.Join("/", filePath)]
	if !ok {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		})
	}
	return f
}

// init populates file index.
func init() {
	files = map[string]file{}
	files["/css/style.css"] = file{
		contents: []byte{
			0x62, 0x6F, 0x64, 0x79, 0x20, 0x7B, 0x0A, 0x09, 0x62, 0x61, 0x63, 0x6B, 0x67,
			0x72, 0x6F, 0x75, 0x6E, 0x64, 0x2D, 0x63, 0x6F, 0x6C, 0x6F, 0x72, 0x3A, 0x20,
			0x77, 0x68, 0x69, 0x74, 0x65, 0x3B, 0x0A, 0x7D, 0x0A,
		},
	}
	files["/index.html"] = file{
		contents: []byte{
			0x3C, 0x21, 0x44, 0x4F, 0x43, 0x54, 0x59, 0x50, 0x45, 0x20, 0x68, 0x74, 0x6D,
			0x6C, 0x3E, 0x0A, 0x3C, 0x68, 0x74, 0x6D, 0x6C, 0x3E, 0x0A, 0x09, 0x3C, 0x68,
			0x65, 0x61, 0x64, 0x3E, 0x0A, 0x09, 0x09, 0x3C, 0x6C, 0x69, 0x6E, 0x6B, 0x20,
			0x72, 0x65, 0x6C, 0x3D, 0x22, 0x73, 0x74, 0x79, 0x6C, 0x65, 0x73, 0x68, 0x65,
			0x65, 0x74, 0x22, 0x20, 0x68, 0x72, 0x65, 0x66, 0x3D, 0x22, 0x2F, 0x63, 0x73,
			0x73, 0x2F, 0x73, 0x74, 0x79, 0x6C, 0x65, 0x2E, 0x63, 0x73, 0x73, 0x22, 0x3E,
			0x0A, 0x09, 0x3C, 0x2F, 0x68, 0x65, 0x61, 0x64, 0x3E, 0x0A, 0x09, 0x3C, 0x62,
			0x6F, 0x64, 0x79, 0x3E, 0x0A, 0x09, 0x09, 0x3C, 0x73, 0x63, 0x72, 0x69, 0x70,
			0x74, 0x20, 0x73, 0x72, 0x63, 0x3D, 0x22, 0x2F, 0x6A, 0x73, 0x2F, 0x73, 0x63,
			0x72, 0x69, 0x70, 0x74, 0x2E, 0x6A, 0x73, 0x22, 0x3E, 0x3C, 0x2F, 0x73, 0x63,
			0x72, 0x69, 0x70, 0x74, 0x3E, 0x0A, 0x09, 0x3C, 0x2F, 0x62, 0x6F, 0x64, 0x79,
			0x3E, 0x0A, 0x3C, 0x2F, 0x68, 0x74, 0x6D, 0x6C, 0x3E, 0x0A,
		},
	}
	files["/js/script.js"] = file{
		contents: []byte{
			0x28, 0x66, 0x75, 0x6E, 0x63, 0x74, 0x69, 0x6F, 0x6E, 0x28, 0x29, 0x20, 0x7B,
			0x0A, 0x09, 0x27, 0x75, 0x73, 0x65, 0x20, 0x73, 0x74, 0x72, 0x69, 0x63, 0x74,
			0x27, 0x0A, 0x09, 0x63, 0x6F, 0x6E, 0x73, 0x6F, 0x6C, 0x65, 0x2E, 0x6C, 0x6F,
			0x67, 0x28, 0x27, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x27, 0x29, 0x0A, 0x7D,
			0x29, 0x28, 0x29, 0x0A,
		},
	}
}
