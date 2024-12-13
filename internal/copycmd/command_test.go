package copycmd

import (
	"os"
	"testing"
)

func TestCopyFile(t *testing.T) {
	sourceContent := []byte("Hello, World!")
	sourceFile, err := os.CreateTemp("", "source*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(sourceFile.Name())

	if _, err := sourceFile.Write(sourceContent); err != nil {
		t.Fatal(err)
	}
	sourceFile.Close()

	destFile, err := os.CreateTemp("", "dest*.txt")
	if err != nil {
		t.Fatal(err)
	}
	destFile.Close()
	defer os.Remove(destFile.Name())

	if err := CopyFile(sourceFile.Name(), destFile.Name()); err != nil {
		t.Errorf("CopyFile failed: %v", err)
	}

	destContent, err := os.ReadFile(destFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if string(destContent) != string(sourceContent) {
		t.Errorf("Expected %s, got %s", sourceContent, destContent)
	}
}
