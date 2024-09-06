package secret

import (
	"bytes"
	"encoding/pem"

	"golang.org/x/crypto/pkcs12"
)

func (k *Secret) Pem(password string) (data []byte, err error) {
	data = k.Bytes()

	if data[0] == '-' {
		return
	}

	pems, err := pkcs12.ToPEM(data, password)
	if err != nil {
		return
	}

	buffer := &bytes.Buffer{}
	for _, p := range pems {
		buffer.Write(pem.EncodeToMemory(p))
	}
	data = buffer.Bytes()

	return
}
