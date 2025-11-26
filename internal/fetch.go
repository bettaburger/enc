package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func hasHTTP(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}
	return url
}

func Fetch(){
	for _, url := range os.Args[1:] {

		url = hasHTTP(url)

		response, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err) 
			os.Exit(1)
		}

		b, err := io.Copy(os.Stdout, response.Body)
		response.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err) 
			os.Exit(1)
		}

		fmt.Printf("%d\n", b)

		fmt.Printf("Response code: %s\n", response.Status)
	}
}