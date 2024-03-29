[![Build Status][ci-status-badge]][ci-status]
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![GoDoc][godoc-badge]][godoc]

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
		Transport: httpreplay.NewReplayOrFetchTransport("./testdata", http.DefaultClient),
	}
	// httpClient will behave like the client that created from NewReplayTransport but DO actual request if local cache is missing.
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
		Transport: httpreplay.NewReplayTransport("./testdata"),
	}
	// httpClient will not do actual request to remote sites but returns the response from local cache files.
}
```

## See also

- https://github.com/vcr/vcr - go-http-replay is heavily inspired from this project

## Author

- aereal

[license]: https://github.com/aereal/go-http-replay/blob/main/LICENSE
[godoc]: https://pkg.go.dev/github.com/aereal/go-http-replay
[godoc-badge]: https://pkg.go.dev/badge/aereal/go-http-replay
[ci-status]: https://github.com/aereal/go-http-replay/actions/workflows/CI
[ci-status-badge]: https://github.com/aereal/go-http-replay/workflows/CI/badge.svg
