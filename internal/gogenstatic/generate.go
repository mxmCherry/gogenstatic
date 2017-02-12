package gogenstatic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// Generate generates static package code from `src` directory (public root).
func Generate(out io.Writer, src string) error {
	const hexDict = "0123456789ABCDEF"
	const hexBytesPerLine = 13 // results in 78 chars per line (+ newline)

	src = filepath.Clean(src)

	buf := make([]byte, hexBytesPerLine) // read buffer

	hexLineBufSize := 3 + hexBytesPerLine*(4+1+1) // 3 tabs for each line + "0x00" + "," + " " or "\n" for each byte
	hexLineBuf := make([]byte, hexLineBufSize)

	jsonBuf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(jsonBuf)

	hexLineBuf[0] = '\t'
	hexLineBuf[1] = '\t'
	hexLineBuf[2] = '\t'
	for i, lastI := 3, hexLineBufSize-1; i <= lastI; {

		hexLineBuf[i] = '0'
		i++

		hexLineBuf[i] = 'x'
		i += 1 + 2 // i++ and skip next 2 hex chars

		hexLineBuf[i] = ','
		i += 2 // i++ and skip next space or newline character (no point of assiging it here, as it may be overriden on hex line generation)
	}

	_, _ = fmt.Fprintf(out, header, time.Now().Unix())

	err := filepath.Walk(src, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filePath != "." && len(filePath) > 0 && filePath[0] == '.' {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.IsDir() {
			return nil
		}

		file, err := os.Open(filePath)
		if err != nil {
			return err
		}

		rel := path.Clean(strings.TrimPrefix(filePath, src))

		_, _ = io.WriteString(out, "\tfiles[")

		_ = encoder.Encode(rel)
		id := jsonBuf.Bytes()
		jsonBuf.Reset()

		_, _ = out.Write(id[0 : len(id)-1])
		_, _ = io.WriteString(out, "] = file{\n")
		_, _ = io.WriteString(out, "\t\tcontents: []byte{\n")

		for {

			n, err := file.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}

			j := 3 // skip 3 tabs
			for i, lastN := 0, n-1; i <= lastN; i++ {
				b := buf[i]

				j += 2 // skip "0x"

				hexLineBuf[j] = hexDict[(b&0xF0)>>4]
				j++

				hexLineBuf[j] = hexDict[b&0x0F]
				j++

				hexLineBuf[j] = ','
				j++

				if i == lastN {
					hexLineBuf[j] = '\n'
				} else {
					hexLineBuf[j] = ' '
				}
				j++

			}

			_, _ = out.Write(hexLineBuf[0:j])
		}

		_, _ = io.WriteString(out, "\t\t},\n") // close content: []byte{}
		_, _ = io.WriteString(out, "\t}\n")    // close file{}

		return nil
	})
	if err != nil {
		return err
	}

	_, _ = io.WriteString(out, footer)

	return nil
}

const header = `// Package static provides HTTP handler for serving embedded static files.
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
var modTime = time.Unix(%d, 0)

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

	http.ServeContent(w, r, path.Base(p), ModTime, buf)
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
`

const footer = `}
`
