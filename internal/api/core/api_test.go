package api

import "testing"

func TestNewAPI(t *testing.T) {
	api := NewAPI()
	if api == nil {
		t.Error("Expected non-nil API instance")
	}
}
