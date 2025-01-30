package utils_test

import (
	"context"
	"io"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/aereal/go-http-replay/internal/utils"
)

func Test_getReplayFilePath(t *testing.T) {
	type args struct {
		req     *http.Request
		dataDir string
	}
	testcases := []struct {
		given    args
		expected string
	}{
		{
			given: args{
				dataDir: "testdata",
				req:     newRequest(t, context.Background(), http.MethodGet, "http://example.com/", nil),
			},
			expected: filepath.Join("testdata", "GET---http%3A%2F%2Fexample.com%2F"),
		},
	}

	for _, tc := range testcases {
		actual := utils.BuildReplayFilePath(tc.given.dataDir, tc.given.req)
		if actual != tc.expected {
			t.Errorf("expected %q but got %q (%#v)", tc.expected, actual, tc.given)
		}
	}
}

func newRequest(t *testing.T, ctx context.Context, method string, url string, body io.Reader) *http.Request {
	t.Helper()
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		t.Fatal(err)
	}
	return req
}
