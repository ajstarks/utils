// imgps -- show GPS coordinates contained in EXIF data
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/rwcarlsen/goexif/exif"
)

// for every file on the command line, report GPS coordinates
func main() {
	for _, f := range os.Args[1:] {
		err := process(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", f, err)
			continue
		}
	}
}

// process retrieves GPS coordinates from a file
func process(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	x, err := exif.Decode(f)
	if err == io.EOF {
		return fmt.Errorf("no exif data")
	}
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
