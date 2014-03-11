package goelasticsearch

import (
	"net/http"
	"testing"
)

func init() {
	c := NewClient("http://localhost:9200")
	c.DeleteIndex("goelasticsearch")
}

func TestDeleteIndex(t *testing.T) {
	// Case when http.NewRequest returns an error.
	c := NewClient("")
	msg := ""
	_, _, err := c.DeleteIndex("goelasticsearch")
	if err != nil {
		msg = err.Error()
	}
	compare("error", `unsupported protocol scheme ""`, msg, t)

	// Case when client.Do returns an error.
	c = NewClient("http://localhost:65536")
	msg = ""
	_, _, err = c.DeleteIndex("goelasticsearch")
	if err != nil {
		msg = err.Error()
	}
	compare("error", "dial tcp: invalid port 65536", msg, t)
}

func TestClientCreate(t *testing.T) {
	// Case when json.Marshal returns an error.
	c := NewClient("http://localhost:9200")
	msg := ""
	_, _, err := c.Create("goelasticsearch", "article", make(chan string))
	if err != nil {
		msg = err.Error()
	}
	compare("error", "json: unsupported type: chan string", msg, t)

	// Case when http.Post returns an error.
	c = NewClient("http://localhost:65536")
	msg = ""
	_, _, err = c.Create("goelasticsearch", "article", "")
	if err != nil {
		msg = err.Error()
	}
	compare("error", "Post http://localhost:65536/goelasticsearch/article: dial tcp: invalid port 65536", msg, t)

	// Case when no error is returned.
	c = NewClient("http://localhost:9200")
	code, b, err := c.Create("goelasticsearch", "article", map[string]string{"a": "b"})
	checkError(err, t)
	compare("code", code, http.StatusCreated, t)
	body := b.(map[string]interface{})
	expectedBody := map[string]interface{}{"_index": "goelasticsearch", "_type": "article", "_version": float64(1), "created": true, "_id": "Automatically generated ID"}
	if body["_index"] != expectedBody["_index"] || body["_type"] != expectedBody["_type"] || body["_version"] != expectedBody["_version"] || body["created"] != expectedBody["created"] {
		t.Errorf("Returned body is invalid. [expected: %s][actual: %s]", expectedBody, body)
	}
}

func compare(name, expected, actual interface{}, t *testing.T) {
	if expected != actual {
		t.Errorf("Returned %s is invalid. [expteced: %v][actual: %v]", name, expected, actual)
	}
}

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Errorf("Error occurred. [error: %s]", err.Error())
	}
}
