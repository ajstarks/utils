package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

func dpi(w, h, d float64) float64 {
	return math.Sqrt((w*w)+(h*h)) / d
}

func main() {
	if len(os.Args) < 4 {
		println("Usage: dpi w h diag [device]")
		os.Exit(1)
	}
	w, _ := strconv.ParseFloat(os.Args[1], 64)
	h, _ := strconv.ParseFloat(os.Args[2], 64)
	d, _ := strconv.ParseFloat(os.Args[3], 64)
	if len(os.Args) == 5 {
		fmt.Printf("%s ", os.Args[4])
	}
	if h > 0 && d > 0 {
		fmt.Printf("w=%.0f h=%0.f diag=%.2f aspect=%.2f dpi=%.2f\n", w, h, d, w/h, dpi(w, h, d))
	}
}
