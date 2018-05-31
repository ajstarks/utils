package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	port    = flag.String("port", ":8080", "http service address")
	webroot = flag.String("root", ".", "web root")
	logging = flag.Bool("log", false, "logging")
)

func main() {
	flag.Parse()
	http.Handle("/", http.HandlerFunc(sf))

	err := http.ListenAndServe(*port, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func sf(c http.ResponseWriter, req *http.Request) {
	if *logging {
		log.Printf("%#v\n", req.URL.Path)
	}
	if len(req.URL.Path) > 2 {
		http.ServeFile(c, req, req.URL.Path[1:])
	}
}
