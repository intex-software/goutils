package obfuscate

import (
	"math/rand"

	"fiurthorn.de/goutils/internal"
)

// ObfuscateString takes a string content and obfuscates it by encoding it as base32.
// It returns the obfuscated string and an error if any.
func ObfuscateString(content []byte) (string, error) {
	return internal.Base32.EncodeToString(ObfuscateBytes(content)), nil
}

// DeobfuscateString takes an obfuscated content string and returns the deobfuscated secret string.
// It decodes the content string using base32 decoding and then applies the Deobfuscate function to the decoded string.
// If an error occurs during decoding, the original content string is returned along with the error.
// Otherwise, the deobfuscated secret string is returned along with a nil error.
func DeobfuscateString(content string) ([]byte, error) {
	raw, err := internal.Base32.DecodeString(content)
	if err != nil {
		return []byte(content), err
	}
	return DeobfuscateBytes(raw), nil
}

// ObfuscateBytes takes a string content and returns an obfuscated version of it.
// It uses a random salt value to perform bitwise XOR operation on each byte of the content.
// The obfuscated content is then returned as a string.
func ObfuscateBytes(content []byte) []byte {
	salt := byte(rand.Intn(256))

	for i := range len(content) {
		content[i] ^= salt
	}

	content = append([]byte{salt}, content...)
	return content
}

// DeobfuscateBytes takes an obfuscated content string and returns the original content string.
// It performs a deobfuscation process by XORing each byte of the content with a salt value,
// and then converting the result back to a string.
func DeobfuscateBytes(content []byte) []byte {
	salt := content[0]
	content = content[1:]

	for i := range len(content) {
		content[i] ^= salt
	}

	return content
}
