// Package goelasticsearch implements an Elasticsearch client.
package goelasticsearch

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// A Client represents an Elasticsearch client.
type Client struct {
	baseURL string
}

// DeleteIndex calls the delete index API.
func (c *Client) DeleteIndex(indexName string) (int, interface{}, error) {
	req, err := http.NewRequest("DELETE", c.baseURL+"/"+indexName, nil)
	if err != nil {
		return 0, nil, err
	}
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	return parseResponse(res)
}

// Create calls the index operation to create a document.
func (c *Client) Create(indexName string, typeName string, data interface{}) (int, interface{}, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return 0, nil, err
	}
	res, err := http.Post(c.baseURL+"/"+indexName+"/"+typeName, "application/json", bytes.NewReader(b))
	if err != nil {
		return 0, nil, err
	}
	return parseResponse(res)
}

// NewClient generates an Elasticsearch client and return it.
func NewClient(baseURL string) *Client {
	return &Client{baseURL: baseURL}
}

// parseResponse parses an HTTP response and returns the result.
func parseResponse(res *http.Response) (int, interface{}, error) {
	code := res.StatusCode
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return code, nil, err
	}
	bodyJSON := make(map[string]interface{})
	if err := json.Unmarshal(body, &bodyJSON); err != nil {
		return code, nil, err
	}
	return code, bodyJSON, nil
}
