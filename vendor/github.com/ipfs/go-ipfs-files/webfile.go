package files

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
)

// WebFile is an implementation of File which reads it
// from a Web URL (http). A GET request will be performed
// against the source when calling Read().
type WebFile struct {
	body          io.ReadCloser
	url           *url.URL
	contentLength int64
}

// NewWebFile creates a WebFile with the given URL, which
// will be used to perform the GET request on Read().
func NewWebFile(url *url.URL) *WebFile {
	return &WebFile{
		url: url,
	}
}

// Read reads the File from it's web location. On the first
// call to Read, a GET request will be performed against the
// WebFile's URL, using Go's default HTTP client. Any further
// reads will keep reading from the HTTP Request body.
func (wf *WebFile) Read(b []byte) (int, error) {
	if wf.body == nil {
		resp, err := http.Get(wf.url.String())
		if err != nil {
			return 0, err
		}
		wf.body = resp.Body
		wf.contentLength = resp.ContentLength
	}
	return wf.body.Read(b)
}

// Close closes the WebFile (or the request body).
func (wf *WebFile) Close() error {
	if wf.body == nil {
		return nil
	}
	return wf.body.Close()
}

// TODO: implement
func (wf *WebFile) Seek(offset int64, whence int) (int64, error) {
	return 0, ErrNotSupported
}

func (wf *WebFile) Size() (int64, error) {
	if wf.contentLength < 0 {
		return -1, errors.New("Content-Length hearer was not set")
	}

	return wf.contentLength, nil
}

func (wf *WebFile) AbsPath() string {
	return wf.url.String()
}

func (wf *WebFile) Stat() os.FileInfo {
	return nil
}

var _ File = &WebFile{}
var _ FileInfo = &WebFile{}
