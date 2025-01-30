// Package httpreplay is for stubbing HTTP requests/responses
package httpreplay

import (
	"bufio"
	"bytes"
	"context"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/aereal/go-http-replay/internal/utils"
)

// NewReplayOrFetchTransport returns new http.RoundTripper that replays HTTP response local cache.
// If the cache is not available, do actual request and record the response to local cache.
func NewReplayOrFetchTransport(dataDir string, httpClient *http.Client) http.RoundTripper {
	return newReplayHandler(dataDir)(newFetchHandler(dataDir, httpClient)(notHandledTransport))
}

// NewReplayTransport returns new http.RoundTripper that only replays HTTP response from local cache, do not request actually.
func NewReplayTransport(dataDir string) http.RoundTripper {
	return newReplayHandler(dataDir)(notHandledTransport)
}

func newFetchHandler(dataDir string, httpClient *http.Client) transportHandler {
	return transportHandler(func(next transportFunc) transportFunc {
		return transportFunc(func(req *http.Request) (*http.Response, error) {
			resp, err := httpClient.Do(req)
			if err != nil {
				return next.RoundTrip(withError(req, err))
			}
			defer resp.Body.Close()

			path := utils.BuildReplayFilePath(dataDir, req)
			f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return next.RoundTrip(withError(req, err))
			}
			dump, err := httputil.DumpResponse(resp, true)
			if err != nil {
				return next.RoundTrip(withError(req, err))
			}
			buf := bytes.NewBuffer(dump)
			if _, err := buf.WriteTo(f); err != nil { //nolint:govet
				return next.RoundTrip(withError(req, err))
			}
			return resp, err
		})
	})
}

func newReplayHandler(dataDir string) transportHandler {
	return transportHandler(func(next transportFunc) transportFunc {
		return transportFunc(func(req *http.Request) (*http.Response, error) {
			f, err := os.Open(utils.BuildReplayFilePath(dataDir, req))
			if err != nil {
				return next.RoundTrip(withError(req, err))
			}
			resp, err := http.ReadResponse(bufio.NewReader(f), req)
			if err != nil {
				return next.RoundTrip(withError(req, err))
			}
			return resp, nil
		})
	})
}

type errCtxKey struct{}

var key = errCtxKey{}

func errorFromContext(ctx context.Context) error {
	if err, ok := ctx.Value(key).(error); ok {
		return err
	}
	return nil
}

func withError(r *http.Request, err error) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), key, err))
}
