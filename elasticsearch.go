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
func (c *Client) DeleteIndex(indexName string) (int, map[string]interface{}, error) {
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

// Create calls the index API.
func (c *Client) Create(indexName string, typeName string, data interface{}) (int, map[string]interface{}, error) {
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

// Get calls the get API.
func (c *Client) Get(indexName string, typeName string, id string) (int, map[string]interface{}, error) {
	res, err := http.Get(c.baseURL + "/" + indexName + "/" + typeName + "/" + id)
	if err != nil {
		return 0, nil, err
	}
	return parseResponse(res)
}

// Delete calls the delete API.
func (c *Client) Delete(indexName string, typeName string, id string) (int, map[string]interface{}, error) {
	req, err := http.NewRequest("DELETE", c.baseURL+"/"+indexName+"/"+typeName+"/"+id, nil)
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

// Update calls the update API
func (c *Client) Update(indexName string, typeName string, id string, data interface{}) (int, map[string]interface{}, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return 0, nil, err
	}
	res, err := http.Post(c.baseURL+"/"+indexName+"/"+typeName+"/"+id+"/_update", "application/json", bytes.NewReader(b))
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
func parseResponse(res *http.Response) (int, map[string]interface{}, error) {
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
