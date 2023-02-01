package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_application_handlers(t *testing.T) {
	tests := []struct {
		name               string
		url                string
		expectedStatusCode int
	}{
		{"home", "/", http.StatusOK},
		{"404", "/fish", http.StatusNotFound},
	}

	var app application
	routes := app.routes()

	// create test server
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	pathToTemplates = "./../../templates/" // fix template path error

	// loop through tests
	for _, tt := range tests {
		resp, err := ts.Client().Get(ts.URL + tt.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != tt.expectedStatusCode {
			t.Errorf("for %s: expected status %d but got %d", tt.name, tt.expectedStatusCode, resp.StatusCode)
		}
	}
}
