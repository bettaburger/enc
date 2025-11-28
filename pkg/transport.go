package pkg

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Date time.Time

var (
    cert        *x509.Certificate
    duration    time.Duration
    expiresIn   string
    days        int
    hours       int
    minutes     int
)

type TLS struct {
    Version     uint16      `json:"version"`
    CipherSuite uint16      `json:"cipher_suite"`
    Server      string      `json:"server"`
// certification
    Issuer      string      `json:"issuer"`
    Org         []string    `json:"organization"`
    Holder      string      `json:"holder"`
    IssuedOn    Date        `json:"issued_on"`
    ExpiresOn   Date        `json:"expires_on"`
    ExpiresIn   string      `json:"expires_in"`
}


func GetTLS(response *http.Response) *TLS {
    if response.TLS == nil {
        fmt.Println("No TLS (connection was HTTP)")
    }


    cert = response.TLS.PeerCertificates[0]
    
    duration = time.Until(cert.NotAfter)
	if duration == 0 {
		expiresIn = "expired"
	} else {
		days = int(duration.Hours()) / 24
		hours = int(duration.Hours()) % 24
		minutes = int(duration.Minutes()) % 60
		expiresIn = fmt.Sprintf("%d days, %d hours, %d minutes", days, hours, minutes)
	}

	return &TLS{
		Version:        response.TLS.Version,
		CipherSuite:    response.TLS.CipherSuite,
		Server:         response.TLS.ServerName,
		Issuer:         cert.Issuer.CommonName,
		Org:            cert.Issuer.Organization,
		Holder:         cert.Subject.CommonName,
		IssuedOn:       Date(cert.NotBefore),
		ExpiresOn:      Date(cert.NotAfter),
        ExpiresIn:      expiresIn,
	}
}

func (t *TLS) JSON() (string, error) {
    b, err := json.MarshalIndent(t, "", "   ") 
    if err != nil {
        return "", err
    }
    return string(b), nil
}

func (d Date) MarshalJSON() ([]byte, error) {
    str := time.Time(d).In(time.Local).Format("Mon Jan 2 15:04:05 MST 2006")
    return []byte(`"` + str + `"`), nil
}


