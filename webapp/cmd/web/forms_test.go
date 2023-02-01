package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Has(t *testing.T) {
	form := NewForm(nil)

	has := form.Has("nil")
	if has {
		t.Error("form show has field when it should not")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	form = NewForm(postedData)

	has = form.Has("a")
	if !has {
		t.Error("form show does not has field when it should")
	}
}

func TestForm_Requireds(t *testing.T) {
	r := httptest.NewRequest("POST", "/apapun", nil)
	form := NewForm(r.PostForm)

	form.Required("a", "b", "c")

	if form.Valid() {
		t.Error("form shows valid when required fields are missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r, _ = http.NewRequest("POST", "/apapun", nil)
	r.PostForm = postedData

	form = NewForm(r.PostForm)
	form.Required("a", "b", "c")

	if !form.Valid() {
		t.Error("form shows does not have required field when it has")
	}
}

func TestForm_Required(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("a", "")

	form := NewForm(postedData)
	form.Required("a")

	if form.Valid() {
		t.Error("form is valid when it should not because a field is required")
	}

	postedData = url.Values{}
	postedData.Add("a", "terisi")
	postedData.Add("b", "isi juga")

	form = NewForm(postedData)
	form.Required("a", "b")

	if !form.Valid() {
		t.Error("Form is invalid when it should valid")
	}
}

func TestForm_Check(t *testing.T) {
	form := NewForm(nil)
	form.Check(false, "password", "password is required")
	if form.Valid() {
		t.Error("form is valid when it should be not valid")
	}
}

func TestForm_ErrorGet(t *testing.T) {
	form := NewForm(nil)
	form.Check(false, "password", "password is required")
	s := form.Errors.Get("password")

	if len(s) == 0 {
		t.Error("should have an error returned from GET, but do not")
	}

	s = form.Errors.Get("hihi")
	if len(s) != 0 {
		t.Error("should not have an error, but got one")
	}
}
