package httpreplay

import (
	"io/ioutil"
	"net/http"
)

type transportHandler func(h transportFunc) transportFunc

type transportFunc func(req *http.Request) (*http.Response, error)

func (f transportFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

var notHandledTransport = transportFunc(func(req *http.Request) (*http.Response, error) {
	return &http.Response{
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Body:       ioutil.NopCloser(nil),
		Header:     make(http.Header),
		StatusCode: 599,
	}, nil
})
