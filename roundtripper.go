package httpreplay

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
)

func NewReplayOrFetchRoundTripper(dataDir string, httpClient *http.Client) http.RoundTripper {
	return newReplayMiddleware(dataDir)(newFetchMiddleware(dataDir, httpClient)(notHandledTripper))
}

func NewReplayRoundTripper(dataDir string) http.RoundTripper {
	return newReplayMiddleware(dataDir)(notHandledTripper)
}

func newFetchMiddleware(dataDir string, httpClient *http.Client) Middleware {
	return Middleware(func(next RoundTripperFunc) RoundTripperFunc {
		return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			resp, err := httpClient.Do(req)
			if err != nil {
				return next.RoundTrip(req)
			}
			defer resp.Body.Close()

			path := getReplayFilePath(dataDir, req)
			f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return next.RoundTrip(req)
			}
			dump, err := httputil.DumpResponse(resp, true)
			if err != nil {
				return next.RoundTrip(req)
			}
			buf := bytes.NewBuffer(dump)
			if _, err := buf.WriteTo(f); err != nil {
				return next.RoundTrip(req)
			}
			return resp, err
		})
	})
}

func newReplayMiddleware(dataDir string) Middleware {
	return Middleware(func(next RoundTripperFunc) RoundTripperFunc {
		return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			f, err := os.Open(getReplayFilePath(dataDir, req))
			if err != nil {
				return next.RoundTrip(req)
			}
			resp, err := http.ReadResponse(bufio.NewReader(f), nil)
			if err != nil {
				return next.RoundTrip(req)
			}
			return resp, nil
		})
	})
}

func getReplayFilePath(dataDir string, req *http.Request) string {
	baseName := url.QueryEscape(req.URL.String())
	return filepath.Join(dataDir, fmt.Sprintf("%s---%s", req.Method, baseName))
}
