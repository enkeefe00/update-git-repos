package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// Get the user's home directory.
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	// Create the log directory path
	logDirPath := filepath.Join(homeDir, "repo_updates")
	if err := os.MkdirAll(logDirPath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating log directory '%s': %v\n", logDirPath, err)
		os.Exit(1)
	}

	// Generate the timestamped filename
	currentTime := time.Now()
	filename := fmt.Sprintf("%s.txt", currentTime.Format("Jan-02-06_15-04"))
	fullFilePath := filepath.Join(logDirPath, filename)

	// Capture the original standard output and standard error for console-only messages
	originalStdout := os.Stdout
	originalStderr := os.Stderr

	// Open the log file for writing (append if exists, create if not)
	fileOnlyLogger, err := os.OpenFile(fullFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// This error is critical and should always go to console
		fmt.Fprintf(originalStderr, "Error opening log file '%s': %v\n", fullFilePath, err)
		os.Exit(1)
	}
	defer fileOnlyLogger.Close() // Ensure the file is closed when main exits

	// Create a MultiWriter that writes to both the console (originalStdout) and the file (f).
	// This will be used for messages that need to appear in both places.
	consoleAndFileLogger := io.MultiWriter(originalStdout, fileOnlyLogger)

	// Initial startup messages that always go to both console and file
	fmt.Fprintln(consoleAndFileLogger, "Starting Git repository update process...")
	fmt.Fprintf(consoleAndFileLogger, "All detailed output will be logged to: %s\n", fullFilePath)

	// Construct the full path to the Git repositories folder.
	gitRoot := filepath.Join(homeDir, "git")

	// Check if the ~/git directory exists.
	if _, err := os.Stat(gitRoot); os.IsNotExist(err) {
		fmt.Fprintf(consoleAndFileLogger, "Error: The directory '%s' does not exist. Please ensure your Git repositories are in this folder.\n", gitRoot)
		os.Exit(1)
	}

	fmt.Fprintf(fileOnlyLogger, "Scanning for Git repositories in: %s\n", gitRoot)

	// Walk through the ~/git directory to find all Git repositories.
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
