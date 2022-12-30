package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// API Info
const (
	NYTAPIkey = "NYTAPIKEY" // obtained from the environment
	NYTfmt    = "http://api.nytimes.com/svc/news/v3/content/all/%s/.json?api-key=%s&limit=5"
)

// NYTHeadlines is the headline info from the New York Times
type NYTHeadlines struct {
	Status     string   `json:"status"`
	Copyright  string   `json:"copyright"`
	NumResults int      `json:"num_results"`
	Results    []result `json:"results"`
}

type result struct {
	Section    string `json:"section"`
	Subsection string `json:"subsection"`
	Title      string `json:"title"`
	Abstract   string `json:"abstract"`
	Thumbnail  string `json:"thumbnail_standard"`
}

func main() {
	var section = flag.String("h", "u.s.", "headline type (arts, health, sports, science, technology, u.s., world)")
	flag.Parse()
	nytheadlines(*section)
}

// apikey returns the API key from the environment, or the empty string if not found.
func apikey(s string) string {
	key, ok := os.LookupEnv(s)
	if !ok {
		return ""
	}
	return key
}

// netread derefernces a URL, returning the Reader, with an error
func netread(url string) (io.ReadCloser, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to get network data for %s (%s)", url, resp.Status)
	}
	return resp.Body, nil
}

// nytheadlines retrieves data from the New York Times API, decodes and displays it.
func nytheadlines(section string) {
	key := apikey(NYTAPIkey)
	if len(key) == 0 {
		fmt.Fprintln(os.Stderr, "invalid API key")
		return
	}
	r, err := netread(fmt.Sprintf(NYTfmt, section, key))
	if err != nil {
		fmt.Fprintf(os.Stderr, "headline read error: %v\n", err)
		return
	}
	defer r.Close()
	var data NYTHeadlines
	err = json.NewDecoder(r).Decode(&data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "decode: %v\n", err)
		return
	}
	for i := 0; i < len(data.Results); i++ {
		fmt.Printf("%v\n", data.Results[i].Title)
	}
}
