package main

import (
	"fmt"
	"io"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// updateGitRepo fetches the latest commits for a given local repository directory.
func updateGitRepo(repoPath string, fileOnlyLogger io.Writer) error {
	// Open the repository & get the config and working directory
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return fmt.Errorf("failed to open repository at %s: %w", repoPath, err)
	}
	config, err := r.Config()
	if err != nil {
		return fmt.Errorf("failed to get repository config: %w", err)
	}
	workingDir, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get repository working directory: %w", err)
	}

	// Get remote and main branch names
	remoteName, err := getRemoteName(*config)
	if err != nil {
		return fmt.Errorf("failed to get a remote name for repository: %w", err)
	}

	mainBranchName, err := getMainBranchName(*config)
	if err != nil {
		return fmt.Errorf("failed to retrieve main branch for repository: %w", err)
	}

	// Checkout and pull form the remote branch
	fmt.Fprintf(fileOnlyLogger, "Checking out '%s' branch...\n", mainBranchName)
	branchRefName := plumbing.NewBranchReferenceName(mainBranchName)
	err = workingDir.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(branchRefName),
		Force:  true,
	})
	if err != nil {
		return fmt.Errorf("failed to checkout %s branch: %w", mainBranchName, err)
	}

	fmt.Fprintf(fileOnlyLogger, "Pulling updates from '%s' remote to local '%s' branch...\n",
		remoteName, mainBranchName)
	err = workingDir.Pull(&git.PullOptions{
		RemoteName: remoteName,
	})
	if err != nil {
		return fmt.Errorf("failed to pull updates from '%s' remote to local '%s' branch: %w",
			remoteName, mainBranchName, err)
	}

	// Only push changes to 'origin' if the repo is a fork
	// if remoteName == "upstream" {
	// 	fmt.Fprintf(fileOnlyLogger, "Pushing updates to 'origin' remote for %s branch...\n",
	// 		mainBranchName)
	// 	err = r.Push(&git.PushOptions{
	// 		RemoteName: "origin",
	// 	})
	// 	if err != nil {
	// 		return fmt.Errorf("failed to push updates to 'origin' remote for %s branch: %w",
	// 			mainBranchName, err)
	// 	}
	// }

	return nil
}
