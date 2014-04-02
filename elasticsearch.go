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
func (c *Client) DeleteIndex(indexName string, result interface{}) (int, error) {
	req, err := http.NewRequest("DELETE", c.baseURL+"/"+indexName, nil)

	if err != nil {
		return 0, err
	}

	client := http.Client{}

	res, err := client.Do(req)

	if err != nil {
		return 0, err
	}

	return parseResponse(res, result)
}

// Create calls the index API.
func (c *Client) Create(indexName string, typeName string, data interface{}, result interface{}) (int, error) {
	b, err := json.Marshal(data)

	if err != nil {
		return 0, err
	}

	res, err := http.Post(c.baseURL+"/"+indexName+"/"+typeName, "application/json", bytes.NewReader(b))

	if err != nil {
		return 0, err
	}

	return parseResponse(res, result)
}

// Get calls the get API.
func (c *Client) Get(indexName string, typeName string, id string, result interface{}) (int, error) {
	res, err := http.Get(c.baseURL + "/" + indexName + "/" + typeName + "/" + id)

	if err != nil {
		return 0, err
	}

	return parseResponse(res, result)
}

// Delete calls the delete API.
func (c *Client) Delete(indexName string, typeName string, id string, result interface{}) (int, error) {
	req, err := http.NewRequest("DELETE", c.baseURL+"/"+indexName+"/"+typeName+"/"+id, nil)

	if err != nil {
		return 0, err
	}

	client := http.Client{}

	res, err := client.Do(req)

	if err != nil {
		return 0, err
	}

	return parseResponse(res, result)
}

// Update calls the update API
func (c *Client) Update(indexName string, typeName string, id string, data interface{}, result interface{}) (int, error) {
	b, err := json.Marshal(data)

	if err != nil {
		return 0, err
	}

	res, err := http.Post(c.baseURL+"/"+indexName+"/"+typeName+"/"+id+"/_update", "application/json", bytes.NewReader(b))

	if err != nil {
		return 0, err
	}

	return parseResponse(res, result)
}

// Search calls the search API.
func (c *Client) Search(indexName string, typeName string, q string, result interface{}) (int, error) {
	url := c.baseURL + "/" + indexName + "/" + typeName + "/_search"

	if q != "" {
		url += "?q=" + q
	}

	res, err := http.Get(url)

	if err != nil {
		return 0, err
	}

	return parseResponse(res, result)
}

// NewClient generates an Elasticsearch client and return it.
func NewClient(baseURL string) *Client {
	return &Client{baseURL: baseURL}
}

// parseResponse parses an HTTP response and returns the result.
func parseResponse(res *http.Response, result interface{}) (int, error) {
	code := res.StatusCode

	if result != nil {
		body, err := ioutil.ReadAll(res.Body)

		res.Body.Close()

		if err != nil {
			return code, err
		}

		if err := json.Unmarshal(body, result); err != nil {
			return code, err
		}
	}

	return code, nil
}
