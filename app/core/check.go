package core

import (
	"regexp"
)

var nameRegexp = regexp.MustCompile("[0-9A-Za-z-_]{1,16}")

func checkName(name string) bool {
	return nameRegexp.Match([]byte(name))
}
