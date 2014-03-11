package goelasticsearch

import (
	"fmt"
	"net/http"
	"runtime"
	"testing"
)

func init() {
	c := NewClient("http://localhost:9200")
	c.DeleteIndex("goelasticsearch")
}

func TestClientDeleteIndex(t *testing.T) {
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

	// Case when no error is returned.
	c = NewClient("http://localhost:9200")
	_, _, err = c.Create("goelasticsearch", "article", map[string]string{"a": "b"})
	checkError(err, t)
	msg = ""
	code, b, err := c.DeleteIndex("goelasticsearch")
	checkError(err, t)
	compare("code", http.StatusOK, code, t)
	expectedBody := map[string]interface{}{"acknowledged": true}
	body := b.(map[string]interface{})
	if body["acknowledged"] != expectedBody["acknowledged"] {
		t.Errorf("Returned body is invalid. [expteced: %v][actual: %v]", expectedBody, body)
	}
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
	compare("code", http.StatusCreated, code, t)
	expectedBody := map[string]interface{}{"_index": "goelasticsearch", "_type": "article", "_version": float64(1), "created": true, "_id": "Automatically generated ID"}
	body := b.(map[string]interface{})
	if body["_index"] != expectedBody["_index"] || body["_type"] != expectedBody["_type"] || body["_version"] != expectedBody["_version"] || body["created"] != expectedBody["created"] {
		t.Errorf("Returned body is invalid. [expected: %s][actual: %s]", expectedBody, body)
	}
}

func TestClientGet(t *testing.T) {
	// Case when http.Get returns an error.
	c := NewClient("http://localhost:65536")
	msg := ""
	_, _, err := c.Get("goelasticsearch", "article", "x")
	if err != nil {
		msg = err.Error()
	}
	compare("error", "Get http://localhost:65536/goelasticsearch/article/x: dial tcp: invalid port 65536", msg, t)

	// Case when no error is returned.
	c = NewClient("http://localhost:9200")
	_, b, _ := c.Create("goelasticsearch", "article", map[string]string{"a": "b"})
	body := b.(map[string]interface{})
	code, b, err := c.Get("goelasticsearch", "article", body["_id"].(string))
	checkError(err, t)
	compare("code", http.StatusOK, code, t)
	body = b.(map[string]interface{})
	if body["found"] != true {
		t.Errorf("Target document was not found.")
	}
}

func TestClientDelete(t *testing.T) {
	// Case when http.NewRequest returns an error.
	c := NewClient("")
	msg := ""
	_, _, err := c.Delete("goelasticsearch", "article", "x")
	if err != nil {
		msg = err.Error()
	}
	compare("error", `unsupported protocol scheme ""`, msg, t)

	// Case when client.Do returns an error.
	c = NewClient("http://localhost:65536")
	msg = ""
	_, _, err = c.Delete("goelasticsearch", "article", "x")
	if err != nil {
		msg = err.Error()
	}
	compare("error", "dial tcp: invalid port 65536", msg, t)

	// Case when no error is returned.
	c = NewClient("http://localhost:9200")
	_, b, _ := c.Create("goelasticsearch", "article", map[string]string{"a": "b"})
	body := b.(map[string]interface{})
	code, b, err := c.Delete("goelasticsearch", "article", body["_id"].(string))
	checkError(err, t)
	compare("code", http.StatusOK, code, t)
	body = b.(map[string]interface{})
	if body["found"] != true {
		t.Errorf("Target document was not found.")
	}
}

func TestClientUpdate(t *testing.T) {
	// Case when json.Marshal returns an error.
	c := NewClient("http://localhost:9200")
	msg := ""
	_, _, err := c.Update("goelasticsearch", "article", "x", make(chan string))
	if err != nil {
		msg = err.Error()
	}
	compare("error", "json: unsupported type: chan string", msg, t)

	// Case when http.Post returns an error.
	c = NewClient("http://localhost:65536")
	msg = ""
	_, _, err = c.Update("goelasticsearch", "article", "x", "")
	if err != nil {
		msg = err.Error()
	}
	compare("error", "Post http://localhost:65536/goelasticsearch/article/x/_update: dial tcp: invalid port 65536", msg, t)

	// Case when no error is returned.
	c = NewClient("http://localhost:9200")
	_, b, _ := c.Create("goelasticsearch", "article", map[string]string{"a": "b", "c": "d"})
	body := b.(map[string]interface{})
	code, b, err := c.Update("goelasticsearch", "article", body["_id"].(string), map[string]string{"script": `ctx._source.a = "bb"`})
	fmt.Println(b)
	checkError(err, t)
	compare("code", http.StatusOK, code, t)
	body = b.(map[string]interface{})
	if int(body["_version"].(float64)) != 2 {
		t.Errorf("Target document was not updated.")
	}
}

func compare(name, expected, actual interface{}, t *testing.T) {
	if expected != actual {
		_, _, line, _ := runtime.Caller(1)
		t.Errorf("Returned %s is invalid. [expteced: %v][actual: %v][line: %d]", name, expected, actual, line)
	}
}

func checkError(err error, t *testing.T) {
	if err != nil {
		_, _, line, _ := runtime.Caller(1)
		t.Errorf("Error occurred. [error: %s][line: %d]", err.Error(), line)
	}
}
