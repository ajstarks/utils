// mdtopdf -- convert markdown to PDF
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mandolyte/mdtopdf"
)

func die(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}

func main() {
	var input, output string
	flag.StringVar(&input, "i", "", "input markdown file (default is standard input)")
	flag.StringVar(&output, "o", "", "output PDF file (required)")
	flag.Parse()

	if output == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	var content []byte
	var err error

	if input == "" {
		content, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			die(err)
		}
	} else {
		content, err = os.ReadFile(input)
		if err != nil {
			die(err)
		}
	}
	pf := mdtopdf.NewPdfRenderer("", "", output, "")
	err = pf.Process(content)
	if err != nil {
		die(err)
	}
}
