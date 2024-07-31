package util

import (
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
)

func GetLocalRegistryRepoPath() string {
	cs2pmPath := os.Getenv("CS2PM_REGISTRY_PATH")

	if cs2pmPath != "" {
		return filepath.Clean(cs2pmPath)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(homeDir, ".cs2pm")
}

func CloneRegistryRepo() *git.Repository {
	cloneDir := GetLocalRegistryRepoPath()
	repo, err := git.PlainClone(cloneDir, false, &git.CloneOptions{
		URL:      "https://github.com/hadley31/cs2pm.git",
		Progress: os.Stdout,
	})

	if err != nil {
		panic(err)
	}

	return repo
}

func GetLocalRegistryRepo() *git.Repository {
	repo, err := git.PlainOpen(GetLocalRegistryRepoPath())

	if err != nil {
		return nil
	}

	return repo
}

func GetOrCloneLocalRegistryRepo() *git.Repository {
	repo := GetLocalRegistryRepo()

	if repo == nil {
		repo = CloneRegistryRepo()
	}

	return repo
}

func PullLatestRegistryChanges() bool {
	worktree, err := GetLocalRegistryRepo().Worktree()

	if err != nil {
		panic(err)
	}

	err = worktree.Pull(&git.PullOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
	})

	return err == nil || err == git.NoErrAlreadyUpToDate
}
