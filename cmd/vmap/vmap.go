package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 6 {
		fmt.Fprintf(os.Stderr, "vmap value low1 high1 low2 high2\n")
		os.Exit(1)
	}
	var value, low1, high1, low2, high2 float64
	value, _ = strconv.ParseFloat(os.Args[1], 64)
	low1, _ = strconv.ParseFloat(os.Args[2], 64)
	high1, _ = strconv.ParseFloat(os.Args[3], 64)
	low2, _ = strconv.ParseFloat(os.Args[4], 64)
	high2, _ = strconv.ParseFloat(os.Args[5], 64)
	fmt.Printf("%g (%g, %g) (%g, %g) = %g\n", value, low1, high1, low2, high2, vmap(value, low1, high1, low2, high2))

}

//vmap maps one interval to another
func vmap(value float64, low1 float64, high1 float64, low2 float64, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}
