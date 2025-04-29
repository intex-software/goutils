package rand

import (
	"crypto/rand"
	"encoding/base32"
	"strings"

	"github.com/intex-software/goutils/secret"
)

var base32Encoder = base32.StdEncoding.WithPadding(base32.NoPadding)

func RandomByteSlice(size int) []byte {
	token := make([]byte, size)
	rand.Read(token)
	return token
}

func RandomString(size int) string {
	token := RandomByteSlice(size)
	return strings.ToLower(base32Encoder.EncodeToString(token))
}

func RandomLength(size int) int {
	return base32Encoder.EncodedLen(size)
}

func NewSecret(size int) *secret.Secret {
	sec := secret.Secret(RandomByteSlice(size))
	return &sec
}

func EncodeBase32(secret []byte) string {
	return strings.ToLower(base32Encoder.EncodeToString(secret))
}

func DecodeBase32(secret string) ([]byte, error) {
	return base32Encoder.DecodeString(strings.ToUpper(secret))
}

func DecodeBase32String(secret string) (dst string, err error) {
	v, err := DecodeBase32(secret)
	dst = string(v)
	return
}
