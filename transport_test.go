package httpreplay_test

import (
	"net/http"
	"testing"
	"time"

	httpreplay "github.com/aereal/go-http-replay"
)

func TestNewReplayHandler(t *testing.T) {
	t.Parallel()
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://aereal.org/", nil)
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{
		Transport:     httpreplay.NewReplayTransport("./testdata"),
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Second * 30,
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
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
