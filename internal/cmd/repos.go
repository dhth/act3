package cmd

import (
	"errors"
	"fmt"

	"github.com/dhth/act3/internal/utils"
)

var (
	errInvalidRepoProvided  = errors.New("invalid repo provided")
	errInvalidReposProvided = errors.New("invalid repos provided")
)

func validateRepos(repos []string) error {
	var invalidRepos []string
	for _, r := range repos {
		if !utils.IsRepoNameValid(r) {
			invalidRepos = append(invalidRepos, r)
		}
	}
	if len(invalidRepos) == 1 {
		return fmt.Errorf(`%w: %q; value needs to be in the format "owner/repo"`, errInvalidRepoProvided, invalidRepos[0])
	}

	if len(invalidRepos) > 1 {
		return fmt.Errorf(`%w: %q; value needs to be in the format "owner/repo"`, errInvalidReposProvided, invalidRepos)
	}

	return nil
}
