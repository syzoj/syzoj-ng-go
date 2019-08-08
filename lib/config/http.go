package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

var log = logrus.StandardLogger()

func GetHttpListenPort() (int, error) {
	return strconv.Atoi(os.Getenv("HTTP_LISTEN_PORT"))
}

func GetHttpURL(name string) (string, error) {
	url := os.Getenv(name + "_HTTP_URL")
	if strings.HasSuffix(url, "/") {
		log.Warning(name + "_HTTP_URL" + " should not end with slash")
	}
	if url == "" {
		return "", fmt.Errorf("Environment variable %s_HTTP_URL not found", name)
	}
	return url, nil
}
