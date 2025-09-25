package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func main() {
	args := os.Args[1:]
	gitRoot := args[0]

	// Create the log directory path
	logDirPath := filepath.Join(gitRoot, "repo_updates")
	if err := os.MkdirAll(logDirPath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating log directory '%s': %v\n", logDirPath, err)
		os.Exit(1)
	}

	// Create the timestamped log filepath
	currentTime := time.Now()
	filename := fmt.Sprintf("%s.txt", currentTime.Format("Jan-02-06_15-04"))
	fullFilePath := filepath.Join(logDirPath, filename)

	originalStdout := os.Stdout
	originalStderr := os.Stderr

	// Create loggers
	fileOnlyLogger, err := os.OpenFile(fullFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(originalStderr, "Error opening log file '%s': %v\n", fullFilePath, err)
		os.Exit(1)
	}
	defer fileOnlyLogger.Close() // Ensure the file is closed when main exits
	consoleAndFileLogger := io.MultiWriter(originalStdout, fileOnlyLogger)

	// Initial startup messages that always go to both console and file
	fmt.Fprintln(consoleAndFileLogger, "Starting Git repository update process...")
	fmt.Fprintf(consoleAndFileLogger, "All detailed output will be logged to: %s\n", fullFilePath)

	// Check if the provided directory exists.
	if _, err := os.Stat(gitRoot); os.IsNotExist(err) {
		fmt.Fprintf(consoleAndFileLogger, "Error: The directory '%s' does not exist. Please ensure your Git repositories are in this folder.\n", gitRoot)
		os.Exit(1)
	}

	// Walk through the provided directory to find all Git repositories.
	fmt.Fprintf(fileOnlyLogger, "Scanning for Git repositories in: %s\n", gitRoot)
	err = filepath.WalkDir(gitRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(consoleAndFileLogger, "Error accessing path %s: %v\n", path, err)
			return nil
		}

		if !d.IsDir() {
			return nil
		}

		gitDir := filepath.Join(path, ".git")
		if info, err := os.Stat(gitDir); err == nil && info.IsDir() {
			fmt.Fprintf(consoleAndFileLogger, "\n--- Found repository: %s ---\n", path)

			err := updateGitRepo(path, fileOnlyLogger)
			if err != nil {
				fmt.Fprintf(consoleAndFileLogger, "!!! Failed to update repository %s: %v\n", path, err)
			} else {
				fmt.Fprintf(consoleAndFileLogger, "+++ Successfully updated repository: %s\n", path)
			}
			return filepath.SkipDir
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(consoleAndFileLogger, "\nAn error occurred during directory traversal: %v\n", err)
	}

	fmt.Fprintln(consoleAndFileLogger, "\nGit repository update process completed.")
}
