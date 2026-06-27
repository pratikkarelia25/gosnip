package shortcode

import (
	"crypto/rand"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Generate returns a random short code of length n made up of
// URL-safe alphanumeric characters.
func Generate(n int) (string, error) {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	code := make([]byte, n)
	for i, b := range buf {
		code[i] = alphabet[int(b)%len(alphabet)]
	}
	return string(code), nil
}
