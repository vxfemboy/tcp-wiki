package main

import (
	//"fmt"

	"os"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func cloneRepository(repoURL, localPath string) error {
	_, err := git.PlainClone(localPath, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})
	return err
}

func pullRepository(localPath, branch string) error {
	repo, err := git.PlainOpen(localPath)
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = worktree.Pull(&git.PullOptions{
		RemoteName:    "origin",
		ReferenceName: plumbing.NewBranchReferenceName(branch),
	})
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

func getAuthorAndLastModification(localPath, filePath string) (string, time.Time, string, time.Time, error) {
	repo, err := git.PlainOpen(localPath)
	if err != nil {
		return "", time.Time{}, "", time.Time{}, err
	}

	iter, err := repo.Log(&git.LogOptions{FileName: &filePath})
	if err != nil {
		return "", time.Time{}, "", time.Time{}, err
	}

	var firstCommit *object.Commit
	var lastCommit *object.Commit

	err = iter.ForEach(func(c *object.Commit) error {
		if firstCommit == nil {
			firstCommit = c
		}
		lastCommit = c
		return nil
	})
	if err != nil {
		return "", time.Time{}, "", time.Time{}, err
	}

	return firstCommit.Author.Name, firstCommit.Author.When, lastCommit.Author.Name, lastCommit.Author.When, nil
}
