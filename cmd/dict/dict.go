package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/net/dict"
)

var (
	db      = flag.String("d", "wn", "Dictionary database")
	dserver = flag.String("s", "dict.org:2628", "Dictionary Server")
)

func main() {
	flag.Parse()
	c, err := dict.Dial("tcp", *dserver)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer c.Close()

	// no args, list dictionaries, exit
	if len(flag.Args()) == 0 {
		dicts, err := c.Dicts()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}
		for _, dl := range dicts {
			fmt.Println(dl.Name, dl.Desc)
		}
		return
	}

	// define each word specified on the command line
	for _, word := range flag.Args() {
		defs, err := c.Define(*db, word)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", word, err)
			continue
		}
		for _, result := range defs {
			fmt.Println(string(result.Text))
		}
	}
}
