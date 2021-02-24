// csv2poly - generate decksh polygons from x,y pairs
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	smallest = -math.MaxFloat64
	largest  = math.MaxFloat64
)

type params struct {
	left, right, bottom, top float64
	label, color             string
}

func main() {
	var p params
	flag.Float64Var(&p.left, "left", 0, "left")
	flag.Float64Var(&p.right, "right", 100, "right")
	flag.Float64Var(&p.bottom, "bottom", 0, "bottom")
	flag.Float64Var(&p.top, "top", 100, "top")
	flag.StringVar(&p.color, "color", "gray", "color")
	flag.StringVar(&p.label, "label", "", "label")
	flag.Parse()

	for _, f := range flag.Args() {
		if err := process(p, f); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
	}
}

// readata reads x, y pairs, checking for errors
func readata(r io.Reader) ([]float64, []float64, error) {
	var x, y []float64
	var xp, yp float64
	var err error
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), ",")
		if len(fields) != 2 {
			continue
		}
		if xp, err = strconv.ParseFloat(fields[0], 64); err != nil {
			continue
		}
		if yp, err = strconv.ParseFloat(fields[1], 64); err != nil {
			continue
		}
		x = append(x, xp)
		y = append(y, yp)
	}
	return x, y, scanner.Err()
}

// process data in the filename
func process(p params, filename string) error {
	r, err := os.Open(filename)
	if err != nil {
		return err
	}

	fmt.Println("#", p.left, p.right, p.bottom, p.top)

	x, y, err := readata(r)
	if err != nil {
		return err
	}

	pminx := largest
	pmaxx := smallest
	fmt.Printf("polygon \"")
	for i := 0; i < len(x); i++ {
		px := x[i]
		if px > pmaxx {
			pmaxx = px
		}
		if px < pminx {
			pminx = px
		}
		fmt.Printf("%.3g ", px)
	}
	fmt.Printf("%.3g\"", x[0])

	pminy := largest
	pmaxy := smallest
	fmt.Printf("  \"")
	for i := 0; i < len(y); i++ {
		py := y[i]
		if py > pmaxy {
			pmaxy = py
		}
		if py < pminy {
			pminy = py
		}
		fmt.Printf("%.3g ", py)
	}
	fmt.Printf("%.3g\" \"%s\"\n", y[0], p.color)
	if len(p.label) > 0 {
		fmt.Printf("ctext \"%s\" %g %g 1\n", p.label, pminx+((pmaxx-pminx)/2), pminy+((pmaxy-pminy)/2))
	}
	return r.Close()
}

// vmap maps one range to another
func vmap(value, low1, high1, low2, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}
