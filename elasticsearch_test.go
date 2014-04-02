package goelasticsearch

import (
	"net/http"
	"runtime"
	"testing"
)

func init() {
	c := NewClient("http://localhost:9200")
	c.DeleteIndex("goelasticsearch", nil)
}

func TestClientDeleteIndex(t *testing.T) {
	// Case when http.NewRequest returns an error.
	c := NewClient("")
	msg := ""
	_, err := c.DeleteIndex("goelasticsearch", nil)
	if err != nil {
		msg = err.Error()
	}
	compare("error", `unsupported protocol scheme ""`, msg, t)

	// Case when client.Do returns an error.
	c = NewClient("http://localhost:65536")
	msg = ""
	_, err = c.DeleteIndex("goelasticsearch", nil)
	if err != nil {
		msg = err.Error()
	}
	compare("error", "dial tcp: invalid port 65536", msg, t)

	// Case when no error is returned.
	c = NewClient("http://localhost:9200")
	body := map[string]interface{}{}
	_, err = c.Create("goelasticsearch", "article", map[string]string{"a": "b"}, nil)
	checkError(err, t)
	msg = ""
	code, err := c.DeleteIndex("goelasticsearch", &body)
	checkError(err, t)
	compare("code", http.StatusOK, code, t)
	expectedBody := map[string]interface{}{"acknowledged": true}
	if body["acknowledged"] != expectedBody["acknowledged"] {
		t.Errorf("Returned body is invalid. [expteced: %v][actual: %v]", expectedBody, body)
	}
}

func TestClientCreate(t *testing.T) {
	// Case when json.Marshal returns an error.
	c := NewClient("http://localhost:9200")
	msg := ""
	_, err := c.Create("goelasticsearch", "article", make(chan string), nil)
	if err != nil {
		msg = err.Error()
	}
	compare("error", "json: unsupported type: chan string", msg, t)

	// Case when http.Post returns an error.
	c = NewClient("http://localhost:65536")
	msg = ""
	_, err = c.Create("goelasticsearch", "article", "", nil)
	if err != nil {
		msg = err.Error()
	}
	compare("error", "Post http://localhost:65536/goelasticsearch/article: dial tcp: invalid port 65536", msg, t)

	// Case when no error is returned.
	c = NewClient("http://localhost:9200")
	body := map[string]interface{}{}
	code, err := c.Create("goelasticsearch", "article", map[string]string{"a": "b"}, &body)
	checkError(err, t)
	compare("code", http.StatusCreated, code, t)
	expectedBody := map[string]interface{}{"_index": "goelasticsearch", "_type": "article", "_version": float64(1), "created": true, "_id": "Automatically generated ID"}
	if body["_index"] != expectedBody["_index"] || body["_type"] != expectedBody["_type"] || body["_version"] != expectedBody["_version"] || body["created"] != expectedBody["created"] {
		t.Errorf("Returned body is invalid. [expected: %s][actual: %s]", expectedBody, body)
	}
}

func TestClientGet(t *testing.T) {
	// Case when http.Get returns an error.
	c := NewClient("http://localhost:65536")
	msg := ""
	_, err := c.Get("goelasticsearch", "article", "x", nil)
	if err != nil {
		msg = err.Error()
	}
	compare("error", "Get http://localhost:65536/goelasticsearch/article/x: dial tcp: invalid port 65536", msg, t)

	// Case when no error is returned.
	c = NewClient("http://localhost:9200")
	body := map[string]interface{}{}
	c.Create("goelasticsearch", "article", map[string]string{"a": "b"}, &body)
	code, err := c.Get("goelasticsearch", "article", body["_id"].(string), &body)
	checkError(err, t)
	compare("code", http.StatusOK, code, t)
	if body["found"] != true {
		t.Errorf("Target document was not found.")
	}
}

func TestClientDelete(t *testing.T) {
	// Case when http.NewRequest returns an error.
	c := NewClient("")
	msg := ""
	_, err := c.Delete("goelasticsearch", "article", "x", nil)
	if err != nil {
		msg = err.Error()
	}
	compare("error", `unsupported protocol scheme ""`, msg, t)

	// Case when client.Do returns an error.
	c = NewClient("http://localhost:65536")
	msg = ""
	_, err = c.Delete("goelasticsearch", "article", "x", nil)
	if err != nil {
		msg = err.Error()
	}
	compare("error", "dial tcp: invalid port 65536", msg, t)

	// Case when no error is returned.
	c = NewClient("http://localhost:9200")
	body := map[string]interface{}{}
	c.Create("goelasticsearch", "article", map[string]string{"a": "b"}, &body)
	code, err := c.Delete("goelasticsearch", "article", body["_id"].(string), &body)
	checkError(err, t)
	compare("code", http.StatusOK, code, t)
	if body["found"] != true {
		t.Errorf("Target document was not found.")
	}
}

func TestClientUpdate(t *testing.T) {
	// Case when json.Marshal returns an error.
	c := NewClient("http://localhost:9200")
	msg := ""
	_, err := c.Update("goelasticsearch", "article", "x", make(chan string), nil)
	if err != nil {
		msg = err.Error()
	}
	compare("error", "json: unsupported type: chan string", msg, t)

	// Case when http.Post returns an error.
	c = NewClient("http://localhost:65536")
	msg = ""
	_, err = c.Update("goelasticsearch", "article", "x", "", nil)
	if err != nil {
		msg = err.Error()
	}
	compare("error", "Post http://localhost:65536/goelasticsearch/article/x/_update: dial tcp: invalid port 65536", msg, t)

	// Case when no error is returned.
	c = NewClient("http://localhost:9200")
	body := map[string]interface{}{}
	c.Create("goelasticsearch", "article", map[string]string{"a": "b", "c": "d"}, &body)
	code, err := c.Update("goelasticsearch", "article", body["_id"].(string), map[string]string{"script": `ctx._source.a = "bb"`}, &body)
	checkError(err, t)
	compare("code", http.StatusOK, code, t)
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
