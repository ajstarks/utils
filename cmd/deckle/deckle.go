// deckle -- generate deck markup for deckled edges (lines and filled)
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/ajstarks/deck/generate"
)

func polygon(deck *generate.Deck, raw bool, x, y []float64, color string) {
	if raw {
		deck.Polygon(x, y, color)
	} else {
		fmt.Fprintf(os.Stdout, "polygon \"")
		for i := 0; i < len(x); i++ {
			fmt.Fprintf(os.Stdout, "%.4g ", x[i])
		}
		fmt.Fprintf(os.Stdout, "\" \"")
		for i := 0; i < len(y); i++ {
			fmt.Fprintf(os.Stdout, "%.4g ", y[i])
		}
		fmt.Fprintf(os.Stdout, "\" %q\n", color)
	}
}

func line(deck *generate.Deck, raw bool, x1, y1, x2, y2, linewidth float64, color string) {
	if raw {
		deck.Line(x1, y1, x2, y2, linewidth, color)
	} else {
		fmt.Fprintf(os.Stdout, "line %.4g %.4g %.4g %.4g %.4g %q\n", x1, y1, x2, y2, linewidth, color)
	}
}

// hfill makes a (width long) horizontal deckled edge starting at (x,y)
func hfill(deck *generate.Deck, mtype bool, x, y, width, height float64, color string, n int) {
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
	polygon(deck, mtype, xp, yp, color)
}

// vfill makes a (height high) vertical deckled edge
func vfill(deck *generate.Deck, mtype bool, x, y, width, height float64, color string, n int) {
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
	polygon(deck, mtype, xp, yp, color)
}

// hline makes a (width long) horizontal deckled edge
func hline(deck *generate.Deck, mtype bool, x, y, width, height, linewidth float64, color string, n int) {
	xincr := width / float64(n)
	hi := xincr / 2
	y1 := y
	for x1 := x; x1 < x+width; x1 += xincr {
		y2 := y1 + rand.Float64()*height
		line(deck, mtype, x1, y1, x1+(hi), y2, linewidth, color)
		line(deck, mtype, x1+(hi), y2, x1+xincr, y1, linewidth, color)
	}
}

// vline makes a (height high) vertical deckled edge
func vline(deck *generate.Deck, mtype bool, x, y, width, height, linewidth float64, color string, n int) {
	yincr := height / float64(n)
	hi := yincr / 2
	x1 := x
	for y1 := y; y1 < y+height; y1 += yincr {
		x2 := x1 + rand.Float64()*width
		line(deck, mtype, x1, y1, x2, y1+(hi), linewidth, color)
		line(deck, mtype, x2, y1+(hi), x1, y1+yincr, linewidth, color)
	}
}

func main() {
	deck := generate.NewSlides(os.Stdout, 0, 0)
	var (
		x, y, width, height, linewidth float64
		n                              int
		color, dtype                   string
		mtype                          bool
	)

	flag.Float64Var(&x, "x", 10, "x")
	flag.Float64Var(&y, "y", 50, "y")
	flag.Float64Var(&width, "w", 80, "width")
	flag.Float64Var(&height, "h", 3, "height")
	flag.Float64Var(&linewidth, "lw", 0.1, "line width")
	flag.BoolVar(&mtype, "raw", false, "type of markup - true for deck, false for decksh")
	flag.IntVar(&n, "n", 50, "number of bumps")
	flag.StringVar(&color, "color", "gray", "color")
	flag.StringVar(&dtype, "type", "lh", "fv: filled vertical, fh: filled horizontal, lv: line vertical, lh: line horizontal")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	switch dtype {
	case "fv":
		vfill(deck, mtype, x, y, width, height, color, n)
	case "fh":
		hfill(deck, mtype, x, y, width, height, color, n)
	case "lv":
		vline(deck, mtype, x, y, width, height, linewidth, color, n)
	case "lh":
		hline(deck, mtype, x, y, width, height, linewidth, color, n)
	}
}
