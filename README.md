# go-http-replay

Record and replay HTTP response for testing

## Synopsis

Replay HTTP response or fetch from the remote:

```go
import (
	"net/http"
	"testing"

	httpreplay "github.com/aereal/go-http-replay"
)

func Test_http_lib(t *testing.T) {
	httpClient := &http.Client{
		Transport: httpreplay.NewReplayOrRoundTripper("./testdata"),
	}
	// httpClient will behave like the client that created from NewReplayRoundTripper but DO actual request if local cache is missing.
}
```

Only replay HTTP response from cache:

```go
import (
	"net/http"
	"testing"

	httpreplay "github.com/aereal/go-http-replay"
)

func Test_http_lib(t *testing.T) {
	httpClient := &http.Client{
		Transport: httpreplay.NewReplayRoundTripper("./testdata"),
	}
	// httpClient will not do actual request to remote sites but returns the response from local cache files.
}
```

## See also

- https://github.com/vcr/vcr - go-http-replay is heavily inspired from this project

## Author

- aereal
