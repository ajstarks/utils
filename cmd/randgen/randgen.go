package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

func vmap(v, l1, h1, l2, h2 float64) float64 {
	return l2 + (h2-l2)*(v-l1)/(h1-l1)
}

func main() {
	nrand := flag.Int("n", 100, "number of items")
	min := flag.Float64("min", 0, "minimum value")
	max := flag.Float64("max", 1e6, "minimum value")
	ndec := flag.Int("dec", 3, "number of decimals")
	flag.Parse()
	rand.Seed(int64(time.Now().Nanosecond()) % 1e9)
	f := fmt.Sprintf("%%.%df\n", *ndec)
	for i := 0; i < *nrand; i++ {
		fmt.Printf(f, vmap(rand.Float64(), 0, 1, *min, *max))
	}

}
