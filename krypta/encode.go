package krypta

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"strings"

	"github.com/intex-software/goutils/internal"
)

func (c krypta) isEncoded(val string) bool {
	return strings.HasPrefix(val, c.prefix) && strings.HasSuffix(val, c.suffix)
}

func (c krypta) Encode(decoded string) (string, error) {
	if c.isEncoded(decoded) {
		return decoded, nil
	}

	if len(decoded) == 0 {
		return decoded, nil
	}

	encoded, err := c.encrypt(decoded)
	if err != nil {
		return decoded, err
	}

	return c.prefix + encoded + c.suffix, nil
}

func (c krypta) encrypt(decoded string) (string, error) {
	plainText := []byte(decoded)
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return decoded, err
	}
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return decoded, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)
	encoded := internal.Base32.EncodeToString(cipherText)
	return encoded, nil
}
