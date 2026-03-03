package config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
    "os"
)

var (
    timeFormat = "Mon Jan 2 15:04:05 MST 2006"
)

type Date time.Time

type Header map[string][]string

type TLS struct {
    Version     uint16      `json:"version"`
    CipherSuite uint16      `json:"cipher_suite"`
    Server      string      `json:"server"`
    Issuer      string      `json:"issuer"`
    Org         []string    `json:"organization"`
    Holder      string      `json:"holder"`
    IssuedOn    Date        `json:"issued_on"`
    ExpiresOn   Date        `json:"expires_on"`
    ExpiresIn   string      `json:"expires_in"`
}

type Response struct {
    Headers     Header      `json:"headers"`
    URL         string      `json:"url"`
    TLS         *TLS        `json:"transport"`
}

func GetTLS(response *http.Response) *TLS {
    if response.TLS == nil {
        fmt.Println("No TLS (connection was HTTP)")
    }
    cert := response.TLS.PeerCertificates[0]
    var expiresIn string
    duration := time.Until(cert.NotAfter)
	if duration == 0 {
		expiresIn = "expired"
	} else {
		days := int(duration.Hours()) / 24
		hours := int(duration.Hours()) % 24
		minutes := int(duration.Minutes()) % 60
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

func BuildResponse(response *http.Response) *Response {
	if response == nil {
		return nil
	}
	return &Response{
		URL:     response.Request.URL.String(),
		Headers: Header(response.Header),
		TLS:     GetTLS(response),
	}
}

func (d Date) MarshalJSON() ([]byte, error) {
    str := time.Time(d).In(time.Local).Format(timeFormat)
    return []byte(`"` + str + `"`), nil
}









func (r *Response) JSON() {
	b, err := json.MarshalIndent(&r, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
	}
	os.Stdout.Write(b)
}


