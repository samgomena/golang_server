package main

import (
    "fmt"
    "encoding/json"
    "io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/status", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	StatusHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("\"/status\" handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	// 	expected := `{"alive": true}`
	// 	if rr.Body.String() != expected {
	// 		t.Errorf("handler returned unexpected body: got %v want %v",
	// 			rr.Body.String(), expected)
	// 	}
}

func TestNotImplementedHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/thisEndpointDoesNotExist", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	NotImplementedHandler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("\"/NotImplementedHandler\" handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestJwtHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/token", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	JwtHandler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("\"/TestJwtHandler\" handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

    var m map[string]string
	rrBuf, _ := ioutil.ReadAll(rr.Body)
	json.Unmarshal(rrBuf, &m)

	if _, ok := m["token"]; !ok {
	    t.Errorf("\"TestJwtHandler\"handler returned unexpected body: got %v want %v",
			m, "token")
	}
}
