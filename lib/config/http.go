package config

import (
	"fmt"
	"net/url"
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.StandardLogger()

// Get listen port for http server from environment variable HTTP_LISTEN_ADDR.
func GetHttpListenAddr() string {
	return os.Getenv("HTTP_LISTEN_ADDR")
}

// Get URL for http endpoint from environment variable HTTP_URL.
func GetHttpURL(prefix string) (*url.URL, error) {
	u := os.Getenv(prefix + "HTTP_URL")
	if u == "" {
		return nil, fmt.Errorf("Environment variable %sHTTP_URL not found", prefix)
	}
	return url.Parse(u)
}
