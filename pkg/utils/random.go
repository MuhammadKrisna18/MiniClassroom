package utils

import (
	"math/rand"
	"time"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const numberCharset = "0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// GenerateRandomString generates a random alphanumeric string of the specified length.
func GenerateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// GenerateRandomNumberString generates a random numeric string of the specified length.
func GenerateRandomNumberString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = numberCharset[seededRand.Intn(len(numberCharset))]
	}
	return string(b)
}
