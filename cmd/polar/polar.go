package main

import (
	"flag"
	"fmt"
	"math"
)

func main() {
	var x, y, r, theta float64
	flag.Float64Var(&x, "x", 50, "x coordinate")
	flag.Float64Var(&y, "y", 50, "y coordinate")
	flag.Float64Var(&r, "r", 10, "radius")
	flag.Float64Var(&theta, "t", 90, "angle (degrees)")
	flag.Parse()

	rad := theta * (math.Pi / 180)
	fmt.Printf("x=%g y=%g r=%g t=%g -> %g %g\n", x, y, r, theta, x+(r*math.Cos(rad)), y+(r*math.Sin(rad)))
}
