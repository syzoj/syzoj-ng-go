package config

import (
	"github.com/elastic/go-elasticsearch"
	"os"
)

// Creates an elasticsearch client from environment variables.
func NewElastic(prefix string) (*elasticsearch.Client, error) {
	return elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{os.Getenv(prefix + "ELASTICSEARCH")},
	})
}
