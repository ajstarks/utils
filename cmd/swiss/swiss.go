// swiss: make a Swiss Railway clock
// $ swiss | pdf -pagesize 500,500
package main

import (
	"math"
	"os"
	"time"

	"github.com/ajstarks/deck/generate"
)

var hrangles = [12]float64{
	90, 60, 30, // 12, 1, 2
	0, 330, 300, // 3, 4, 5
	270, 240, 210, // 6, 7, 8
	180, 150, 120, // 9, 10, 11
}

var minangles = [60]float64{
	90, 84, 78, 72, 66, // 0 - 4 min
	60, 54, 48, 42, 36, // 5 - 9 min
	30, 24, 18, 12, 6, // 10 - 14 min
	0, 354, 348, 342, 336, // 15 - 19 min
	330, 324, 318, 312, 306, // 20 - 24
	300, 294, 288, 282, 276, // 25 - 29
	270, 264, 258, 252, 246, // 30 - 34
	240, 234, 228, 222, 216, // 35 - 39
	210, 204, 198, 192, 186, // 40 - 44
	180, 174, 168, 162, 156, // 45 - 49
	150, 144, 138, 132, 126, // 50 - 54
	120, 114, 108, 102, 96, // 55 - 59
}

// polar converts from polar to cartesian coordinates
func polar(cx, cy, r, degrees float64) (float64, float64) {
	a := degrees * (math.Pi / 180)
	x := r * math.Cos(a)
	y := r * math.Sin(a)
	return x + cx, y + cy
}

// clock draws a clock face with hour and minute markers
func clock(deck *generate.Deck, x, y, r float64) {
	n := 0
	var r2, s2 float64

	// scale the dimensions of the hour and minute ticks
	// by the size of the clock face
	hsize := r / 5
	msize := hsize / 3

	hlen := r * 0.05
	mlen := hlen / 3

	deck.Circle(x, y, (r*2)+msize, "silver")
	deck.Circle(x, y, r*2, "white")

	// Around the circle, make hour and second ticks
	// for every 5, mark the hour, else mark the minute
	for t := 0.0; t < 360; t += 6 {
		if n%5 == 0 { // hours
			r2 = r - hsize
			s2 = hlen
		} else { // seconds
			r2 = r - msize
			s2 = mlen
		}
		n++
		px1, py1 := polar(x, y, r, t)
		px2, py2 := polar(x, y, r2, t)
		deck.Line(px1, py1, px2, py2, s2, "black")
	}
}

// oppangle computes the opposite angle
func oppangle(a float64) float64 {
	if a >= 0 && a <= 180 {
		return a + 180
	}
	return a - 180
}

// drawtime draws a clock face with  hour, minute and second hands
func drawtime(deck *generate.Deck, x, y, r float64, h, m, s int) {
	if (m > 59 || m < 0) || (s > 59 || s < 0) {
		return
	}
	clock(deck, x, y, r)
	linesize := r * 0.09
	extrar := r / 4 // length of the line past the centerline

	// get angles for hour, minute, second
	ha := hrangles[h%12]
	ma := minangles[m]
	sa := minangles[s]

	if m > 30 { // if the minute is > 30, adjust the hour angle proportionally
		ha = ha - (30.0 * float64(m) / 60.0)
	}

	// the hour, minute, and second hands are drawn in two parts.
	// part 1 is the line between the center point (x,y) and radius
	// part 2 is the extra line past the center whose angle is opposite the part 1 line.

	// hour line
	hx, hy := polar(x, y, r*0.7, ha)
	hx2, hy2 := polar(x, y, extrar, oppangle(ha))
	deck.Line(x, y, hx, hy, linesize, "gray")
	deck.Line(x, y, hx2, hy2, linesize, "gray")

	// minute line
	mx, my := polar(x, y, r*0.95, ma)
	mx2, my2 := polar(x, y, extrar, oppangle(ma))
	deck.Line(x, y, mx, my, linesize, "black")
	deck.Line(x, y, mx2, my2, linesize, "black")

	// second line -- includes dot at the end
	dotsize := r * 0.2
	slinesize := linesize * 0.275
	sx, sy := polar(x, y, (r*0.7)-dotsize/2, sa)
	sx2, sy2 := polar(x, y, extrar, oppangle(sa))
	cx, cy := polar(x, y, r*0.7, sa)
	deck.Line(x, y, sx, sy, slinesize, "red")
	deck.Line(x, y, sx2, sy2, slinesize, "red")
	deck.Circle(cx, cy, dotsize, "red")
}

func main() {
	now := time.Now()
	deck := generate.NewSlides(os.Stdout, 0, 0)
	deck.StartDeck()
	deck.StartSlide()
	deck.TextMid(50, 2, now.Format(time.Kitchen), "sans", 4, "")
	drawtime(deck, 50, 50, 40, now.Hour(), now.Minute(), now.Second())
	deck.EndSlide()

	deck.StartSlide()
	deck.TextMid(20, 30, "Now-3H", "sans", 2, "")
	deck.TextMid(50, 30, "Now", "sans", 2, "")
	deck.TextMid(80, 30, "Now+5H", "sans", 2, "")

	drawtime(deck, 20, 50, 10, now.Hour()-3, now.Minute(), now.Second())
	drawtime(deck, 50, 50, 10, now.Hour(), now.Minute(), now.Second())
	drawtime(deck, 80, 50, 10, now.Hour()+5, now.Minute(), now.Second())
	deck.EndSlide()

	deck.StartSlide("black", "white")
	m := 0
	s := 5
	for y := 80.0; y >= 20; y -= 20 {
		for x := 20.0; x <= 80; x += 30 {
			drawtime(deck, x, y, 8, 12, m, s%60)
			m += 5
			s += 5
		}
	}
	deck.EndSlide()
	deck.EndDeck()
}
