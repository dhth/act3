package utils

import "regexp"

var repoRegex = regexp.MustCompile("^[a-zA-Z0-9_.-]+/[a-zA-Z0-9_.-]+$")

func IsRepoNameValid(name string) bool {
	return repoRegex.MatchString(name)
}
