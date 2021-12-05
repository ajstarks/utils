package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	port    = flag.String("port", ":8080", "http service address")
	logging = flag.Bool("log", false, "logging")
	cert    = flag.String("cert", "", "cert path")
	key     = flag.String("key", "", "key path")
)

func main() {
	flag.Parse()
	http.Handle("/", http.HandlerFunc(sf))
	log.Printf("ws: listen on %s", *port)

	var err error
	if *cert != "" && *key != "" {
		err = http.ListenAndServeTLS(*port, *cert, *key, nil)
	} else {
		err = http.ListenAndServe(*port, nil)
	}
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func sf(c http.ResponseWriter, req *http.Request) {
	if *logging {
		log.Printf("%v %v %v\n", req.RemoteAddr, req.URL.Path, req.UserAgent())
	}
	if len(req.URL.Path) > 2 {
		http.ServeFile(c, req, req.URL.Path[1:])
	}
}
