package httpreplay

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type transportHandler func(h transportFunc) transportFunc

type transportFunc func(req *http.Request) (*http.Response, error)

func (f transportFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

var notHandledTransport = transportFunc(func(req *http.Request) (*http.Response, error) {
	var body io.Reader
	if err := errorFromContext(req.Context()); err != nil {
		body = strings.NewReader(err.Error())
	} else {
		body = strings.NewReader("")
	}
	return &http.Response{
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Body:       ioutil.NopCloser(body),
		Header:     make(http.Header),
		StatusCode: 599,
	}, nil
})
