package cmd

//
import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	ghapi "github.com/cli/go-gh/v2/pkg/api"
	"github.com/dhth/act3/internal/gh"
	"github.com/go-git/go-git/v5"
)

var (
	errCouldntGetRepo        = errors.New("couldn't get repository")
	errNoRemotesFound        = errors.New("no remotes found")
	errRemoteURLEmpty        = errors.New("remote URL is empty")
	errInvalidURLFormat      = errors.New("remote URL has invalid format")
	errCouldntParseRemoteURL = errors.New("couldn't parse remote URL")
)

func getWorkflowsForCurrentRepo(ghClient *ghapi.RESTClient, repo string) ([]gh.Workflow, error) {
	wd, err := gh.GetWorkflowDetails(ghClient, repo)
	if err != nil {
		return nil, err
	}

	workflows := make([]gh.Workflow, len(wd.Workflows))
	for i, w := range wd.Workflows {
		workflows[i] = gh.Workflow{
			ID:   w.NodeID,
			Repo: repo,
			Name: w.Name,
		}
	}

	return workflows, nil
}

func getCurrentRepo() (string, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return "", fmt.Errorf("%w: %s", errCouldntGetRepo, err.Error())
	}

	remotes, err := repo.Remotes()
	if err != nil {
		return "", fmt.Errorf("%w: %s", errNoRemotesFound, err.Error())
	}

	if len(remotes) == 0 {
		return "", fmt.Errorf("%w", errNoRemotesFound)
	}

	remote := remotes[0]
	if remote == nil {
		return "", fmt.Errorf("%w", errNoRemotesFound)
	}

	remoteURL := remote.Config().URLs[0]

	userRepo, err := extractRepoName(remoteURL)
	if err != nil {
		return "", fmt.Errorf("%w: %s", errCouldntParseRemoteURL, err.Error())
	}

	return userRepo, nil
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
