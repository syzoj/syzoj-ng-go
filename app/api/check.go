package api

import (
	"regexp"
)

var userNamePattern = regexp.MustCompile("^[0-9A-Za-z]{3,32}$")

func checkUserName(userName string) bool {
	return userNamePattern.MatchString(userName)
}
