package httpreplay

import (
	"context"
	"io"
	"net/http"
	"path/filepath"
	"testing"
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
		actual := getReplayFilePath(tc.given.dataDir, tc.given.req)
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

func TestNewReplayHandler(t *testing.T) {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://aereal.org/", nil)
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{
		Transport: NewReplayTransport("./testdata"),
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status code=%d", resp.StatusCode)
	}
	if resp.Request == nil {
		t.Fatalf("incoming request not filled")
	}
	if resp.Request.URL.String() != req.URL.String() {
		t.Errorf("request URI mismatch")
	}
}
