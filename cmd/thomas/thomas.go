// Iris, Tulips, Jonquils, and Crocuses, 1969, Alma Thomas
package main

import (
	"fmt"
	"math/rand"
)

type brushStack struct {
	n         int     // number of items
	intensity int     // how many chips/item
	height    float64 // height in percent
	color     string  // the chip color
	opacity   float64 // chip  opacity
}

type brushStacks []brushStack

const (
	red1      = "#9b524e"
	red2      = "#8a2a1b"
	blue1     = "#204271"
	blue2     = "#68a8b2"
	blue3     = "#5c8bb9"
	blue4     = "#9dcbd3"
	pink      = "#c5abb0"
	yellow1   = "#e7d367"
	yellow2   = "#cba842"
	yellow3   = "#d3d9b3"
	orange    = "#b36735"
	violet1   = "#4b4b7a"
	violet2   = "#904561"
	violet3   = "#908eba"
	green1    = "#376769"
	green2    = "#c2cfbb"
	green3    = "#759899"
	bluegreen = "#62a5ab"
	defop     = 100.0
)

var alldata = []brushStacks{
	{
		{n: 3, color: green2, opacity: defop},
		{n: 33, color: blue1, opacity: defop},
	},
	{
		{n: 8, color: green2, opacity: defop},
		{n: 12, color: green1, opacity: defop},
		{n: 16, color: blue3, opacity: defop},
	},
	{
		{n: 10, color: green2, opacity: defop},
		{n: 26, color: blue1, opacity: defop},
	},
	{
		{n: 8, color: green2, opacity: defop},
		{n: 14, color: blue1, opacity: defop},
		{n: 7, color: blue3, opacity: defop},
		{n: 7, color: blue1, opacity: defop},
	},
	{
		{n: 9, color: green2, opacity: defop},
		{n: 27, color: blue1, opacity: defop},
	},
	{
		{n: 7, color: green2, opacity: defop},
		{n: 29, color: violet1, opacity: defop},
	},
	{
		{n: 21, color: yellow1, opacity: defop},
		{n: 15, color: blue1, opacity: defop},
	},
	{
		{n: 18, color: yellow1, opacity: defop},
		{n: 18, color: green2, opacity: defop},
	},
	{
		{n: 36, color: yellow2, opacity: defop},
	},
	{
		{n: 36, color: orange, opacity: defop},
	},
	{
		{n: 36, color: red2, opacity: defop},
	},
	{
		{n: 36, color: violet1, opacity: defop},
	},
	{
		{n: 3, color: yellow1, opacity: defop},
		{n: 5, color: green1, opacity: defop},
		{n: 7, color: violet1, opacity: defop},
		{n: 10, color: blue1, opacity: defop},
		{n: 11, color: violet1, opacity: defop},
	},
	{
		{n: 4, color: yellow1, opacity: defop},
		{n: 6, color: green1, opacity: defop},
		{n: 21, color: blue1, opacity: defop},
		{n: 5, color: green1, opacity: defop},
	},
	{
		{n: 32, color: blue1, opacity: defop},
		{n: 4, color: blue2, opacity: defop},
	},
	{
		{n: 33, color: blue1, opacity: defop},
		{n: 3, color: blue2, opacity: defop},
	},
	{
		{n: 9, color: red1, opacity: defop},
		{n: 27, color: blue1, opacity: defop},
	},
	{
		{n: 10, color: yellow1, opacity: defop},
		{n: 26, color: blue1, opacity: defop},
	},
	{
		{n: 4, color: yellow1, opacity: defop},
		{n: 32, color: blue1, opacity: defop},
	},
	{
		{n: 5, color: yellow1, opacity: defop},
		{n: 31, color: green1, opacity: defop},
	},
	{
		{n: 36, color: yellow1, opacity: defop},
	},
	{
		{n: 36, color: red1, opacity: defop},
	},
	{
		{n: 36, color: red1, opacity: defop},
	},
	{
		{n: 36, color: blue2, opacity: defop},
	},
	{
		{n: 36, color: red1, opacity: defop},
	},
	{
		{n: 36, color: red1, opacity: defop},
	},
	{
		{n: 36, color: orange, opacity: defop},
	},
	{
		{n: 36, color: green1, opacity: defop},
	},
	{
		{n: 36, color: blue1, opacity: defop},
	},
	{
		{n: 36, color: blue1, opacity: defop},
	},
	{
		{n: 36, color: blue1, opacity: defop},
	},
	{
		{n: 28, color: violet2, opacity: defop},
		{n: 8, color: blue1, opacity: defop},
	},
	{
		{n: 11, color: orange, opacity: defop},
		{n: 25, color: blue1, opacity: defop},
	},
	{
		{n: 14, color: green2, opacity: defop},
		{n: 22, color: blue1, opacity: defop},
	},
	{
		{n: 36, color: green2, opacity: defop},
	},
	{
		{n: 36, color: yellow1, opacity: defop},
	},
	{
		{n: 36, color: red1, opacity: defop},
	},
	{
		{n: 36, color: pink, opacity: defop},
	},
	{
		{n: 36, color: blue1, opacity: defop},
	},
	{
		{n: 36, color: blue1, opacity: defop},
	},
	{
		{n: 36, color: blue4, opacity: defop},
	},
	{
		{n: 36, color: yellow3, opacity: defop},
	},
	{
		{n: 36, color: blue1, opacity: defop},
	},
	{
		{n: 6, color: red1, opacity: defop},
		{n: 30, color: blue1, opacity: defop},
	},
	{
		{n: 15, color: red1, opacity: defop},
		{n: 21, color: blue1, opacity: defop},
	},
	{
		{n: 17, color: red1, opacity: defop},
		{n: 19, color: violet3, opacity: defop},
	},
	{
		{n: 18, color: red1, opacity: defop},
		{n: 18, color: bluegreen, opacity: defop},
	},
	{
		{n: 15, color: pink, opacity: defop},
		{n: 21, color: blue1, opacity: defop},
	},
	{
		{n: 16, color: pink, opacity: defop},
		{n: 20, color: green3, opacity: defop},
	},
}

// chip makes a four-sided polygon bounded to the rectangle
// defined by (x, y) at its center with dimensions (w,h).
// wd defines the depth inside the left and right sides,
// hd defines the depth inside the top and bottom
// The effect is a brush Stack when the coordinates of the polygon are
// defined randomly within these constraints.
func chip(x, y, w, h, wd, hd float64, color string, opacity float64) {
	xp := make([]float64, 4)
	yp := make([]float64, 4)

	l := x - (w / 2)
	r := x + (w / 2)
	t := y + (h / 2)
	b := y - (h / 2)

	xp[0] = l + (w * rand.Float64())
	yp[0] = t - (hd * rand.Float64())

	xp[1] = r - (wd * rand.Float64())
	yp[1] = b + (h * rand.Float64())

	xp[2] = l + (w * rand.Float64())
	yp[2] = b + (hd * rand.Float64())

	xp[3] = l + (wd * rand.Float64())
	yp[3] = b + (h * rand.Float64())

	fmt.Printf("polygon \"%.2f %.2f %.2f %.2f\" \"%.2f %.2f %.2f %.2f\" \"%s\" %.2f\n",
		xp[0], xp[1], xp[2], xp[3], yp[0], yp[1], yp[2], yp[3], color, opacity)
}

// blob makes n number of ellipses bounded by the rectangle
// centered at (x,y) with dimensions (w,h)
// the width and height of the ellipses are determined as
// some fraction of (w,h) determined by (wh, hd)
func blob(x, y, w, h, wd, hd float64, n int, color string, opacity float64) {
	l := x - (w / 2)
	b := y - (h / 2)
	//fmt.Printf("rect %.2f %.2f %.2f %.2f \"%s\" %.2f\n", x, y, w, h, "black", 10.0)
	for i := 0; i < n; i++ {
		xp, yp := l+(w*rand.Float64()), b+(h*rand.Float64())
		ew, eh := w*wd, h*hd
		fmt.Printf("ellipse %.2f %.2f %.2f %.2f \"%s\" %.2f\n", xp, yp, ew, eh, color, opacity)
	}
}

func tower(data brushStacks, x, y, w, h, wd, hd float64, nb int) {
	yp := y
	for i := 0; i < len(data); i++ {
		d := data[i]
		for j := 0; j < d.n; j++ {
			for n := 0; n < nb; n++ {
				chip(x, yp, w, h, wd, hd, d.color, d.opacity)
			}
			yp += h
		}
	}
}

func blobtower(data brushStacks, x, y, w, h, wd, hd float64, nb int) {
	yp := y
	for i := 0; i < len(data); i++ {
		d := data[i]
		for j := 0; j < d.n; j++ {
			blob(x, yp, w, h, wd, hd, nb, d.color, d.opacity)
			yp += h
		}
	}
}

func alltower(data []brushStacks, x, y, w, h, wd, hd float64, nb int) {
	xp := x
	for _, f := range data {
		tower(f, xp, y, w, h, wd, hd, nb)
		xp += w
	}
}

func grid(x1, x2, y1, y2, w, h, wd, hd float64, nb int, palette []string) {
	for x := x1; x <= x2; x += w {
		for y := y1; y <= y2; y += h {
			for n := 0; n < nb; n++ {
				c := rand.Intn(len(palette))
				blob(x, y, w, h, wd, hd, 1, palette[c], 10)
			}
		}
	}
}

func main() {
	bgcolor := "white"
	fmt.Println("deck")
	fmt.Printf("slide \"%s\"\n", bgcolor)
	w := 1.5
	h := w * 1.6
	wd := w * 0.05
	hd := h * 0.05
	alltower(alldata, 15, 10, w, h, wd, hd, 1)
	fmt.Println("eslide")

	// fmt.Printf("slide \"%s\"\n", bgcolor)
	// blob(50, 50, 10, 15, 1, 1, 10, red1, 60)
	// fmt.Println("eslide")

	// fmt.Printf("slide \"%s\"\n", bgcolor)
	// grid(20, 80, 20, 80, w, h, 1, 1, 10, []string{"red", "green", "blue", "pink", "violet", "yellow", "orange"})
	// fmt.Println("eslide")

	fmt.Println("edeck")

}
