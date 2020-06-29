package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/wamuir/go-jsonapi-core"
)

const mime = `application/vnd.api+json`

type Client struct {
	URL url.URL
}

type Response struct {
	StatusCode int
	Header     http.Header
	Document   core.Document
	Trailer    http.Header
}

func New(u url.URL) Client {

	return Client{URL: u}

}

func (c *Client) Get(path string, parms url.Values) (*Response, error) {

	var response Response

	u := c.URL
	u.Path = path
	u.RawQuery = parms.Encode()

	ctx, canx := context.WithTimeout(context.Background(), 5*time.Second)
	defer canx()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.api+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response.StatusCode = resp.StatusCode
	response.Header = resp.Header
	response.Trailer = resp.Trailer

	if err := json.NewDecoder(resp.Body).Decode(&response.Document); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) Post(path string, parms url.Values, body core.Document) (*Response, error) {

	var response Response

	u := c.URL
	u.Path = path
	u.RawQuery = parms.Encode()

	j, err := json.MarshalIndent(body, "", "    ")
	if err != nil {
		return nil, err
	}

	ctx, canx := context.WithTimeout(context.Background(), 5*time.Second)
	defer canx()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/vnd.api+json")
	req.Header.Set("Accept", "application/vnd.api+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response.StatusCode = resp.StatusCode
	response.Header = resp.Header
	response.Trailer = resp.Trailer

	if err := json.NewDecoder(resp.Body).Decode(&response.Document); err != nil {
		return nil, err
	}

	return &response, nil
}
