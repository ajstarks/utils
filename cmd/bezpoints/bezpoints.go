// bezpoints generated decksh code for a filled quadratic bezier
package main

import (
	"flag"
	"fmt"
)

type point struct {
	x, y float64
}

// bezpoints computes the coordinates of a quadratic bezier curve
// source: https://en.wikipedia.org/wiki/B%C3%A9zier_curve
// for t between 0 and 1, where p0 is the start, p1 is conttol, p2 is end
// p_x = (1-t)^2*p0_x + 2(1-t)*t*p1_x + t^2*p2_x
// p_y = (1-t)^2*p0_y + 2(1-t)*t*p1_y + t^2*p2_y
func bezpoints(start, end, control point, npoints int) []point {
	p := make([]point, npoints)
	step := 1.0 / float64(npoints)
	t := 0.0
	for i := 0; i < npoints; i++ {
		p[i].x = (1-t)*(1-t)*start.x + 2*(1-t)*t*control.x + t*t*end.x
		p[i].y = (1-t)*(1-t)*start.y + 2*(1-t)*t*control.y + t*t*end.y
		t += step
	}
	return p
}

// showcurve generates the decksh code for the computed coordinates
func showcurve(start, end, control point, n int, color string, opacity float64) {
	coordinates := bezpoints(start, end, control, n)
	//fmt.Printf("curve %.3f %.3f %.3f %.3f %.3f %.3f 0.2 \"red\"\n", start.x, start.y, control.x, control.y, end.x, end.y)
	fmt.Printf("polygon \"%.3f ", start.x)
	for _, p := range coordinates {
		fmt.Printf("%.3f ", p.x)
	}
	fmt.Printf("%.3f\" \"%.3f ", end.x, start.y)
	for _, p := range coordinates {
		fmt.Printf("%.3f ", p.y)
	}
	fmt.Printf("%.3f\" \"%s\" %g\n", end.y, color, opacity)
}

func main() {
	var sx, sy, ex, ey, cx, cy, opacity float64
	var color string
	var npoints int

	flag.Float64Var(&sx, "sx", 20, "start x")
	flag.Float64Var(&sy, "sy", 50, "start y")
	flag.Float64Var(&cx, "cx", 50, "control x")
	flag.Float64Var(&cy, "cy", 80, "control y")
	flag.Float64Var(&ex, "ex", 80, "end x")
	flag.Float64Var(&ey, "ey", 50, "end y")
	flag.IntVar(&npoints, "n", 100, "number of points")
	flag.StringVar(&color, "color", "gray", "color")
	flag.Float64Var(&opacity, "opacity", 50, "opacity")
	flag.Parse()

	start := point{sx, sy}
	end := point{ex, ey}
	control := point{cx, cy}
	showcurve(start, end, control, npoints, color, opacity)
}
