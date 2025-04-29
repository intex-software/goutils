package krypta

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"

	"github.com/intex-software/goutils/internal"
)

func (c krypta) Decode(encoded string) (string, error) {
	if !c.isEncoded(encoded) {
		return encoded, nil
	}

	from, till := len(c.prefix), len(encoded)-len(c.suffix)
	decoded, err := c.decrypt(encoded[from:till])
	if err != nil {
		return encoded, err
	}

	return decoded, nil
}

func (c krypta) decrypt(encoded string) (string, error) {
	if len(encoded) == 0 {
		return encoded, nil
	}
	cipherText, err := internal.Base32.DecodeString(encoded)
	if err != nil {
		return encoded, err
	}
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return encoded, err
	}
	if len(cipherText) < aes.BlockSize {
		err = errors.New("ciphertext block size is too short")
		return encoded, err
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)
	decoded := string(cipherText)

	return decoded, nil
}
