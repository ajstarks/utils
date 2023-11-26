// dotspiral -- concentric circle designs
package main

import (
	"flag"
	"fmt"
)

const (
	deckfmt   = "deck\ncanvas %d %d\n"
	circlefmt = "circle %.2f %.2f %.2f %q\n"
	polarfmt  = "p=polar %.2f %.2f %.2f %.2f\ncircle p_x  p_y %.2f %q %.2f\n"
	slidefmt  = "slide \"%s\"\n"
)

type config struct {
	start, end, r, rincr, dincr, tincr, dotsize, dotop float64
	dotcolor, bgcolor                                  string
}

// circle draws a circle
func circle(x, y, size float64, color string) {
	fmt.Printf(circlefmt, x, y, size, color)
}

// cpolar places a circle at a polar coordinate
func cpolar(x, y, r, t, size float64, color string, op float64) {
	fmt.Printf(polarfmt, x, y, r, t, size, color, op)
}

// begindeck begins decksh markup
func begindeck(w, h int) {
	fmt.Printf(deckfmt, w, h)
}

// enddeck end a deck
func enddeck() {
	fmt.Println("edeck")
}

// beginslide begins a slide
func beginslide(color string) {
	fmt.Printf(slidefmt, color)
}

// endslide ends a slide
func endslide() {
	fmt.Println("eslide")
}

// dotspiral makes a dot spiral
func dotspiral(cx, cy float64, c config) {
	r := c.r
	dotsize := c.dotsize
	beginslide(c.bgcolor)
	for t := c.start; t <= c.end; t += c.tincr {
		cpolar(cx, cy, r, t, dotsize, c.dotcolor, c.dotop)
		r += c.rincr
		dotsize += c.dincr
	}
	endslide()
}

// configure set command line options
func configure() config {
	var c config
	flag.Float64Var(&c.start, "start", 180, "start angle")
	flag.Float64Var(&c.end, "end", 360, "end angle")
	flag.Float64Var(&c.r, "r", 10.0, "radius")
	flag.Float64Var(&c.rincr, "rincr", 1.0, "radius increment")
	flag.Float64Var(&c.tincr, "tincr", 10.0, "angle increment")
	flag.Float64Var(&c.dincr, "dincr", 0.5, "size increment")
	flag.Float64Var(&c.dotsize, "size", 0.5, "dot size")
	flag.Float64Var(&c.dotop, "op", 50, "dot opacity")
	flag.StringVar(&c.dotcolor, "color", "red", "dot color")
	flag.StringVar(&c.bgcolor, "bgcolor", "white", "background color")
	flag.Parse()
	return c

}

func main() {
	begindeck(500, 500)
	dotspiral(50, 50, configure())
	enddeck()
}
