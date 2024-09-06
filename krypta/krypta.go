package krypta

import (
	"crypto/sha256"
)

type Krypta interface {
	Encode(decoded string) (string, error)
	Decode(encoded string) (string, error)

	encrypt(decoded string) (string, error)
	decrypt(encoded string) (string, error)
}

type krypta struct {
	key    []byte
	prefix string
	suffix string
}

func NewKrypta(key string, prefix, suffix string) Krypta {
	hash := sha256.New()
	hash.Write([]byte(key))

	return krypta{
		key:    hash.Sum(nil),
		prefix: prefix,
		suffix: suffix,
	}
}
