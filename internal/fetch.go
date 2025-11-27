package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func hasHTTP(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}
	return url
}

func printTLSInfo(resp *http.Response) {
    if resp.TLS == nil {
        fmt.Println("No TLS (connection was HTTP)")
        return
    }

    cs := resp.TLS.CipherSuite
    ver := resp.TLS.Version
    cert := resp.TLS.PeerCertificates[0]

    fmt.Println("=== TLS Information ===")
    fmt.Printf("Version: %x\n", ver)
    fmt.Printf("Cipher Suite: %x\n", cs)
    fmt.Printf("Server Name: %s\n", resp.TLS.ServerName)
    fmt.Printf("Certificate Issuer: %s\n", cert.Issuer.CommonName)
    fmt.Printf("Certificate Subject: %s\n", cert.Subject.CommonName)
    fmt.Printf("Valid From: %s\n", cert.NotBefore)
    fmt.Printf("Valid Until: %s\n", cert.NotAfter)
}

func Fetch(){
	/*TRANSPORT */
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		MaxIdleConns:        20,
		IdleConnTimeout:     30 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		DisableCompression:  true,
	}

	client := &http.Client{
    	Transport: tr,
    	Timeout:   15 * time.Second,
	}

	for _, url := range os.Args[1:] {

		url = hasHTTP(url)

		response, err := client.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err) 
			os.Exit(1)
		}

		b, err := io.Copy(os.Stdout, response.Body)
		response.Body.Close()

		printTLSInfo(response)


		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err) 
			os.Exit(1)
		}

		fmt.Printf("%d\n", b)

		fmt.Printf("Response code: %s\n", response.Status)
	}
}