package main

import (
	"flag"
	"fmt"
)

func main() {
	var value, low1, high1, low2, high2 float64
	var format string

	flag.Float64Var(&value, "value", 1, "value")
	flag.Float64Var(&low1, "low1", 0, "low1")
	flag.Float64Var(&high1, "high1", 10, "high1")
	flag.Float64Var(&low2, "low2", 0, "low12")
	flag.Float64Var(&high2, "high2", 100, "high2")
	flag.StringVar(&format, "fmt", "%v", "format")
	flag.Parse()

	fmt.Printf(fmt.Sprintf("%s (%s, %s) (%s, %s) = %s\n", format, format, format, format, format, format), value, low1, high1, low2, high2, vmap(value, low1, high1, low2, high2))

}

// vmap maps one interval to another
func vmap(value float64, low1 float64, high1 float64, low2 float64, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}
