package config

import (
	"github.com/elastic/go-elasticsearch"
	"os"
)

func OpenElastic(name string) (*elasticsearch.Client, error) {
	return elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{os.Getenv(name + "_ELASTICSEARCH")},
	})
}
