package requests

import (
	"compress/flate"
	"compress/gzip"
	"fmt"
	"net/http"
)

type DecompressingTransport struct {
	Transport http.RoundTripper
}

func (t *DecompressingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error creating gzip reader: %v", err)
		}
		resp.Body = reader
		resp.Header.Del("Content-Encoding")
		resp.Header.Del("Content-Length")
	case "deflate":
		reader := flate.NewReader(resp.Body)
		resp.Body = reader
		resp.Header.Del("Content-Encoding")
		resp.Header.Del("Content-Length")
	}

	return resp, nil
}
