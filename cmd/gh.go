package cmd

import (
	"context"
	"errors"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

var (
	GHTokenNotProvided = errors.New("Github Access Token not provided")
)

func getGHClient() (*githubv4.Client, error) {
	accessToken := os.Getenv("ACT3_GH_ACCESS_TOKEN")

	if accessToken == "" {
		return nil, GHTokenNotProvided
	}
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)
	return client, nil
}
