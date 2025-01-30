package httpreplay_test

import (
	"context"
	"net/http"
	"testing"

	httpreplay "github.com/aereal/go-http-replay"
)

func TestNewReplayHandler(t *testing.T) {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://aereal.org/", nil)
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{
		Transport: httpreplay.NewReplayTransport("./testdata"),
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
