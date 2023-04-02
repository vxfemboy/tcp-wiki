package main

import (
	//"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	git "github.com/go-git/go-git/v5"
)

func cloneRepository(repoURL, localPath string) error {
	_, err := git.PlainClone(localPath, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})
	return err
}

func pullRepository(localPath string) error {
	repo, err := git.PlainOpen(localPath)
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = worktree.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}
	return nil
}

func readFileFromRepo(localPath, filePath string) ([]byte, error) {
	absPath := filepath.Join(localPath, filePath)
	content, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}
	return content, nil
}
