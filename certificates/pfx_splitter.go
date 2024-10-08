package certificates

import (
	"golang.org/x/crypto/pkcs12"

	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func PfxSplitter(pfxFile, password string) (err error) {
	if _, err = os.Stat(pfxFile); errors.Is(err, os.ErrNotExist) {
		return
	}

	content, err := os.ReadFile(pfxFile)
	if err != nil {
		return
	}

	base := filepath.Join(
		filepath.Dir(pfxFile),
		strings.TrimSuffix(filepath.Base(pfxFile), filepath.Ext(pfxFile)),
	)

	pems, err := pkcs12.ToPEM(content, password)
	if err != nil {
		return
	}

	for i, p := range pems {
		p.Headers = nil
		os.WriteFile(fmt.Sprintf("%s-%d-%s.pem", base, 1+i, p.Type), pem.EncodeToMemory(p), 0644)
	}

	return
}
