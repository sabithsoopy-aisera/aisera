package aisera

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type offering struct {
	auhtKey string
}

func (o offering) Bots(ctx context.Context) (bots Bots, err error) {
	headers := http.Header{}
	headers.Add("get-over-post", "true")
	headers.Add("Content-Type", "application/json")
	filter := Filter{
		Fields: []string{"id", "name", "domain", "bot_type", "config"},
		SortCriteria: []SortCriteria{
			{
				Field: "name",
			},
		},
	}
	_, err = o.post(ctx, "/aisera/bots", filter.JSONReader(), &bots, headers)
	if err != nil {
		return nil, fmt.Errorf("error getting bots: %w", err)
	}
	return
}

func (o offering) cookie() *http.Cookie {
	return &http.Cookie{
		Name:     aiseraAdminCookie,
		Value:    o.auhtKey,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}
}

func (o offering) post(ctx context.Context, path string, payload io.Reader, response any, httpHeader http.Header) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, MustURL(path).String(), payload)
	if err != nil {
		return 0, fmt.Errorf("error creating request: %w", err)
	}
	for k, v := range httpHeader {
		for i := range v {
			req.Header.Add(k, v[i])
		}
	}
	req.Header.Add(sessionHeader, o.auhtKey)
	return Parse(req, response)
}

func (o offering) get(ctx context.Context, path string, val any) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, MustURL(path).String(), nil)
	if err != nil {
		return 0, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Add(sessionHeader, o.auhtKey)
	return Parse(req, val)
}
