package certificates

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/pkcs12"

	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"

	des "software.sslmate.com/src/go-pkcs12"
)

type MsGraphCertificateResult struct {
	Keys []byte
	Cert []byte
}

func resolveSibling(filename, extension string) string {
	oldExt := filepath.Ext(filename)
	return strings.TrimSuffix(filename, oldExt) + extension
}

func WriteSelfSingedGraphAppCertificates(filename string, days time.Duration) (err error) {
	log.Println("Create Azure MsGraph Certificate")

	pemFile := resolveSibling(filename, ".pem")
	cerFile := resolveSibling(filename, ".cer")

	if _, err := os.Stat(pemFile); !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("file already exists: %s", pemFile)
	}

	if _, err := os.Stat(cerFile); !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("file already exists: %s", cerFile)
	}

	pem, cer, err := createSelfSingedGraphAppCertificate(days)
	if err != nil {
		return
	}

	err = os.WriteFile(pemFile, pem, 0644)
	if err != nil {
		return
	}

	err = os.WriteFile(cerFile, cer, 0644)
	if err != nil {
		return
	}

	return
}

func CreateSelfSingedGraphAppCertificate(days time.Duration) (result *MsGraphCertificateResult, err error) {
	keys, cer, err := createSelfSingedGraphAppCertificate(days)
	if err != nil {
		return
	}
	result = &MsGraphCertificateResult{
		Keys: keys,
		Cert: cer,
	}
	return
}

func createSelfSingedGraphAppCertificate(days time.Duration) (keys, cer []byte, err error) {
	pKeys, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to generate private keys, error: %s", err)
	}

	serial, err := rand.Int(rand.Reader, (&big.Int{}).Exp(big.NewInt(2), big.NewInt(159), nil))
	if err != nil {
		return nil, nil, err
	}

	now := time.Now()
	template := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: "localhost"},
		NotBefore:             now.Add(-time.Minute).UTC(),
		NotAfter:              now.Add(days),
		BasicConstraintsValid: true,
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	}

	rawCertificate, err := x509.CreateCertificate(rand.Reader, template, template, &pKeys.PublicKey, pKeys)
	if err != nil {
		err = fmt.Errorf("failed to generate certificate, error: %s", err)
		return
	}

	certificate, err := x509.ParseCertificate(rawCertificate)
	if err != nil {
		err = fmt.Errorf("failed to parse certificate, error: %s", err)
		return
	}

	pfx, err := des.LegacyDES.Encode(pKeys, certificate, nil, "")
	if err != nil {
		err = fmt.Errorf("failed to generate pfx, error: %s", err)
		return
	}

	buf := bytes.Buffer{}
	pems, err := pkcs12.ToPEM(pfx, "")
	if err != nil {
		return
	}
	for _, p := range pems {
		p.Headers = nil
		buf.Write(pem.EncodeToMemory(p))
	}

	keys = buf.Bytes()
	cer = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: rawCertificate})

	return
}
