package rand

import (
	"crypto/rand"
	"strings"

	"github.com/intex-software/goutils/internal"
	"github.com/intex-software/goutils/secret"
)

func RandomByteSlice(size int) []byte {
	token := make([]byte, size)
	rand.Read(token)
	return token
}

func RandomString(size int) string {
	token := RandomByteSlice(size)
	return internal.Base32.EncodeToString(token)
}

func RandomLength(size int) int {
	return internal.Base32.EncodedLen(size)
}

func NewSecret(size int) *secret.Secret {
	sec := secret.Secret(RandomByteSlice(size))
	return &sec
}

func EncodeBase32(secret []byte) string {
	return internal.Base32.EncodeToString(secret)
}

func DecodeBase32(secret string) ([]byte, error) {
	return internal.Base32.DecodeString(strings.ToLower(secret))
}

func DecodeBase32String(secret string) (dst string, err error) {
	v, err := DecodeBase32(secret)
	dst = string(v)
	return
}
