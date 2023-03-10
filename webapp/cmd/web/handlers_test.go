package main

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func Test_application_handlers(t *testing.T) {
	tests := []struct {
		name                    string
		url                     string
		expectedStatusCode      int
		expectedURL             string
		expectedFirstStatusCode int
	}{
		{"home", "/", http.StatusOK, "/", http.StatusOK},
		{"404", "/fish", http.StatusNotFound, "/fish", http.StatusNotFound},
		{"profile", "/user/profile", http.StatusOK, "/", http.StatusTemporaryRedirect},
	}

	routes := app.routes()

	// create test server
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

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

		if resp.Request.URL.Path != tt.expectedURL {
			t.Errorf("%s: expected final url of %s but got %s", tt.name, tt.expectedURL, resp.Request.URL.Path)
		}

		resp2, _ := client.Get(ts.URL + tt.url)
		if resp2.StatusCode != tt.expectedFirstStatusCode {
			t.Errorf("%s: expected first returned status code to be %d but got %d", tt.name, tt.expectedFirstStatusCode, resp.StatusCode)
		}
	}
}

func TestAppHomeOld(t *testing.T) {
	// create a request
	req, _ := http.NewRequest("GET", "/", nil)

	req = addContextAndSessionToRequest(req, app)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(app.Home)

	handler.ServeHTTP(rr, req)

	// check status code
	if rr.Code != http.StatusOK {
		t.Errorf("TestAppHome returned wrong status code; expected 200 but got %d", rr.Code)
	}

	body, _ := io.ReadAll(rr.Body)
	if !strings.Contains(string(body), `<small>From Session:`) {
		t.Error("did not find correct text in html")
	}
}

func TestAppHome(t *testing.T) {
	var tests = []struct {
		name         string
		putInSession string
		expectedHTML string
	}{
		{"first visit", "", "<small>From Session:"},
		{"second visit", "hello, session!", "<small>From Session: hello, session!"},
	}

	for _, tt := range tests {
		// create a request
		req, _ := http.NewRequest("GET", "/", nil)

		req = addContextAndSessionToRequest(req, app)
		_ = app.Session.Destroy(req.Context())

		if tt.putInSession != "" {
			app.Session.Put(req.Context(), "test", tt.putInSession)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(app.Home)

		handler.ServeHTTP(rr, req)

		// check status code
		if rr.Code != http.StatusOK {
			t.Errorf("TestAppHome returned wrong status code; expected 200 but got %d", rr.Code)
		}

		body, _ := io.ReadAll(rr.Body)
		if !strings.Contains(string(body), tt.expectedHTML) {
			t.Errorf("%s: did not find `%s` in response body", tt.name, tt.expectedHTML)
		}
	}
}

func TestApp_renderWithBadTemplate(t *testing.T) {
	// set templatepath to bad template
	pathToTemplates = "./testdata/"

	req, _ := http.NewRequest("GET", "/", nil)
	req = addContextAndSessionToRequest(req, app)
	rr := httptest.NewRecorder()

	err := app.render(rr, req, "bad.page.gohtml", &TemplateData{})
	if err == nil {
		t.Error("expected error from bad template, but didn't get one")
	}

	pathToTemplates = "./../../templates/" // fix template path error
}

func getCtx(req *http.Request) context.Context {
	ctx := context.WithValue(req.Context(), contextUserKey, "unknown")
	return ctx
}

func addContextAndSessionToRequest(req *http.Request, app application) *http.Request {
	req = req.WithContext(getCtx(req))

	ctx, _ := app.Session.Load(req.Context(), req.Header.Get("X-Session"))

	return req.WithContext(ctx)
}

func Test_app_Login(t *testing.T) {
	var tests = []struct {
		name               string
		postedData         url.Values
		expectedStatusCode int
		expectedLoc        string
	}{
		{
			name: "valid login",
			postedData: url.Values{
				"email":    {"admin@example.com"},
				"password": {"secret"},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLoc:        "/user/profile",
		},
		{
			name:               "invalid form request",
			postedData:         url.Values{},
			expectedStatusCode: http.StatusSeeOther,
			expectedLoc:        "/",
		},
		{
			name: "user is not exists",
			postedData: url.Values{
				"email":    {"admin@notfound.com"},
				"password": {"secret"},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLoc:        "/",
		},
		{
			name: "incorrect password",
			postedData: url.Values{
				"email":    {"admin@example.com"},
				"password": {"secret-invalid"},
			},
			expectedStatusCode: http.StatusSeeOther,
			expectedLoc:        "/",
		},
	}

	for _, tt := range tests {
		req, _ := http.NewRequest("POST", "/login", strings.NewReader(tt.postedData.Encode()))
		req = addContextAndSessionToRequest(req, app)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(app.Login)
		handler.ServeHTTP(rr, req)

		if rr.Code != tt.expectedStatusCode {
			t.Errorf("%s: returned wrong status code; expected %d, but got %d", tt.name, tt.expectedStatusCode, rr.Code)
		}

		actualLoc, err := rr.Result().Location()
		if err == nil {
			if actualLoc.String() != tt.expectedLoc {
				t.Errorf("%s: expected location %s but end up in %s", tt.name, tt.expectedLoc, actualLoc.String())
			}
		} else {
			t.Errorf("%s: not location header set", tt.name)
		}
	}
}
