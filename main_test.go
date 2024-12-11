package main

import (
	"os"
	"syscall"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	t.Run("server lifecycle", func(t *testing.T) {
		os.Setenv("APP_PORT", "0")
		defer os.Unsetenv("APP_PORT")

		errChan := make(chan error, 1)
		exitChan := make(chan struct{})

		go func() {
			main()
			close(exitChan)
		}()

		go func() {
			time.Sleep(100 * time.Millisecond)
			proc, _ := os.FindProcess(os.Getpid())
			if err := proc.Signal(syscall.SIGTERM); err != nil {
				errChan <- err
			}
		}()

		select {
		case err := <-errChan:
			t.Errorf("unexpected error: %v", err)
		case <-exitChan:
		case <-time.After(2 * time.Second):
			t.Error("test timed out")
		}
	})
}
