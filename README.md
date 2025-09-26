# Git Repository Updater

A Go-based command-line tool that automatically updates all Git repositories in a specified directory with their latest remote commits. This tool is particularly useful for managing multiple forked repositories, as it intelligently handles both regular repositories and forks by pulling from the appropriate remote.

*Non-fast-forward merges are not yet handled.*

## Installation

## Dependencies

- [go-git/go-git](https://github.com/go-git/go-git) - Pure Go Git implementation

### Prerequisites

- Go 1.23.4 or later
- Git repositories must have:either an `origin` or `upstream` remote configured
- Git repositories must have a `main` or `master` branch
- Write permissions to the target directory for log file creation

### Build from Source

```bash
git clone https://github.com/enkeefe0/git-repo-update.git
cd git-repo-update
go build -o git-repo-update
```

## Usage

```bash
./git-repo-update <directory-path>
```

## How It Works

1. **Directory Scanning**: The tool walks through the specified directory to find all subdirectories containing a `.git` folder
2. **Repository Analysis**: For each repository, it:
   - Opens the repository and reads its configuration
   - Determines the appropriate remote (`origin` or `upstream`)
   - Identifies the main branch (`main` or `master`)
3. **Update Process**: 
   - Checks out the main branch
   - Pulls the latest changes from the remote
   - For forked repositories (with `upstream` remote), also pushes changes to `origin`

## Logging

The tool creates detailed logs in a `repo_updates` directory within the target directory. Log files are named with timestamps (e.g., `Jan-01-25_15-04.txt`).
