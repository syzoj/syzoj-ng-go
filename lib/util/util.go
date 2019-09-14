// Utility library.
package util

import (
	"crypto/rand"
	"encoding/hex"
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
