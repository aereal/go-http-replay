package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
)

func BuildReplayFilePath(dataDir string, req *http.Request) string {
	baseName := url.QueryEscape(req.URL.String())
	return filepath.Join(dataDir, fmt.Sprintf("%s---%s", req.Method, baseName))
}
