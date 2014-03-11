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
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return res.StatusCode, nil, err
	}
	bodyJSON := make(map[string]interface{})
	if err := json.Unmarshal(body, &bodyJSON); err != nil {
		return res.StatusCode, nil, err
	}
	return res.StatusCode, bodyJSON, nil
}

// NewClient generates an Elasticsearch client and return it.
func NewClient(baseURL string) *Client {
	return &Client{baseURL: baseURL}
}
