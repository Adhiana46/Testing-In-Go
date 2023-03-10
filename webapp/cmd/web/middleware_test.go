package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"webapp/pkg/data"
)

func Test_application_addIPToContext(t *testing.T) {
	tests := []struct {
		headerName  string
		headerValue string
		addr        string
		emptyAddr   bool
	}{
		{"", "", "", false},
		{"", "", "", true},
		{"X-Forwarded-For", "192.3.2.1", "", false},
		{"", "", "hello:world", false},
	}

	// create a dummy handler, use to check context
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// make sure that the value exists in the context
		val := r.Context().Value(contextUserKey)
		if val == nil {
			t.Error(string(contextUserKey), "not present")
		}

		// make sure we got a string back
		ip, ok := val.(string)
		if !ok {
			t.Errorf("not string")
		}
		t.Log(ip) // TODO: delete
	})

	for _, tt := range tests {
		// create the handler
		handlerToTest := app.addIPToContext(nextHandler)

		req := httptest.NewRequest("GET", "http://testing", nil)

		if tt.emptyAddr {
			req.RemoteAddr = ""
		}

		if len(tt.headerName) > 0 {
			req.Header.Add(tt.headerName, tt.headerValue)
		}

		if len(tt.addr) > 0 {
			req.RemoteAddr = tt.addr
		}

		handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
	}
}

func Test_application_ipFromContext(t *testing.T) {
	// get a context
	ctx := context.Background()

	// put something in the context
	ctx = context.WithValue(ctx, contextUserKey, "0.0.0.0")

	// call the function
	ip := app.ipFromContext(ctx)

	// perform the test
	if ip != "0.0.0.0" {
		t.Errorf("unexpected result from ipFromContext(): %s", ip)
	}
}

func Test_app_auth(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//
	})

	var tests = []struct {
		name   string
		isAuth bool
	}{
		{"logged in", true},
		{"not logged in", false},
	}

	for _, tt := range tests {
		handlerToTest := app.auth(nextHandler)
		req := httptest.NewRequest("GET", "http://testing", nil)
		req = addContextAndSessionToRequest(req, app)
		if tt.isAuth {
			app.Session.Put(req.Context(), "user", data.User{ID: 1})
		}
		rr := httptest.NewRecorder()
		handlerToTest.ServeHTTP(rr, req)

		if tt.isAuth && rr.Code != http.StatusOK {
			t.Errorf("%s: expected status code of 200, but got %d", tt.name, rr.Code)
		}

		if !tt.isAuth && rr.Code != http.StatusTemporaryRedirect {
			t.Errorf("%s: expected status code of 307, but got %d", tt.name, rr.Code)
		}
	}
}
