package main

import (
	"os"
	"path/filepath"
	"testing"

	git "github.com/go-git/go-git/v5"
	gitcfg "github.com/go-git/go-git/v5/config"
)

func TestGetRemoteName(t *testing.T) {
	// Create a temporary directory for test repositories
	tempDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name         string
		setupRepoCfg func(string) (*gitcfg.Config, error)
		expectedName string
		expectError  bool
	}{
		{
			name: "repository with upstream remote",
			setupRepoCfg: func(repoPath string) (*gitcfg.Config, error) {
				// Initialize a new repository
				repo, err := git.PlainInit(repoPath, false)
				if err != nil {
					return nil, err
				}

				// Add upstream remote
				_, err = repo.CreateRemote(&gitcfg.RemoteConfig{
					Name: "upstream",
					URLs: []string{"https://github.com/owner/repo.git"},
				})
				if err != nil {
					return nil, err
				}

				return repo.Config()
			},
			expectedName: "upstream",
			expectError:  false,
		},
		{
			name: "repository with origin remote",
			setupRepoCfg: func(repoPath string) (*gitcfg.Config, error) {
				// Initialize a new repository
				repo, err := git.PlainInit(repoPath, false)
				if err != nil {
					return nil, err
				}

				// Add upstream remote
				_, err = repo.CreateRemote(&gitcfg.RemoteConfig{
					Name: "origin",
					URLs: []string{"https://github.com/owner/repo.git"},
				})
				if err != nil {
					return nil, err
				}

				return repo.Config()
			},
			expectedName: "origin",
			expectError:  false,
		},
		{
			name: "repository with origin & upstream remotes",
			setupRepoCfg: func(repoPath string) (*gitcfg.Config, error) {
				repo, err := git.PlainInit(repoPath, false)
				if err != nil {
					return nil, err
				}

				// Add upstream remote with multiple URLs
				_, err = repo.CreateRemote(&gitcfg.RemoteConfig{
					Name: "upstream",
					URLs: []string{
						"https://github.com/owner/repo.git",
					},
				})
				if err != nil {
					return nil, err
				}

				_, err = repo.CreateRemote(&gitcfg.RemoteConfig{
					Name: "origin",
					URLs: []string{
						"https://github.com/owner/repo.git",
					},
				})
				if err != nil {
					return nil, err
				}

				return repo.Config()
			},
			expectedName: "upstream",
			expectError:  false,
		},
		{
			name: "repository with empty upstream URLs",
			setupRepoCfg: func(repoPath string) (*gitcfg.Config, error) {
				repo, err := git.PlainInit(repoPath, false)
				if err != nil {
					return nil, err
				}

				// Add upstream remote with no URLs
				_, err = repo.CreateRemote(&gitcfg.RemoteConfig{
					Name: "upstream",
					URLs: []string{""},
				})
				if err != nil {
					return nil, err
				}

				return repo.Config()
			},
			expectedName: "upstream",
			expectError:  false,
		},
		{
			name: "invalid remote name",
			setupRepoCfg: func(repoPath string) (*gitcfg.Config, error) {
				// Initialize a new repository
				repo, err := git.PlainInit(repoPath, false)
				if err != nil {
					return nil, err
				}

				// Add upstream remote
				_, err = repo.CreateRemote(&gitcfg.RemoteConfig{
					Name: "invalid",
					URLs: []string{"https://github.com/owner/repo.git"},
				})
				if err != nil {
					return nil, err
				}

				return repo.Config()
			},
			expectedName: "",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a unique subdirectory for this test
			testRepoPath := filepath.Join(tempDir, tt.name)
			err := os.MkdirAll(testRepoPath, 0755)
			if err != nil {
				t.Fatalf("Failed to create test directory: %v", err)
			}

			// Setup the repository according to the test case
			repoCfg, err := tt.setupRepoCfg(testRepoPath)
			if err != nil {
				t.Fatalf("Failed to setup test repository: %v", err)
			}

			// Test the getUpstreamURL function
			url, err := getRemoteName(*repoCfg)

			// Check error expectations
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			// Check URL expectations
			if url != tt.expectedName {
				t.Errorf("Expected '%s', but got '%s'", tt.expectedName, url)
			}
		})
	}
}

func TestGetMainBranchName(t *testing.T) {
	tests := []struct {
		name           string
		repoCfg        gitcfg.Config
		expectedBranch string
		expectError    bool
	}{
		{
			name: "repository with main branch",
			repoCfg: gitcfg.Config{
				Branches: map[string]*gitcfg.Branch{
					"main": {
						Name: "main",
					},
				},
			},
			expectedBranch: "main",
			expectError:    false,
		},
		{
			name: "repository with master branch",
			repoCfg: gitcfg.Config{
				Branches: map[string]*gitcfg.Branch{
					"master": {
						Name: "master",
					},
				},
			},
			expectedBranch: "master",
			expectError:    false,
		},
		{
			name: "repository with both main and master branches",
			repoCfg: gitcfg.Config{
				Branches: map[string]*gitcfg.Branch{
					"main": {
						Name: "main",
					},
					"master": {
						Name: "master",
					},
				},
			},
			expectedBranch: "main",
			expectError:    false,
		},
		{
			name: "repository with neither main nor master branch",
			repoCfg: gitcfg.Config{
				Branches: map[string]*gitcfg.Branch{
					"develop": {
						Name: "develop",
					},
					"feature": {
						Name: "feature",
					},
				},
			},
			expectedBranch: "",
			expectError:    true,
		},
		{
			name: "repository with empty branch configuration",
			repoCfg: gitcfg.Config{
				Branches: map[string]*gitcfg.Branch{},
			},
			expectedBranch: "",
			expectError:    true,
		},
		{
			name: "repository with nil branch configuration",
			repoCfg: gitcfg.Config{
				Branches: nil,
			},
			expectedBranch: "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the getMainBranchName function
			branchName, err := getMainBranchName(tt.repoCfg)

			// Check error expectations
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			// Check branch name expectations
			if branchName != tt.expectedBranch {
				t.Errorf("Expected '%s', but got '%s'", tt.expectedBranch, branchName)
			}
		})
	}
}
