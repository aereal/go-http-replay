package httpreplay

import (
	"bufio"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func NewReplayRoundTripper(dataDir string) http.RoundTripper {
	return newReplayMiddleware(dataDir)(notHandledTripper)
}

func newReplayMiddleware(dataDir string) Middleware {
	return Middleware(func(next RoundTripperFunc) RoundTripperFunc {
		return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			baseName := url.QueryEscape(req.URL.String())
			path := filepath.Join(dataDir, baseName)
			f, err := os.Open(path)
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
