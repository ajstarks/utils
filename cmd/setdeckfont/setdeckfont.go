package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/flopp/go-findfont"
)

func main() {
	sample := ""
	if len(os.Args) > 1 {
		sample = os.Args[1]
	}
	fontPath, err := findfont.Find(sample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	fmt.Printf("DECKFONTS=%s\n", filepath.Dir(fontPath))

}
