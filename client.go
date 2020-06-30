package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
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
	Raw        []byte
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response.StatusCode = resp.StatusCode
	response.Header = resp.Header
	response.Raw = body
	response.Trailer = resp.Trailer

	decoder := json.NewDecoder(bytes.NewBuffer(body))
	if err := decoder.Decode(&response.Document); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) Post(path string, parms url.Values, document core.Document) (*Response, error) {

	var response Response

	u := c.URL
	u.Path = path
	u.RawQuery = parms.Encode()

	j, err := json.MarshalIndent(document, "", "    ")
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response.StatusCode = resp.StatusCode
	response.Header = resp.Header
	response.Raw = body
	response.Trailer = resp.Trailer

	decoder := json.NewDecoder(bytes.NewBuffer(body))
	if err := decoder.Decode(&response.Document); err != nil {
		return nil, err
	}

	return &response, nil
}
