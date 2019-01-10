// Package httpreplay is for stubbing HTTP requests/responses
package httpreplay

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
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
				return next.RoundTrip(embeddingError(req, err))
			}
			defer resp.Body.Close()

			path := getReplayFilePath(dataDir, req)
			f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return next.RoundTrip(embeddingError(req, err))
			}
			dump, err := httputil.DumpResponse(resp, true)
			if err != nil {
				return next.RoundTrip(embeddingError(req, err))
			}
			buf := bytes.NewBuffer(dump)
			if _, err := buf.WriteTo(f); err != nil {
				return next.RoundTrip(embeddingError(req, err))
			}
			return resp, err
		})
	})
}

func newReplayHandler(dataDir string) transportHandler {
	return transportHandler(func(next transportFunc) transportFunc {
		return transportFunc(func(req *http.Request) (*http.Response, error) {
			f, err := os.Open(getReplayFilePath(dataDir, req))
			if err != nil {
				return next.RoundTrip(embeddingError(req, err))
			}
			resp, err := http.ReadResponse(bufio.NewReader(f), nil)
			if err != nil {
				return next.RoundTrip(embeddingError(req, err))
			}
			return resp, nil
		})
	})
}

func getReplayFilePath(dataDir string, req *http.Request) string {
	baseName := url.QueryEscape(req.URL.String())
	return filepath.Join(dataDir, fmt.Sprintf("%s---%s", req.Method, baseName))
}

type ctxKey string

const lastErrorContextKey = ctxKey("last_error")

func getError(ctx context.Context) error {
	err := ctx.Value(lastErrorContextKey)
	if err != nil {
		return err.(error)
	}
	return nil
}

func GetError(req *http.Request) error {
	return getError(req.Context())
}

func withError(pctx context.Context, err error) context.Context {
	return context.WithValue(pctx, lastErrorContextKey, err)
}

func embeddingError(req *http.Request, err error) *http.Request {
	return req.WithContext(withError(req.Context(), err))
}
