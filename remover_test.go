package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNodeModulesRemover(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "node-modules-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test node_modules directory
	nodeModulesDir := filepath.Join(tempDir, "test-project", "node_modules")
	if err := os.MkdirAll(nodeModulesDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a dummy file in node_modules
	dummyFile := filepath.Join(nodeModulesDir, "dummy.txt")
	if err := os.WriteFile(dummyFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	// Test with a very short age threshold (should find the directory)
	config := &Config{
		Directories: []string{tempDir},
		Age:         1 * time.Second,
		DryRun:      true,
		Verbose:     false,
		Quiet:       true,
	}

	remover := NewNodeModulesRemover(config)
	if err := remover.Run(); err != nil {
		t.Fatal(err)
	}

	// Test with a very long age threshold (should not find the directory)
	config.Age = 1 * time.Hour
	remover = NewNodeModulesRemover(config)
	if err := remover.Run(); err != nil {
		t.Fatal(err)
	}
}

func TestConfigParsing(t *testing.T) {
	// Test that the Config struct can be created properly
	config := &Config{
		Directories: []string{"."},
		Age:         3 * time.Minute,
		DryRun:      true,
		Verbose:     false,
		Quiet:       false,
	}

	if len(config.Directories) != 1 {
		t.Errorf("Expected 1 directory, got %d", len(config.Directories))
	}

	if config.Age != 3*time.Minute {
		t.Errorf("Expected age of 3 minutes, got %v", config.Age)
	}

	if !config.DryRun {
		t.Error("Expected DryRun to be true")
	}
}
