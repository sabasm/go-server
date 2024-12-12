package main

import (
	"os"
	"syscall"
	"testing"
	"time"
)

func TestServerMain(t *testing.T) {
	os.Setenv("APP_PORT", "0")
	defer os.Unsetenv("APP_PORT")

	done := make(chan struct{})
	go func() {
		main()
		close(done)
	}()

	time.Sleep(100 * time.Millisecond)

	proc, _ := os.FindProcess(os.Getpid())
	if err := proc.Signal(syscall.SIGTERM); err != nil {
		t.Errorf("failed to send signal: %v", err)
	}

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Error("timeout waiting for server shutdown")
	}
}
