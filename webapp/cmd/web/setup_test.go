package main

import (
	"os"
	"testing"
)

var app application

func TestMain(m *testing.M) {
	pathToTemplates = "./../../templates/" // fix template path error

	app.Session = getSession()

	os.Exit(m.Run())
}
