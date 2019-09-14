package model

// File contains information about a file.
type File struct {
	// A URL to download from, like a pre-signed S3 GetObject.
	URL string `json:"url,omitempty"`
}

// Testset contains information about how submissions should be judged.
type Testset struct {
	Id string `json:"id,omitempty"`
}
