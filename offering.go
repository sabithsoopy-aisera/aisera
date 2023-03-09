package aisera

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type Offerings interface {
	Bots(ctx context.Context) (Bots, error)
	CreateBot(ctx context.Context, bot Bot) (int, error)
	DeleteBot(ctx context.Context, id int) error
}

func (o offering) httpHeaders() http.Header {
	headers := http.Header{}
	headers.Add("Content-Type", "application/json")
	headers.Add(sessionHeader, o.auhtKey)
	return headers
}

var _ Offerings = offering{}

type offering struct {
	auhtKey string
}

func (o offering) DeleteBot(ctx context.Context, id int) error {
	statusCode, err := o.delete(ctx, "/aisera/bots", DeleteEntityRequest{
		EntityID: id,
	}.JSONReader(), &map[string]interface{}{}, o.httpHeaders())
	if err != nil {
		return fmt.Errorf("error deleting bot: %w", err)
	}
	if statusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d", statusCode)
	}
	return nil
}

func (o offering) CreateBot(ctx context.Context, bot Bot) (int, error) {
	createRequest := CreateEntityRequest{
		Entity: bot,
	}
	statusCode, err := o.post(ctx, "/aisera/bots", createRequest.JSONReader(), &bot, o.httpHeaders())
	if err != nil {
		return 0, fmt.Errorf("error creating bots: %w", err)
	}
	if statusCode != http.StatusOK {
		return 0, fmt.Errorf("invalid status code: %d", statusCode)
	}
	return bot.ID, nil
}

func (o offering) Bots(ctx context.Context) (bots Bots, err error) {
	headers := o.httpHeaders()
	headers.Add("get-over-post", "true")

	filter := Filter{
		Fields: []string{"id", "name", "domain", "bot_type", "config"},
		SortCriteria: []SortCriteria{
			{Field: "name"},
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

func (o offering) delete(ctx context.Context, path string, payload io.Reader, response any, httpHeaders http.Header) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, MustURL(path).String(), payload)
	if err != nil {
		return 0, fmt.Errorf("error creating request: %w", err)
	}
	addHeaders(req, httpHeaders)
	return Parse(req, response)
}

func (o offering) post(ctx context.Context, path string, payload io.Reader, response any, httpHeaders http.Header) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, MustURL(path).String(), payload)
	if err != nil {
		return 0, fmt.Errorf("error creating request: %w", err)
	}
	addHeaders(req, httpHeaders)
	return Parse(req, response)
}

func addHeaders(req *http.Request, headers http.Header) {
	for k, vs := range headers {
		for i := range vs {
			req.Header.Add(k, vs[i])
		}
	}
}

func (o offering) get(ctx context.Context, path string, val any, httpHeaders http.Header) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, MustURL(path).String(), nil)
	if err != nil {
		return 0, fmt.Errorf("error creating request: %w", err)
	}
	addHeaders(req, httpHeaders)
	return Parse(req, val)
}
