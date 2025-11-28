package pkg

import (
	"net/http"
	"fmt"
    "time"
)

type TLS struct {
    Version     uint16      `json:"version"`
    CipherSuite uint16      `json:"cipher_suite"`
    Server      string      `json:"server"`
// certification
    Issuer      string      `json:"issuer"`
    Org         []string    `json:"organization"`
    Holder      string      `json:"holder"`
    IssuedOn    time.Time   `json:"issued_on"`
    ExpiresOn   time.Time   `json:"expires_on"`
}


func DisplayTLS(response *http.Response) *TLS {
    if response.TLS == nil {
        fmt.Println("No TLS (connection was HTTP)")
    }

    cert := response.TLS.PeerCertificates[0]

	return &TLS{
		Version:      response.TLS.Version,
		CipherSuite:  response.TLS.CipherSuite,
		Server:   response.TLS.ServerName,
		Issuer:       cert.Issuer.CommonName,
		Org: cert.Issuer.Organization,
		Holder:      cert.Subject.CommonName,
		IssuedOn:    cert.NotBefore,
		ExpiresOn:   cert.NotAfter,
	}
}

