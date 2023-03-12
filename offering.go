package aisera

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

const getOverPost = "get-over-post"

type Offerings interface {

	// Bot operations
	Bots(ctx context.Context, filterCriteria Filter) (Bots, error)
	CreateBot(ctx context.Context, bot Bot) (int, error)
	DeleteBot(ctx context.Context, id int) error
	MapBotToChannel(ctx context.Context, botID int, channelIDs []int) error
	RetrainApp(ctx context.Context, botID int) error

	//Channel operations
	Channels(ctx context.Context, filterCriteria Filter) (Channels, error)
	CreateChannel(ctx context.Context, channel Channel) (int, error)
	DeleteChannel(ctx context.Context, id int) error

	//Executions
	Executions(ctx context.Context, filter Filter) ([]Execution, error)

	// HTTP Calls
	Get(ctx context.Context, path string, val any, httpHeaders http.Header) (int, error)
	Post(ctx context.Context, path string, payload io.Reader, response any, httpHeaders http.Header) (int, error)
	Delete(ctx context.Context, path string, payload io.Reader, response any, httpHeaders http.Header) (int, error)
}

func (o offering) HttpHeaders() http.Header {
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
	o.addHeaders(req, o.HttpHeaders())
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
	statusCode, err := o.Post(ctx, "/aisera/bots", createRequest.JSONReader(), &bot, o.HttpHeaders())
	if err != nil {
		return 0, fmt.Errorf("error creating bots: %w", err)
	}
	if statusCode != http.StatusOK {
		return 0, fmt.Errorf("invalid status code: %d", statusCode)
	}
	return bot.ID, nil
}

func (o offering) Bots(ctx context.Context, filterCriteria Filter) (bots Bots, err error) {
	headers := o.HttpHeaders()
	headers.Add(getOverPost, "true")

	_, err = o.Post(ctx, "/aisera/bots", filterCriteria.JSONReader(), &bots, headers)
	if err != nil {
		return nil, fmt.Errorf("error getting bots: %w", err)
	}
	return
}

func (o offering) Delete(ctx context.Context, path string, payload io.Reader, response any, httpHeaders http.Header) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, MustURL(path).String(), payload)
	if err != nil {
		return 0, fmt.Errorf("error creating request: %w", err)
	}
	o.addHeaders(req, httpHeaders)
	return Parse(req, response)
}

func (o offering) Post(ctx context.Context, path string, payload io.Reader, response any, httpHeaders http.Header) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, MustURL(path).String(), payload)
	if err != nil {
		return 0, fmt.Errorf("error creating request: %w", err)
	}
	o.addHeaders(req, httpHeaders)
	return Parse(req, response)
}

func (o offering) Get(ctx context.Context, path string, val any, httpHeaders http.Header) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, MustURL(path).String(), nil)
	if err != nil {
		return 0, fmt.Errorf("error creating request: %w", err)
	}
	o.addHeaders(req, httpHeaders)
	return Parse(req, val)
}

func (o offering) Channels(ctx context.Context, filterCriteria Filter) (channels Channels, err error) {
	headers := o.HttpHeaders()
	headers.Add("get-over-post", "true")

	_, err = o.Post(ctx, "/aisera/channels", filterCriteria.JSONReader(), &channels, headers)
	if err != nil {
		return nil, fmt.Errorf("error getting bots: %w", err)
	}
	return
}

func (o offering) CreateChannel(ctx context.Context, channel Channel) (int, error) {
	createRequest := CreateEntityRequest{
		Entity: channel,
	}
	statusCode, err := o.Post(ctx, "/aisera/channels", createRequest.JSONReader(), &channel, o.HttpHeaders())
	if err != nil {
		return 0, fmt.Errorf("error creating bots: %w", err)
	}
	if statusCode != http.StatusOK {
		return 0, fmt.Errorf("invalid status code: %d", statusCode)
	}
	return channel.ID, nil
}

func (o offering) DeleteChannel(ctx context.Context, id int) error {
	statusCode, err := o.Delete(ctx, "/aisera/channels", DeleteEntityRequest{
		EntityID: id,
	}.JSONReader(), nil, o.HttpHeaders())
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	if statusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d", statusCode)
	}
	return nil
}

func (o offering) MapBotToChannel(ctx context.Context, botID int, channelIDs []int) error {
	request := MappingRequest{
		AddedEntityIds: channelIDs,
		BotID:          botID,
		EntityType:     "aisera.channel",
	}
	statusCode, err := o.Post(ctx, "aisera/bots/mappings", ToJSONReader(request), nil, o.HttpHeaders())
	if err != nil {
		return fmt.Errorf("error mapping bot: %w", err)
	}
	if statusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d", statusCode)
	}
	return nil
}

func (o offering) RetrainApp(ctx context.Context, botID int) error {
	statusCode, err := o.Post(ctx, "aisera/jobs/userJobs", RetrainRequest{
		JobName: "icm-v2-training", //TODO (sabith): revisit  the enum
		JobParams: JobParams{
			BotID: botID,
		},
	}.JSONReader(), nil, o.HttpHeaders())
	if err != nil {
		return fmt.Errorf("error posting request: %w", err)
	}
	if statusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d", statusCode)
	}
	return nil
}

// Executions
func (o offering) Executions(ctx context.Context, filter Filter) ([]Execution, error) {
	queryOutput := &QueryOutput{}
	httpHeaders := o.HttpHeaders()
	httpHeaders.Add(getOverPost, "true")
	statusCode, err := o.Post(ctx, "aisera/userJobs/executions", filter.JSONReader(), &queryOutput, httpHeaders)
	if err != nil {
		return nil, err
	}
	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code: %d", statusCode)
	}
	return queryOutput.Data.([]Execution), nil
}

func (o offering) addHeaders(req *http.Request, headers http.Header) {
	for k, vs := range headers {
		for i := range vs {
			req.Header.Add(k, vs[i])
		}
	}
	if req.Header.Get(sessionHeader) == "" {
		headers.Add(sessionHeader, o.auhtKey)
	}
}
