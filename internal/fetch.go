package internal

import (
	"fmt"
	//"io"
	"net/http"
	"os"
	"strings"
	"time"

	"enc/pkg"
)

var (
	MAX_IDLE = 2
)

func HasHTTP(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}
	return url
}

func Fetch(){
	/*TRANSPORT */
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		MaxIdleConnsPerHost: MAX_IDLE,		
		IdleConnTimeout:     30 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		DisableCompression:  true,
	}

	client := &http.Client{
    	Transport: tr,
    	Timeout:   15 * time.Second,
	}

	for _, url := range os.Args[1:] {
		
		url = HasHTTP(url)

		start := time.Now()

		response, err := client.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err) 
			os.Exit(1)
		}

		response.Body.Close()
		tls := pkg.GetTLS(response)

		s, _ := tls.JSON()
		fmt.Println(s) 

		fmt.Printf("Response code: %s\nResponse time: %dns %fs\n", response.Status, time.Since(start).Nanoseconds(), time.Since(start).Seconds())

		// CLOSE ANY idle connections, tr makes a max of 2 connections. 
		tr.CloseIdleConnections()
	}
}