// gurl - get url
package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// netread derefernces a URL, returning the Reader, with an error
func netread(url string, timeout int) (io.ReadCloser, error) {
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to get network data for %s (%s)", url, resp.Status)
	}
	return resp.Body, nil
}

func main() {
	timeout := flag.Int("timeout", 30, "time out (sec)")
	flag.Parse()
	for _, url := range flag.Args() {
		r, err := netread(url, *timeout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		io.Copy(os.Stdout, r)
		r.Close()
	}
}
