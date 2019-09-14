package config

import (
	"os"

	"github.com/minio/minio-go"
)

// Creates an Amazon S3 compatible client from environment variables.
// The environment variables are ${prefix}ENDPOINT, ${prefix}ACCESS_KEY, ${prefix}SECRET_KEY.
// The prefix defaults to "S3_" if not specified.
func NewMinio(prefix string) (*minio.Client, error) {
	if prefix == "" {
		prefix = "S3_"
	}
	endpoint := os.Getenv(prefix + "ENDPOINT")
	accessKey := os.Getenv(prefix + "ACCESS_KEY")
	secretKey := os.Getenv(prefix + "SECRET_KEY")
	return minio.New(endpoint, accessKey, secretKey, false)
}
