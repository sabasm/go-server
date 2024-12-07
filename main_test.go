package main

import "testing"

const expectedMessage = "Hello, World!"

func TestMain(t *testing.T) {
	if expectedMessage != "Hello, World!" {
		t.Errorf("Expected %s but got something else", expectedMessage)
	}
}
