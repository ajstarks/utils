// popio: import and export images for popi
package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"io"
	"os"
)

// popin: image to raw
func popin(w io.Writer, r io.Reader) (int, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return 0, err
	}
	// get image dimensions
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	// convert image pixels to grayscale
	data := make([]byte, width*height)
	i := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			data[i] = uint8((19595*r + 38470*g + 7471*b + 1<<15) >> 24)
			i++
		}
	}
	return bufio.NewWriter(w).Write(data)
}

// popout: raw to PNG
func popout(w io.Writer, r io.Reader, width, height int) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	// convert raw data to grayscale pixels
	img := image.NewGray(image.Rect(0, 0, width, height))
	i := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.Gray{data[i]})
			i++
		}
	}
	// write the png
	if err := png.Encode(w, img); err != nil {
		return err
	}
	return nil
}

func main() {
	var read, write bool
	var width, height int
	flag.BoolVar(&read, "import", false, "image to raw popi grayscale")
	flag.BoolVar(&write, "export", false, "popi raw grayscale to PNG")
	flag.IntVar(&width, "width", 512, "image width")
	flag.IntVar(&height, "height", 512, "image height")
	flag.Parse()

	if read && write {
		fmt.Fprintln(os.Stderr, "pick one: -import or -export")
		os.Exit(3)
	}

	if read {
		n, err := popin(os.Stdout, os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v (%d bytes written)\n", err, n)
			os.Exit(1)
		}
		os.Exit(0)
	}
	if write {
		err := popout(os.Stdout, os.Stdin, width, height)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(2)
		}
		os.Exit(0)
	}
}
