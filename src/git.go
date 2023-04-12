package main

import (
	//"fmt"

	"os"

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

func readFileFromRepo(localPath string, filePath string) ([]byte, error) {
	// Open the local repository
	repo, err := git.PlainOpen(localPath)
	if err != nil {
		return nil, err
	}

	// Get the head reference
	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	// Get the commit object
	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	// Get the file contents from the commit tree
	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	file, err := tree.File(filePath)
	if err != nil {
		return nil, err
	}

	content, err := file.Contents()
	if err != nil {
		return nil, err
	}

	return []byte(content), nil
}
