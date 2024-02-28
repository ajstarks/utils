// randgen makes random numbers
package main

import (
	"flag"
	"fmt"
	"math/rand"
)

func vmap(v, l1, h1, l2, h2 float64) float64 {
	return l2 + (h2-l2)*(v-l1)/(h1-l1)
}

func main() {
	nrand := flag.Int("n", 100, "number of items")
	min := flag.Float64("min", 0, "minimum value")
	max := flag.Float64("max", 1e6, "minimum value")
	ndec := flag.Int("dec", 3, "number of decimals")
	xint := flag.Float64("xint", 0, "x value interval")
	flag.Parse()
	xval := 0.0
	f := fmt.Sprintf("%%.%df", *ndec)
	for i := 0; i < *nrand; i++ {
		if *xint > 0 {
			fmt.Printf(f+"\t", xval)
			xval += *xint
		}
		fmt.Printf(f+"\n", vmap(rand.Float64(), 0, 1, *min, *max))
	}

}
