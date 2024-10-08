package certificates

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"strings"
	"time"
)

type TlsCertificateResult struct {
	PrivateKey  []byte
	Certificate []byte

	Accepted []string
}

func CreateSelfSignedTlsCertificate(acceptedHostNames []string, days int) (data *TlsCertificateResult, err error) {
	var dnsnames []string
	var ips []net.IP

	if len(acceptedHostNames) == 0 {
		dnsnames, ips, err = getDnsNames()
		if err != nil {
			return nil, err
		}
	} else {
		for _, name := range acceptedHostNames {
			if ip := net.ParseIP(name); ip != nil {
				ips = append(ips, ip)
			} else {
				dnsnames = append(dnsnames, name)
			}
		}
	}

	data, err = createCertificateAuthority(
		dnsnames,
		ips,
		pkix.Name{
			Country:            []string{"DE"},
			Organization:       []string{"intex"},
			OrganizationalUnit: []string{"research"},
			Locality:           []string{"Lauf"},
			Province:           []string{"BY", "FR"},
			CommonName:         "intex.software",
			StreetAddress:      []string{"Am Winkelsteig 1a"},
			PostalCode:         []string{"91207"},
		},
		days,
		4096,
	)
	return
}

func getIPs() (netIps []net.IP, err error) {
	ifaces, err := net.Interfaces()
	for _, i := range ifaces {
		if addrs, e := i.Addrs(); e != nil {
			return nil, e
		} else {
			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}
				netIps = append(netIps, ip)
			}
		}
	}
	return
}

func getDnsNames() (dnsnames []string, addrs []net.IP, err error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, nil, err
	}

	nameSet := map[string]struct{}{
		hostname: {},
	}

	addrs, err = getIPs()
	if err != nil {
		return nil, nil, err
	}
	for _, ip := range addrs {
		names, err := net.LookupAddr(ip.String())
		if err != nil {
			return nil, nil, err
		}
		for _, name := range names {
			name := strings.TrimSuffix(name, ".")
			nameSet[name] = struct{}{}
		}
	}

	for name := range nameSet {
		dnsnames = append(dnsnames, name)
	}

	return
}

func createCertificateAuthority(dnsnames []string, ips []net.IP, names pkix.Name, days int, size int) (*TlsCertificateResult, error) {
	keys, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("unable to generate private keys, error: %s", err)
	}

	serial, err := rand.Int(rand.Reader, (&big.Int{}).Exp(big.NewInt(2), big.NewInt(159), nil))
	if err != nil {
		return nil, err
	}

	now := time.Now()
	template := x509.Certificate{
		SerialNumber:          serial,
		Subject:               names,
		NotBefore:             now.Add(-time.Minute).UTC(),
		NotAfter:              now.Add(time.Duration(days) * Day),
		BasicConstraintsValid: true,
		IsCA:                  true,
		DNSNames:              dnsnames,
		IPAddresses:           ips,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	}

	rawCertificate, err := x509.CreateCertificate(rand.Reader, &template, &template, &keys.PublicKey, keys)
	if err != nil {
		return nil, fmt.Errorf("failed to generate certificate, error: %s", err)
	}

	cert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: rawCertificate})

	var pKey []byte
	if keyData, err := x509.MarshalPKCS8PrivateKey(keys); err == nil {
		pKey = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyData})
	} else if err != nil {
		return nil, err
	}

	accepted := []string{}
	accepted = append(accepted, dnsnames...)
	for _, ip := range ips {
		accepted = append(accepted, ip.String())
	}

	return &TlsCertificateResult{
		PrivateKey:  pKey,
		Certificate: cert,
		Accepted:    accepted,
	}, nil
}
