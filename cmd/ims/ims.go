package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func main() {
	for _, fname := range os.Args[1:] {
		f, ferr := os.Open(fname)
		if ferr != nil {
			fmt.Fprintf(os.Stderr, "%v\n", ferr)
			continue
		}
		im, _, err := image.DecodeConfig(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", fname, err)
			continue
		}
		fmt.Printf("%s %d %d\n", fname, im.Width, im.Height)
		f.Close()

	}
}
