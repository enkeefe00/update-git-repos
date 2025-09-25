package main

import (
	"errors"

	gitcfg "github.com/go-git/go-git/v5/config"
)

// getRemoteName retrieves the remote name of a repository. If the repository has an 'upstream'
// entry, 'upstream' is returned; otherwise 'origin' is returned.
func getRemoteName(repoCfg gitcfg.Config) (string, error) {
	remotes, ok := repoCfg.Remotes["upstream"]
	if ok {
		if len(remotes.URLs) != 0 {
			return "upstream", nil
		}
	}

	remotes, ok = repoCfg.Remotes["origin"]
	if ok {
		if len(remotes.URLs) != 0 {
			return "origin", nil
		}
	}

	return "", errors.New("repository doesn't have an 'upstream' or 'origin' remote URL")
}

// getMainBranchName returns the 'main' branch name, which should be 'main' or 'master'
func getMainBranchName(repoCfg gitcfg.Config) (string, error) {
	branch, ok := repoCfg.Branches["main"]
	if !ok {
		branch, ok = repoCfg.Branches["master"]
		if !ok {
			return "", errors.New("repository doesn't have a 'main' or 'master' branch")
		}
	}

	return branch.Name, nil
}
