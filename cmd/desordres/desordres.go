// desordres -- tile blocks of lines as in Vera Molnar's Des Ordres
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// random returns a random number between a range
func random(min, max float64) float64 {
	return vmap(rand.Float64(), 0, 1, min, max)
}

// vmap maps one interval to another
func vmap(value, low1, high1, low2, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

// csquare makes a square with lines, using a specified width and color
// if a hue range is set, the color is randomly selected in that range,
// otherwise, the named color is used.
func csquare(x, y, size, maxlw, h1, h2 float64, color string) {
	lw := random(0.1, maxlw)
	if h1 > -1 && h2 > -1 { // hue range set
		color = fmt.Sprintf("hsv(%v,100,100)", random(h1, h2))
	}
	// define the corners
	hs := size / 2
	tlx, tly := x-hs, y+hs
	trx, try := x+hs, y+hs
	blx, bly := x-hs, y-hs
	brx, bry := x+hs, y-hs
	// make the boundaries
	hline(tlx, tly, size, lw, color)
	hline(blx, bly, size, lw, color)
	vline(blx, bly, size, lw, color)
	vline(brx, bry, size, lw, color)
	// make the corners
	square(tlx, tly, lw, color)
	square(blx, bly, lw, color)
	square(brx, bry, lw, color)
	square(trx, try, lw, color)
}

// square makes a square
func square(x, y, size float64, color string) {
	fmt.Printf("square %v %v %v \"%s\"\n", x, y, size, color)
}

// hline makes a horizontal line
func hline(x, y, size, lw float64, color string) {
	fmt.Printf("hline %v %v %v %v \"%s\"\n", x, y, size, lw, color)
}

// vline makes a vertical line
func vline(x, y, size, lw float64, color string) {
	fmt.Printf("vline %v %v %v %v \"%s\"\n", x, y, size, lw, color)
}

// desordres makes a series of concentric squares
func desordres(x, y, minsize, maxsize, maxlw, h1, h2 float64, color string) {
	step := random(1, 5)
	for v := minsize; v < maxsize; v += step {
		csquare(x, y, v, maxlw, h1, h2, color)
	}
}

// parseHues parses a color string: if the string is of the form "h1:h2",
// where h1, and h2 are numbers between 0 and 360, they are a range of hues.
// Otherwise, set to -1 for invalid entries (use named colors instead)
func parseHues(color string) (float64, float64) {
	var h1, h2 float64 = -1.0, -1.0
	hb := strings.Split(color, ":")
	if len(hb) == 2 {
		var err error
		h1, err = strconv.ParseFloat(hb[0], 64)
		if err != nil {
			h1 = -1
		}
		h2, err = strconv.ParseFloat(hb[1], 64)
		if err != nil {
			h2 = -1
		}
	}
	return h1, h2
}

// slide generation functions
func beginDeck()              { fmt.Println("deck") }
func endDeck()                { fmt.Println("edeck") }
func beginSlide(color string) { fmt.Printf("slide %q\n", color) }
func endSlide()               { fmt.Println("eslide") }

func main() {
	var tiles, maxlw float64
	var bgcolor, color string

	flag.Float64Var(&tiles, "tiles", 10, "tiles/row")
	flag.Float64Var(&maxlw, "maxlw", 1, "maximum line thickness")
	flag.StringVar(&bgcolor, "bgcolor", "white", "background color")
	flag.StringVar(&color, "color", "gray", "pen color")
	flag.Parse()

	size := 100 / tiles        // size of each tile
	top := 100 - (size / 2)    // top of the beginning row
	left := 100 - top          // left of the beginning row
	h1, h2 := parseHues(color) // set hue range, or named color

	beginDeck()
	beginSlide(bgcolor)
	for y := top; y > 0; y -= size {
		for x := left; x < 100; x += size {
			desordres(x, y, 2, size, maxlw, h1, h2, color)
		}
	}
	endSlide()
	endDeck()
}
