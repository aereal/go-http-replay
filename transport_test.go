package httpreplay

import (
	"io"
	"net/http"
	"testing"
)

func Test_getReplayFilePath(t *testing.T) {
	type args struct {
		dataDir string
		req     *http.Request
	}
	testcases := []struct {
		given    args
		expected string
	}{
		{
			given: args{
				dataDir: "testdata",
				req:     newRequest(t, http.MethodGet, "http://example.com/", nil),
			},
			expected: "testdata/GET---http%3A%2F%2Fexample.com%2F",
		},
	}

	for _, tc := range testcases {
		actual := getReplayFilePath(tc.given.dataDir, tc.given.req)
		if actual != tc.expected {
			t.Errorf("expected %q but got %q (%#v)", tc.expected, actual, tc.given)
		}
	}
}

func newRequest(t *testing.T, method string, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	}
	return req
}
