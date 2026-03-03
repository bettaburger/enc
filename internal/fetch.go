package internal

import (
	"bufio"
	"enc"
	"fmt"

	//"io"
	"enc/config"
	"net/http"
	"os"
	"strings"
	"time"
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

func DisplayHelp(str string) {
	if enc.StrCmp(str, "--help") == 0 {
		fmt.Printf("\n")
		fmt.Printf("--help		show this message, then exit\n")
		fmt.Printf("\n")
		fmt.Printf("--exit		terminate process\n")
	}	
}

func Fetch(){
	//var filename string
	
	scanner := bufio.NewScanner(os.Stdin)
	/*TRANSPORT */
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		MaxIdleConnsPerHost: MAX_IDLE,		
		IdleConnTimeout:     0, // no limit, if the connection to the url is longer than the timeout?
		TLSHandshakeTimeout: 10 * time.Second,
		DisableCompression:  true,
	}

	client := &http.Client{
    	Transport: tr,
    	Timeout:   15 * time.Second,
	}

	fmt.Printf("\n")

	for {
		fmt.Print("$ https://")
		if !scanner.Scan() {
			break
		}
		url := strings.TrimSpace(scanner.Text())
		if enc.StrCmp(url, "--exit") == 0 {
			os.Exit(1)
		}
		DisplayHelp(url)
		url = HasHTTP(url)
		start := time.Now()
		response, err := client.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v. Try --help\n", err)
			continue
		}
		fmt.Printf("\n")
		response.Body.Close()
		r := config.BuildResponse(response)
		r.JSON()

		host, err := os.Hostname()
		if err != nil {
			fmt.Fprintf(os.Stderr, "hostname: %v\n", err)
		}
		/*
		file, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error retrieving %s\n", filename)
		}*/
		fmt.Printf("\n")
		fmt.Printf("%s %s\nResponse time: %dns %fs\nLocal: %s\n", response.Proto, response.Status, time.Since(start).Nanoseconds(), time.Since(start).Seconds(), host)
		fmt.Printf("\n")
		// CLOSE ANY idle connections, tr makes a max of 2 connections. 
		tr.CloseIdleConnections()
	}
}