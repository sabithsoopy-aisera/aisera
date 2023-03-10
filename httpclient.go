package aisera

import (
	"net"
	"net/http"
	"time"

	"golang.org/x/exp/slog"
)

var httpClient = http.DefaultClient

func init() {
	httpClient.Transport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
		Proxy:               http.ProxyFromEnvironment,
	}
	httpClient.Timeout = 10 * time.Second
}

func SetHttpClient(client *http.Client) {
	httpClient = client
}

func Do(req *http.Request) (*http.Response, error) {
	logger.Debug("making the request", slog.String("method", req.Method), slog.String("url", req.URL.String()))
	return httpClient.Do(req)
}

func Parse(req *http.Request, val any) (int, error) {
	resp, err := Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if val == nil {
		return resp.StatusCode, nil
	}
	return resp.StatusCode, JSONReaderToVal(resp.Body, val)
}

func addHeaders(req *http.Request, headers http.Header) {
	for k, vs := range headers {
		for i := range vs {
			req.Header.Add(k, vs[i])
		}
	}
}
