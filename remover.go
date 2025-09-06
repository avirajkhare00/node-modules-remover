package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type NodeModulesRemover struct {
	config *Config
}

type NodeModuleInfo struct {
	Path     string
	Size     int64
	Modified time.Time
	Age      time.Duration
}

func NewNodeModulesRemover(config *Config) *NodeModulesRemover {
	return &NodeModulesRemover{
		config: config,
	}
}

func (r *NodeModulesRemover) Run() error {
	var totalFound int
	var totalRemoved int
	var totalSize int64

	for _, dir := range r.config.Directories {
		if r.config.Verbose && !r.config.Quiet {
			fmt.Printf("Scanning directory: %s\n", dir)
		}

		found, removed, size, err := r.processDirectory(dir)
		if err != nil {
			if !r.config.Quiet {
				fmt.Fprintf(os.Stderr, "Error processing directory %s: %v\n", dir, err)
			}
			continue
		}

		totalFound += found
		totalRemoved += removed
		totalSize += size
	}

	// Summary output
	if !r.config.Quiet {
		if r.config.DryRun {
			fmt.Printf("\nDRY RUN SUMMARY:\n")
			fmt.Printf("Found %d node_modules directories older than %v\n", totalFound, r.config.Age)
			fmt.Printf("Would free approximately %.2f MB\n", float64(totalSize)/(1024*1024))
		} else {
			fmt.Printf("\nSUMMARY:\n")
			fmt.Printf("Found %d node_modules directories older than %v\n", totalFound, r.config.Age)
			fmt.Printf("Removed %d directories\n", totalRemoved)
			fmt.Printf("Freed approximately %.2f MB\n", float64(totalSize)/(1024*1024))
		}
	}

	return nil
}

func (r *NodeModulesRemover) processDirectory(rootDir string) (found, removed int, totalSize int64, err error) {
	err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Skip directories we can't access
			return nil
		}

		// Only process directories named "node_modules"
		if !info.IsDir() || info.Name() != "node_modules" {
			return nil
		}

		// Check if the directory is old enough
		age := time.Since(info.ModTime())
		if age < r.config.Age {
			if r.config.Verbose && !r.config.Quiet {
				fmt.Printf("  Skipping %s (age: %v, threshold: %v)\n", path, age, r.config.Age)
			}
			return filepath.SkipDir // Skip subdirectories of this node_modules
		}

		found++

		// Calculate directory size
		size, err := r.calculateDirSize(path)
		if err != nil {
			if r.config.Verbose && !r.config.Quiet {
				fmt.Printf("  Warning: Could not calculate size for %s: %v\n", path, err)
			}
			size = 0
		}

		totalSize += size

		if r.config.DryRun {
			if !r.config.Quiet {
				fmt.Printf("  [DRY RUN] Would remove: %s (age: %v, size: %.2f MB)\n",
					path, age, float64(size)/(1024*1024))
			}
		} else {
			if err := os.RemoveAll(path); err != nil {
				if !r.config.Quiet {
					fmt.Fprintf(os.Stderr, "  Error removing %s: %v\n", path, err)
				}
				return nil
			}

			if !r.config.Quiet {
				fmt.Printf("  Removed: %s (age: %v, size: %.2f MB)\n",
					path, age, float64(size)/(1024*1024))
			}
		}

		removed++

		// Skip subdirectories of node_modules to avoid processing nested node_modules
		return filepath.SkipDir
	})

	return found, removed, totalSize, err
}

func (r *NodeModulesRemover) calculateDirSize(path string) (int64, error) {
	var size int64

	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files we can't access
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	return size, err
}
