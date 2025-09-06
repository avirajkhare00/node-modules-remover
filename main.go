package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var version = "dev"

type Config struct {
	Directories []string
	Age         time.Duration
	DryRun      bool
	Verbose     bool
	Quiet       bool
}

func main() {
	config := parseFlags()

	if len(config.Directories) == 0 {
		// Default to current directory if no directories specified
		config.Directories = []string{"."}
	}

	if config.Verbose && !config.Quiet {
		fmt.Printf("Scanning directories: %s\n", strings.Join(config.Directories, ", "))
		fmt.Printf("Age threshold: %v\n", config.Age)
		if config.DryRun {
			fmt.Println("DRY RUN MODE - No files will be deleted")
		}
		fmt.Println()
	}

	remover := NewNodeModulesRemover(config)
	if err := remover.Run(); err != nil {
		log.Fatal(err)
	}
}

func parseFlags() *Config {
	var ageStr string
	var directories string
	var dryRun bool
	var verbose bool
	var quiet bool

	flag.StringVar(&ageStr, "age", "3m", "Remove node_modules older than this duration (e.g., 3m, 90d, 24h)")
	flag.StringVar(&directories, "dirs", "", "Comma-separated list of directories to scan (default: current directory)")
	flag.BoolVar(&dryRun, "dry-run", false, "Show what would be deleted without actually deleting")
	flag.BoolVar(&verbose, "verbose", false, "Show detailed output")
	flag.BoolVar(&quiet, "quiet", false, "Minimal output (good for cron)")

	showHelp := flag.Bool("help", false, "Show help")
	showVersion := flag.Bool("version", false, "Show version")

	flag.Parse()

	if *showHelp {
		showHelpText()
		os.Exit(0)
	}

	if *showVersion {
		fmt.Printf("node-modules-remover %s\n", version)
		os.Exit(0)
	}

	// Parse age duration
	age, err := time.ParseDuration(ageStr)
	if err != nil {
		log.Fatalf("Invalid age duration '%s': %v", ageStr, err)
	}

	config := &Config{
		Age:     age,
		DryRun:  dryRun,
		Verbose: verbose,
		Quiet:   quiet,
	}

	// Parse directories
	if directories != "" {
		config.Directories = strings.Split(directories, ",")
		for i, dir := range config.Directories {
			config.Directories[i] = strings.TrimSpace(dir)
		}
	}

	return config
}

func showHelpText() {
	fmt.Println("node-modules-remover - Remove old node_modules directories")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  node-modules-remover [OPTIONS] [DIRECTORIES...]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -age duration    Remove node_modules older than this (default: 3m)")
	fmt.Println("  -dirs string     Comma-separated list of directories to scan")
	fmt.Println("  -dry-run         Show what would be deleted without actually deleting")
	fmt.Println("  -verbose         Show detailed output")
	fmt.Println("  -quiet           Minimal output (good for cron)")
	fmt.Println("  -help            Show this help message")
	fmt.Println("  -version         Show version")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  node-modules-remover")
	fmt.Println("  node-modules-remover -age 7d -verbose")
	fmt.Println("  node-modules-remover -dirs /path/to/projects,/another/path")
	fmt.Println("  node-modules-remover -dry-run -verbose")
}
