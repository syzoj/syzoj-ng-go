package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

var log = logrus.StandardLogger()

// Get listen port for http server from environment variable HTTP_LISTEN_PORT.
func GetHttpListenPort() (int, error) {
	return strconv.Atoi(os.Getenv("HTTP_LISTEN_PORT"))
}

// Get URL for http endpoint from environment variable HTTP_URL.
func GetHttpURL(prefix string) (*url.URL, error) {
	u := os.Getenv(prefix + "HTTP_URL")
	if u == "" {
		return nil, fmt.Errorf("Environment variable %sHTTP_URL not found", prefix)
	}
	return url.Parse(u)
}
