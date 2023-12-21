package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/ajstarks/kml"
)

// vmap maps one interval to another
func vmap(value float64, low1 float64, high1 float64, low2 float64, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

// readData reads lat/long pairs (separated by white space) from a file, mapping to deck coordinates
func readData(r io.Reader, g kml.Geometry) ([]float64, []float64, error) {
	x := []float64{}
	y := []float64{}
	s := bufio.NewScanner(r)
	for s.Scan() {
		t := s.Text()
		f := strings.Fields(t)
		if len(f) != 2 {
			continue
		}
		xp, err := strconv.ParseFloat(f[1], 64) // latitude
		if err != nil {
			continue
		}
		yp, err := strconv.ParseFloat(f[0], 64) // longitude
		if err != nil {
			continue
		}
		x = append(x, vmap(xp, g.Longmin, g.Longmax, g.Xmin, g.Xmax))
		y = append(y, vmap(yp, g.Latmin, g.Latmax, g.Ymin, g.Ymax))
	}
	return x, y, s.Err()
}

// readStats reads lat/long pairs, report on the computed bounding box and center
func readStats(r io.Reader) {
	maxxval := -100000000.0
	minxval := 100000000.0
	maxyval := -100000000.0
	minyval := 100000000.0

	s := bufio.NewScanner(r)
	for s.Scan() {
		t := s.Text()
		f := strings.Fields(t)
		if len(f) != 2 {
			continue
		}
		xp, err := strconv.ParseFloat(f[1], 64) // latitude
		if err != nil {
			continue
		}
		yp, err := strconv.ParseFloat(f[0], 64) // longitude
		if err != nil {
			continue
		}
		if xp > maxxval {
			maxxval = xp
		}
		if xp < minxval {
			minxval = xp
		}

		if yp > maxyval {
			maxyval = yp
		}
		if yp < minyval {
			minyval = yp
		}
	}
	centerLong := minxval + (maxxval-minxval)/2
	centerLat := minyval + (maxyval-minyval)/2
	fmt.Fprintf(os.Stdout, "--center=%v,%v -bbox=\"%v,%v|%v,%v\" --longmin=%v --longmax=%v --latmin=%v --latmax=%v\n",
		centerLat, centerLong, maxyval, minxval, minyval, maxxval, minxval, maxxval, minyval, maxyval)
}

// process processing input and options, making markup
func process(filename string, info bool, shape, style, color, bbox string, linewidth float64, mapgeo kml.Geometry) {

	// read from stdin by default, open a file if specified
	r := os.Stdin
	if len(filename) > 0 {
		var rerr error
		r, rerr = os.Open(filename)
		if rerr != nil {
			fmt.Fprintf(os.Stderr, "%v\n", rerr)
			return
		}
	}

	// just show info, then return, if specified
	if info {
		readStats(r)
		return
	}

	// read coordinates
	x, y, err := readData(r, mapgeo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	// make a bounding box, if specified
	if len(bbox) > 0 {
		kml.BoundingBox(mapgeo, bbox, style)
	}
	// make the drawing
	kml.Deckshape(shape, style, x, y, linewidth, color, mapgeo)
	r.Close()
}

func main() {
	var mapgeo kml.Geometry
	var fulldeck, info bool
	var linewidth float64
	var color, bbox, shape, bgcolor, style string

	// options
	flag.Float64Var(&mapgeo.Xmin, "xmin", 5, "canvas x minimum")
	flag.Float64Var(&mapgeo.Xmax, "xmax", 95, "canvas x maxmum")
	flag.Float64Var(&mapgeo.Ymin, "ymin", 5, "canvas y minimum")
	flag.Float64Var(&mapgeo.Ymax, "ymax", 95, "canvas y maximum")
	flag.Float64Var(&mapgeo.Latmin, "latmin", -90, "latitude x minimum")
	flag.Float64Var(&mapgeo.Latmax, "latmax", 90, "latitude x maxmum")
	flag.Float64Var(&mapgeo.Longmin, "longmin", -180, "longitude y minimum")
	flag.Float64Var(&mapgeo.Longmax, "longmax", 180, "longitude y maximum")
	flag.Float64Var(&linewidth, "linewidth", 0.1, "line width")
	flag.StringVar(&color, "color", "black", "line color")
	flag.StringVar(&bbox, "bbox", "", "bounding box color (\"\" no box)")
	flag.StringVar(&shape, "shape", "polyline", "polygon, polyline")
	flag.StringVar(&style, "style", "deck", "deck, decksh, plain")
	flag.StringVar(&bgcolor, "bgcolor", "", "background color")
	flag.BoolVar(&fulldeck, "fulldeck", true, "make a full deck")
	flag.BoolVar(&info, "info", false, "only report center and bounding box info")
	flag.Parse()

	// don't do any generation if info only
	if info {
		fulldeck = false
	}
	// add deck/slide markup, if specified
	if fulldeck {
		kml.Deckshbegin(bgcolor)
	}
	// for every file (or stdin if no files are specified), make markup
	if len(flag.Args()) == 0 {
		process("", info, shape, style, color, bbox, linewidth, mapgeo)
	} else {
		for _, filename := range flag.Args() {
			process(filename, info, shape, style, color, bbox, linewidth, mapgeo)
		}
	}
	// end the deck, if specified
	if fulldeck {
		kml.Deckshend()
	}
}
