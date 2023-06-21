// cc -- concentric circle designs
package main

import (
	"fmt"
	"math"
)

// circle draws a circle
func circle(x, y, size float64, color string) {
	fmt.Printf("circle %.2f %.2f %.2f %q\n", x, y, size, color)
}

// cpolar places a circle at a polar coordinate
func cpolar(x, y, r, t, size float64, color string) {
	fmt.Printf("p=polar %.2f %.2f %.2f %.2f\ncircle p_x  p_y %.2f %q\n", x, y, r, t, size, color)
}

// polar returns Cartiesian coordinates from polar
func polar(x, y, r, deg float64) (float64, float64) {
	theta := deg * (math.Pi / 180)
	px := x + (r * math.Cos(theta))
	py := y + (r * math.Sin(theta))
	return px, py
}

// begindeck begins decksh markup
func begindeck() {
	fmt.Println("deck\ncanvas 1000 1000")
}

// enddeck end a deck
func enddeck() {
	fmt.Println("edeck")
}

// beginslide begins a slide
func beginslide(color string) {
	fmt.Printf("slide \"%s\"\n", color)
}

// endslide ends a slide
func endslide() {
	fmt.Println("eslide")
}

// planet makes circles around a point
func planet(x, y, size, radius, a1, a2, steps float64, color string) {
	for t := a1; t < a2; t += steps {
		cpolar(x, y, radius, t, size, color)
	}
}

// solar makes circles around a central circle
func solar(x, y, csize, radius, psize, steps float64, ccolor, pcolor string) {
	circle(x, y, csize, ccolor)
	planet(x, y, psize, radius, 0, 360, steps, pcolor)
}

// d1 maes two concentric rings
func d1(step float64) {
	beginslide("black")
	circle(50, 50, step, "red")
	for t := 0.0; t <= 360; t += step {
		px, py := polar(50, 50, 25, t)
		solar(px, py, 5, 5, 1, step, "red", "orange")
	}
	for t := step / 2; t <= 360; t += step {
		px, py := polar(50, 50, 40, t)
		solar(px, py, 5, 5, 1, step, "red", "orange")
	}
	endslide()
}

// hsv specifies a hue value in the hsv color space
func hsv(hue, sat, value int) string {
	return fmt.Sprintf("hsv(%d,%d,%d)", hue, sat, value)
}

// cchue makes a series of 7 concentric rings, varying bu nue
func cchue(r, step float64, starthue int, bgcolor string) {
	beginslide(bgcolor)
	cstep := 1.0
	c := 1.0
	halfstep := step / 2
	csize := r * 1.5
	hue := starthue

	circle(50, 50, csize, hsv(hue, 100, 100))
	planet(50, 50, c, r, 0, 360, step, hsv(hue, 100, 100))
	r += 2
	c += cstep
	hue += 7
	planet(50, 50, c, r, halfstep, 360, step, hsv(hue, 100, 100))
	r += 3
	c += cstep
	hue += 7
	planet(50, 50, c, r, 0, 360, step, hsv(hue, 100, 100))
	r += 4
	c += cstep
	hue += 7
	planet(50, 50, c, r, halfstep, 360, step, hsv(hue, 100, 100))
	r += 5
	c += cstep
	hue += 7
	planet(50, 50, c, r, 0, 360, step, hsv(hue, 100, 100))
	r += 6
	c += cstep
	hue += 7
	planet(50, 50, c, r, halfstep, 360, step, hsv(hue, 100, 100))
	r += 8
	hue += 7
	planet(50, 50, 7, r, 0, 360, step, hsv(hue, 100, 100))

	for t := 0.0; t <= 360; t += step {
		px, py := polar(50, 50, r, t)
		planet(px, py, 1, 5, 0, 360, 30, hsv(starthue, 100, 100))
	}
	endslide()
}

func main() {
	begindeck()
	cchue(10, 20, 0, "black")
	d1(30)
	enddeck()
}
