package main

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/ajstarks/deck/generate"
)

// bar3d makes a 3D bar
func bar3d(deck *generate.Deck, x, y, w, h float64, tcolor, lcolor string) {
	wh := w / 2
	th := w * 0.5
	th2 := th / 2
	yh := y + h
	topx := []float64{x, x - wh, x, x + wh}
	topy := []float64{yh - th, yh - th2, yh, yh - th2}
	leftx := []float64{x, x - wh, x - wh, x}
	liney := []float64{y, y + th2, yh - th2, yh - th}
	rightx := []float64{x, x + wh, x + wh, x}

	deck.Polygon(topx, topy, tcolor)
	deck.Polygon(leftx, liney, lcolor)
	deck.Polygon(rightx, liney, lcolor, 60)
}

//vmap maps one interval to another
func vmap(value float64, low1 float64, high1 float64, low2 float64, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

// bardata reads data from the io.Reader, and plots bars
func bardata(deck *generate.Deck, r io.Reader, left, bottom, top float64, tcolor, lcolor string) {
	x := left
	y := bottom
	width := 5.0
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 2 {
			continue
		}
		value, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			continue
		}
		yp := vmap(value, 0, 200, 0, top-bottom)
		bar3d(deck, x, y, width, yp, tcolor, lcolor)
		deck.TextMid(x, y-2, fields[0], "sans", 1.5, "")
		x += width
	}
}

func main() {
	deck := generate.NewSlides(os.Stdout, 0, 0)
	deck.StartDeck()
	deck.StartSlide("rgb(30,10,10)", "linen")
	bardata(deck, os.Stdin, 20, 10, 80, "maroon", "linen")
	deck.EndSlide()
	deck.EndDeck()
}
