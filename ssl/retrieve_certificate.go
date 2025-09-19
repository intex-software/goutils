package keymaster

import (
	"bytes"
	"crypto/tls"
	"encoding/pem"
	"fmt"
)

func GetTlsCertificateAsPem(host string) (content string, err error) {
	config := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", host, config)
	if err != nil {
		return "", fmt.Errorf("failed to dial TLS connection: %w", err)
	}
	defer conn.Close()

	state := conn.ConnectionState()

	if len(state.PeerCertificates) == 0 {
		return "", fmt.Errorf("no peer certificates found")
	}

	var buffer bytes.Buffer
	for _, cert := range state.PeerCertificates {
		err = pem.Encode(&buffer, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		})
		if err != nil {
			return "", fmt.Errorf("failed to PEM encode certificate: %w", err)
		}
		return buffer.String(), nil
	}

	return "", nil
}
