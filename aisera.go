package aisera

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
)

const DefaultBaseURL = "https://demo9.login.aisera.cloud"

type Offerings interface {
	Bots(ctx context.Context) (Bots, error)
}

func URL(path string) (*url.URL, error) {
	baseURL := DefaultBaseURL
	if val, ok := os.LookupEnv("AISERA_BASE_URL"); ok {
		baseURL = val
	}
	return url.Parse(fmt.Sprintf("%s/%s", baseURL, strings.TrimPrefix(path, "/")))
}

func MustURL(path string) *url.URL {
	u, _ := URL(path)
	return u
}
