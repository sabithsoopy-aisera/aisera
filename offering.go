package aisera

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type Offerings interface {

	// Bot operations
	Bots(ctx context.Context, filterCriteria Filter) (Bots, error)
	CreateBot(ctx context.Context, bot Bot) (int, error)
	DeleteBot(ctx context.Context, id int) error

	//Channel operations
	Channels(ctx context.Context, filterCriteria Filter) (Channels, error)
	CreateChannel(ctx context.Context, channel Channel) (int, error)
	DeleteChannel(ctx context.Context, id int) error
}

func (o offering) httpHeaders() http.Header {
	headers := http.Header{}
	headers.Add("Content-Type", "application/json")
	headers.Add(sessionHeader, o.auhtKey)
	return headers
}

var _ Offerings = offering{}

type offering struct {
	auhtKey       string
	loginResponse LoginResponse
}

func (o offering) DeleteBot(ctx context.Context, id int) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, MustURL("/aisera/bots").String(), DeleteEntityRequest{
		EntityID: id,
	}.JSONReader())
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	addHeaders(req, o.httpHeaders())
	resp, err := Do(req)
	if err != nil {
		return fmt.Errorf("error deleting bot: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d", resp.StatusCode)
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

func (o offering) Bots(ctx context.Context, filterCriteria Filter) (bots Bots, err error) {
	headers := o.httpHeaders()
	headers.Add("get-over-post", "true")

	_, err = o.post(ctx, "/aisera/bots", filterCriteria.JSONReader(), &bots, headers)
	if err != nil {
		return nil, fmt.Errorf("error getting bots: %w", err)
	}
	return
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

func (o offering) get(ctx context.Context, path string, val any, httpHeaders http.Header) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, MustURL(path).String(), nil)
	if err != nil {
		return 0, fmt.Errorf("error creating request: %w", err)
	}
	addHeaders(req, httpHeaders)
	return Parse(req, val)
}

func (o offering) Channels(ctx context.Context, filterCriteria Filter) (channels Channels, err error) {
	headers := o.httpHeaders()
	headers.Add("get-over-post", "true")

	_, err = o.post(ctx, "/aisera/channels", filterCriteria.JSONReader(), &channels, headers)
	if err != nil {
		return nil, fmt.Errorf("error getting bots: %w", err)
	}
	return
}

func (o offering) CreateChannel(ctx context.Context, channel Channel) (int, error) {
	createRequest := CreateEntityRequest{
		Entity: channel,
	}
	statusCode, err := o.post(ctx, "/aisera/channels", createRequest.JSONReader(), &channel, o.httpHeaders())
	if err != nil {
		return 0, fmt.Errorf("error creating bots: %w", err)
	}
	if statusCode != http.StatusOK {
		return 0, fmt.Errorf("invalid status code: %d", statusCode)
	}
	return channel.ID, nil
}

func (o offering) DeleteChannel(ctx context.Context, id int) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, MustURL("/aisera/channels").String(), DeleteEntityRequest{
		EntityID: id,
	}.JSONReader())
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	addHeaders(req, o.httpHeaders())
	resp, err := Do(req)
	if err != nil {
		return fmt.Errorf("error deleting bot: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}
	return nil
}
