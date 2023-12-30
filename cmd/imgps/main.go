// imgps -- show GPS coordinates contained in EXIF data
package main

import (
	"fmt"
	"os"

	"github.com/rwcarlsen/goexif/exif"
)

func main() {
	for _, f := range os.Args[1:] {
		err := process(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", f, err)
			continue
		}
	}
}

func process(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	x, err := exif.Decode(f)
	if err != nil {
		return err
	}
	lat, long, err := x.LatLong()
	if err != nil {
		return err
	}
	fmt.Printf("%s %.8f %.8f\n", filename, lat, long)
	return nil
}
