package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func Test_isPrime(t *testing.T) {
	primeTests := []struct {
		name     string
		testNum  int
		expected bool
		msg      string
	}{
		{"prime", 7, true, "7 is a prime number!"},
		{"not prime", 8, false, "8 is not a prime number because it is divisible by 2!"},
		{"zero", 0, false, "0 is not prime, by definition!"},
		{"one", 1, false, "1 is not prime, by definition!"},
		{"negative number", -11, false, "Negative numbers are not prime, by definition!"},
	}

	for _, e := range primeTests {
		result, msg := isPrime(e.testNum)
		if e.expected && !result {
			t.Errorf("%s: expected true but got false", e.name)
		}

		if !e.expected && result {
			t.Errorf("%s: expected false but got true", e.name)
		}

		if e.msg != msg {
			t.Errorf("%s: expected %s but got %s", e.name, e.msg, msg)
		}
	}
}

func Test_prompt(t *testing.T) {
	// save copy of os.Stdout
	oldStdout := os.Stdout

	// create a read and write pipe
	r, w, _ := os.Pipe()

	// set os.Stdout to our write pipe
	os.Stdout = w

	prompt()

	// close writer
	_ = w.Close()

	// reset os.Stdout to oldStdout
	os.Stdout = oldStdout

	// read the output of prompt function from read pipe
	out, _ := io.ReadAll(r)

	if string(out) != "-> " {
		t.Errorf("unexpected prompt: expected -> but got %s", string(out))
	}
}

func Test_intro(t *testing.T) {
	// save copy of os.Stdout
	oldStdout := os.Stdout

	// create a read and write pipe
	r, w, _ := os.Pipe()

	// set os.Stdout to our write pipe
	os.Stdout = w

	intro()

	// close writer
	_ = w.Close()

	// reset os.Stdout to oldStdout
	os.Stdout = oldStdout

	// read the output of prompt function from read pipe
	out, _ := io.ReadAll(r)

	if !strings.Contains(string(out), "Enter a whole number") {
		t.Errorf("unexpected intro text; got %s", string(out))
	}
}

func Test_checkNumbers(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedMsg    string
		expectedIsDone bool
	}{
		{"prime", "7", "7 is a prime number!", false},
		{"not prime", "8", "8 is not a prime number because it is divisible by 2!", false},
		{"zero", "0", "0 is not prime, by definition!", false},
		{"one", "1", "1 is not prime, by definition!", false},
		{"negative number", "-11", "Negative numbers are not prime, by definition!", false},
		{"not number", "ini text", "Please enter a whole number!", false},
		{"empty", "", "Please enter a whole number!", false},
		{"lowercase q", "q", "", true},
		{"uppercase Q", "Q", "", true},
	}

	errFmt := "Error '%s'\n\tExpected: %s\n\tGot: %s"

	for _, tt := range tests {
		input := strings.NewReader(tt.input)
		reader := bufio.NewScanner(input)
		msg, isDone := checkNumbers(reader)

		if msg != tt.expectedMsg {
			t.Errorf(errFmt, tt.name, tt.expectedMsg, msg)
		}

		if isDone != tt.expectedIsDone {
			t.Errorf(errFmt, tt.name, tt.expectedIsDone, isDone)
		}
	}
}

func Test_readUserInput(t *testing.T) {
	doneChan := make(chan bool)

	var stdin bytes.Buffer

	stdin.Write([]byte("1\nq\n")) // press 1 -> enter -> q -> enter

	go readUserInput(&stdin, doneChan)
	<-doneChan
	close(doneChan)
}
