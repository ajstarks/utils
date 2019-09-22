package main

import (
	"flag"
	"math/rand"
	"os"
	"time"

	"github.com/ajstarks/deck/generate"
)

// hdeckle makes a horizontal deckled edge running from (x1, y1) to (x2, y2)
func hfill(deck *generate.Deck, x, y, width, height float64, color string, n int) {
	xp := make([]float64, n)
	yp := make([]float64, n)

	// left point
	xp[0] = x
	yp[0] = y

	// right point
	xp[n-2] = x + width
	yp[n-2] = y

	// back to the left
	xp[n-1] = x
	yp[n-1] = y

	xincr := width / float64(n)
	for i := 1; i <= n-3; i++ {
		xp[i] = xp[i-1] + xincr
		yp[i] = y + rand.Float64()*height
	}
	deck.Polygon(xp, yp, color)
}

func vfill(deck *generate.Deck, x, y, width, height float64, color string, n int) {
	xp := make([]float64, n)
	yp := make([]float64, n)

	// bottom point
	xp[0] = x
	yp[0] = y

	// top point
	xp[n-2] = x
	yp[n-2] = y + height

	// back to the bottom
	xp[n-1] = x
	yp[n-1] = y

	yincr := height / float64(n)
	for i := 1; i <= n-3; i++ {
		yp[i] = yp[i-1] + yincr
		xp[i] = x + rand.Float64()*width
	}
	deck.Polygon(xp, yp, color)
}

func hline(deck *generate.Deck, x, y, width, height, linewidth float64, color string, n int) {
	xincr := width / float64(n)
	hi := xincr / 2
	y1 := y
	for x1 := x; x1 < x+width; x1 += xincr {
		y2 := y1 + rand.Float64()*height
		deck.Line(x1, y1, x1+(hi), y2, linewidth, color)
		deck.Line(x1+(hi), y2, x1+xincr, y1, linewidth, color)
	}
}

func vline(deck *generate.Deck, x, y, width, height, linewidth float64, color string, n int) {
	yincr := height / float64(n)
	hi := yincr / 2
	x1 := x
	for y1 := y; y1 < y+height; y1 += yincr {
		x2 := x1 + rand.Float64()*width
		deck.Line(x1, y1, x2, y1+(hi), linewidth, color)
		deck.Line(x2, y1+(hi), x1, y1+yincr, linewidth, color)
	}
}

func main() {
	deck := generate.NewSlides(os.Stdout, 0, 0)
	var x, y, width, height, linewidth float64
	var n int
	var color, dtype string

	flag.Float64Var(&x, "x", 10, "x")
	flag.Float64Var(&y, "y", 50, "y")
	flag.Float64Var(&width, "w", 80, "width")
	flag.Float64Var(&height, "h", 3, "height")
	flag.Float64Var(&linewidth, "lw", 0.1, "line width")
	flag.IntVar(&n, "n", 50, "number of bumps")
	flag.StringVar(&color, "color", "gray", "color")
	flag.StringVar(&dtype, "type", "lh", "fv: filled vertical, fh: filled horizontal, lv: line vertical, lh: line horizontal")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	switch dtype {
	case "fv":
		vfill(deck, x, y, width, height, color, n)
	case "fh":
		hfill(deck, x, y, width, height, color, n)
	case "lv":
		vline(deck, x, y, width, height, linewidth, color, n)
	case "lh":
		hline(deck, x, y, width, height, linewidth, color, n)
	}
}
