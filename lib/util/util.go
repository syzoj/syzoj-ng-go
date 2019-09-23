// Utility library.
package util

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// Creates a random hex string. n is the number of bytes.
func RandomHex(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

// Allocates a string.
func String(s string) *string {
	return &s
}

// Allocates a time.Time.
func Time(t time.Time) *time.Time {
	return &t
}

// Allocates a *int64.
func Int64(v int64) *int64 {
	return &v
}

// Allocates a *int.
func Int(v int) *int {
	return &v
}
