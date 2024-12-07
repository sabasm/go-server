package main

import "testing"

func TestMain(t *testing.T) {
	expected := "Hello, World!"
	if expected != "Hello, World!" {
		t.Errorf("Expected %s but got something else", expected)
	}
}
