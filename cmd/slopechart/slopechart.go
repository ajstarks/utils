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

	"github.com/ajstarks/deck/generate"
)

type nameval struct {
	name  string
	value float64
	label string
}

type options struct {
	min, max, left, right, bottom, top, textsize, linewidth float64
	color, vcolor                                           string
}

var xmlmap = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;")

const (
	largest  = math.MaxFloat64
	smallest = -math.MaxFloat64
)

// xmlesc escapes XML
func xmlesc(s string) string {
	return xmlmap.Replace(s)
}

func readData(r io.ReadCloser) ([]nameval, string, float64, float64, error) {
	var d nameval
	var data []nameval
	var err error
	maxval := smallest
	minval := largest
	title := ""
	scanner := bufio.NewScanner(r)
	// read a line, parse into name, value pairs
	// compute min and max values
	for scanner.Scan() {
		t := scanner.Text()
		if len(t) == 0 { // skip blank lines
			continue
		}
		if len(t) > 2 && t[0] == '#' {
			title = strings.TrimSpace(t[1:])
		}
		fields := strings.Split(t, "\t")
		if len(fields) < 2 {
			continue
		}
		if len(fields) > 2 {
			d.label = fields[2]
		}
		d.name = fields[0]
		d.value, err = strconv.ParseFloat(fields[1], 64)
		if err != nil {
			d.value = 0
		}
		if d.value > maxval {
			maxval = d.value
		}
		if d.value < minval {
			minval = d.value
		}
		data = append(data, d)
	}
	r.Close()
	return data, title, minval, maxval, err

}

// vmap maps one range into another
func vmap(value float64, low1 float64, high1 float64, low2 float64, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

func slopechart(deck *generate.Deck, opts options, r io.ReadCloser) error {
	data, title, _, _, err := readData(r)
	if err != nil {
		return err
	}
	if len(data) < 2 {
		return fmt.Errorf("need at least two data points")
	}
	min := opts.min
	max := opts.max
	left := opts.left
	right := opts.right
	top := opts.top
	bottom := opts.bottom
	color := opts.color
	vcolor := opts.vcolor
	textsize := opts.textsize
	linewidth := opts.linewidth
	lw := linewidth / 2

	lsize := textsize * 0.75
	tsize := textsize * 1.5
	w := right - left
	h := top - bottom
	if len(title) > 0 {
		hsize := textsize * 2
		deck.Text(left, top+10, title, "sans", hsize, "")
	}

	hskip := w * .60
	vskip := h * 1.4
	x1 := left
	x2 := right
	for i := 0; i < len(data)-1; i += 2 {
		if len(data[i].label) > 0 {
			deck.TextMid(x1+(w/2), top+3, data[i].label, "sans", tsize, "")
		}
		v1 := data[i].value
		v2 := data[i+1].value
		v1y := vmap(v1, min, max, bottom, top)
		v2y := vmap(v2, min, max, bottom, top)
		deck.Line(x1, bottom, x1, top, lw, "black")
		deck.Line(x2, bottom, x2, top, lw, "black")
		deck.Circle(x1, v1y, textsize, color)
		deck.Circle(x2, v2y, textsize, color)
		deck.Line(x1, v1y, x2, v2y, linewidth, color)
		deck.TextMid(x1, bottom-2, data[i].name, "sans", textsize, "")
		deck.TextMid(x2, bottom-2, data[i+1].name, "sans", textsize, "")
		deck.TextEnd(x1-1, top, fmt.Sprintf("%g", max), "sans", lsize, "")
		deck.TextEnd(x1-1, v1y, fmt.Sprintf("%g", v1), "sans", lsize, vcolor)
		deck.Text(x2+1, v2y, fmt.Sprintf("%g", v2), "sans", lsize, vcolor)
		x1 += w + hskip
		x2 += w + hskip
		if x2 > 100 {
			x1 = left
			x2 = right
			top -= vskip
			bottom -= vskip
		}
	}
	return err
}

func main() {

	left := flag.Float64("left", 20, "left")
	right := flag.Float64("right", 40, "right")
	bottom := flag.Float64("bottom", 20, "bottom")
	top := flag.Float64("top", 60, "top")
	color := flag.String("color", "steelblue", "color")
	vcolor := flag.String("vcolor", "maroon", "value color")
	max := flag.Float64("max", 100, "max value")
	textsize := flag.Float64("textsize", 1.5, "text size")
	linewidth := flag.Float64("linewidth", 0.2, "line width")
	flag.Parse()

	opts := options{
		min:       0,
		max:       *max,
		left:      *left,
		right:     *right,
		top:       *top,
		bottom:    *bottom,
		linewidth: *linewidth,
		textsize:  *textsize,
		color:     *color,
		vcolor:    *vcolor,
	}

	deck := generate.NewSlides(os.Stdout, 0, 0)
	deck.StartDeck()
	deck.StartSlide()
	if err := slopechart(deck, opts, os.Stdin); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	deck.EndSlide()
	deck.EndDeck()
	os.Exit(0)
}
