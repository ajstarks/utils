// mkpoly - generate decksh polygons from x,y pairs
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
	left, right, bottom, top, minx, maxx, miny, maxy float64
	label, color                                     string
}

func main() {
	var p params
	flag.Float64Var(&p.left, "left", 10, "left")
	flag.Float64Var(&p.right, "right", 90, "right")
	flag.Float64Var(&p.bottom, "bottom", 10, "bottom")
	flag.Float64Var(&p.top, "top", 90, "top")
	flag.Float64Var(&p.minx, "minx", smallest, "top")
	flag.Float64Var(&p.maxx, "maxx", largest, "minx")
	flag.Float64Var(&p.miny, "miny", smallest, "maxx")
	flag.Float64Var(&p.maxy, "maxy", largest, "miny")
	flag.StringVar(&p.color, "color", "gray", "color")
	flag.StringVar(&p.label, "label", "", "label")
	flag.Parse()

	for _, f := range flag.Args() {
		err := process(p, f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
	}
}

// readata reads name, value pairs, and returns the min and max values
func readata(r io.Reader) ([]float64, []float64) {
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
	return x, y
}

// process data in the filename
func process(p params, filename string) error {
	r, err := os.Open(filename)
	if err != nil {
		return err
	}

	fmt.Println("#", p.minx, p.maxx, p.miny, p.maxy, p.left, p.right, p.bottom, p.top)
	pminx := largest
	pmaxx := smallest
	x, y := readata(r)

	fmt.Printf("polygon \"")
	for i := 0; i < len(x); i++ {
		px := vmap(x[i], p.minx, p.maxx, p.left, p.right)
		if px > pmaxx {
			pmaxx = px
		}
		if px < pminx {
			pminx = px
		}
		fmt.Printf("%.3g ", px)
	}
	fmt.Printf("%.3g\"", vmap(x[0], p.minx, p.maxx, p.left, p.right))

	pminy := largest
	pmaxy := smallest
	fmt.Printf("  \"")
	for i := 0; i < len(y); i++ {
		py := vmap(y[i], p.miny, p.maxy, p.bottom, p.top)
		if py > pmaxy {
			pmaxy = py
		}
		if py < pminy {
			pminy = py
		}
		fmt.Printf("%.3g ", py)
	}
	fmt.Printf("%.3g\" \"%s\"\n", vmap(y[0], p.miny, p.maxy, p.bottom, p.top), p.color)
	if len(p.label) > 0 {
		fmt.Printf("ctext \"%s\" %g %g 1\n", p.label, pminx+((pmaxx-pminx)/2), pminy+((pmaxy-pminy)/2))
	}
	return r.Close()
}

// vmap maps one range to another
func vmap(value, low1, high1, low2, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}
