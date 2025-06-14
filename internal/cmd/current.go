package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-git/go-git/v5"
)

var (
	ErrNotInAGitRepo         = errors.New("not in a git repo")
	errCouldntGetRepo        = errors.New("couldn't get repository")
	errCouldntFindRemotes    = errors.New("couldn't find remotes")
	errNoRemotesFound        = errors.New("no remotes found")
	errRemoteURLEmpty        = errors.New("remote URL is empty")
	errInvalidURLFormat      = errors.New("remote URL has invalid format")
	errCouldntParseRemoteURL = errors.New("couldn't parse remote URL")
)

func getCurrentRepo() (string, error) {
	var zero string
	repo, err := git.PlainOpen(".")
	if err != nil {
		if errors.Is(err, git.ErrRepositoryNotExists) {
			return zero, ErrNotInAGitRepo
		}
		return "", fmt.Errorf("%w: %s", errCouldntGetRepo, err.Error())
	}

	remotes, err := repo.Remotes()
	if err != nil {
		return zero, fmt.Errorf("%w: %s", errCouldntFindRemotes, err.Error())
	}

	if len(remotes) == 0 {
		return zero, errNoRemotesFound
	}

	remote := remotes[0]
	if remote == nil {
		return zero, errNoRemotesFound
	}

	remoteURL := remote.Config().URLs[0]

	repoName, err := extractRepoName(remoteURL)
	if err != nil {
		return zero, fmt.Errorf("%w: %s", errCouldntParseRemoteURL, err.Error())
	}

	return repoName, nil
}

func extractRepoName(remoteURL string) (string, error) {
	if remoteURL == "" {
		return "", fmt.Errorf("%w", errRemoteURLEmpty)
	}
	if strings.HasPrefix(remoteURL, "git@") {
		parts := strings.Split(remoteURL, ":")
		if len(parts) < 2 {
			return "", fmt.Errorf("%w", errInvalidURLFormat)
		}
		return strings.TrimSuffix(parts[1], ".git"), nil
	}

	parsedURL, err := url.Parse(remoteURL)
	if err != nil {
		return "", fmt.Errorf("%w: %s", errInvalidURLFormat, err.Error())
	}

	if parsedURL.Scheme == "" {
		return "", fmt.Errorf("%w: URL scheme is empty", errInvalidURLFormat)
	}

	if strings.Count(parsedURL.Path, "/") > 2 {
		return "", fmt.Errorf("%w", errInvalidURLFormat)
	}

	userRepo := strings.TrimSuffix(parsedURL.Path, ".git")
	userRepo = strings.TrimPrefix(userRepo, "/")

	return userRepo, nil
}
