package server

import (
	"testing"
	"time"
)

func TestWithTimeout(t *testing.T) {
	cfg := &Config{Host: "localhost", Port: 8080}
	builder := NewBuilder(cfg).WithTimeout(5*time.Second, 10*time.Second, 15*time.Second)

	if builder.timeouts.read != 5*time.Second {
		t.Errorf("expected read timeout of 5s, got %v", builder.timeouts.read)
	}

	if builder.timeouts.write != 10*time.Second {
		t.Errorf("expected write timeout of 10s, got %v", builder.timeouts.write)
	}

	if builder.timeouts.idle != 15*time.Second {
		t.Errorf("expected idle timeout of 15s, got %v", builder.timeouts.idle)
	}
}
