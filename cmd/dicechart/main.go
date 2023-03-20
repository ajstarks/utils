// dicechart: make a Negro Year Book style dice chart using deck markup
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

type dicedata struct {
	name  string
	value int
}

type config struct {
	diceunit    int
	cw          float64
	ch          float64
	top         float64
	vskip       float64
	textsize    float64
	valuesize   float64
	labelx      float64
	datax       float64
	dicewidth   float64
	dicespacing float64
	dotsize     float64
	dotcolor    string
	title       string
}

const (
	cw          = 792.0
	ch          = 612.0
	top         = 80.0
	vskip       = 7.0
	textsize    = 2.0
	valuesize   = 2.0
	labelx      = 10.0
	datax       = 35.0
	dicewidth   = 1.5
	dotsize     = 0.75
	dicespacing = 5.0
	dotcolor    = "black"
	diceunit    = 5
	legendy     = 5.0
)

// xmlmap defines the XML substitutions
var xmlmap = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;")

// xmlesc XML escapes a string
func xmlesc(s string) string {
	return xmlmap.Replace(s)
}

// polar to Cartesian coordinates, corrected for aspect ratio
func polar(cx, cy, r, theta, cw, ch float64) (float64, float64) {
	ry := r * (cw / ch)
	t := theta * (math.Pi / 180)
	return cx + (r * math.Cos(t)), cy + (ry * math.Sin(t))
}

// beginDeck starts a deck
func beginDeck(w io.Writer) {
	fmt.Fprintln(w, "<deck><slide>")
}

// endDeck ends a deck
func endDeck(w io.Writer) {
	fmt.Fprintln(w, "</slide></deck>")
}

// text renders text at specified location and size
func text(w io.Writer, s string, x, y, size float64) {
	fmt.Fprintf(w, "<text xp=\"%v\" yp=\"%v\" sp=\"%v\">%s</text>\n", x, y, size, xmlesc(s))
}

// ctext makes centered text
func ctext(w io.Writer, s string, x, y, size float64) {
	fmt.Fprintf(w, "<text align=\"c\" xp=\"%v\" yp=\"%v\"sp=\"%v\">%s</text>\n", x, y, size, xmlesc(s))
}

// circle makes a filled circle
func circle(w io.Writer, x, y, r float64, color string) {
	fmt.Fprintf(w, "<ellipse xp=\"%v\" yp=\"%v\" wp=\"%v\" hr=\"100\" color=%q/>\n", x, y, r, color)
}

// readData reads in name,value pairs in CSV format
func readData(r io.Reader) []dicedata {
	var item dicedata
	var datum []dicedata
	input := csv.NewReader(r)
	for {
		record, err := input.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		if len(record) != 2 {
			continue
		}
		item.name = record[0]
		item.value, _ = strconv.Atoi(record[1])
		datum = append(datum, item)
	}
	return datum
}

// dicerow makes a labeled row of dice
func dicerow(w io.Writer, d dicedata, cfg config, y float64) {
	ly := y - (cfg.textsize / 3)
	text(w, d.name, cfg.labelx, ly, cfg.textsize)
	xp := cfg.datax
	for i := 0; i < d.value/cfg.diceunit; i++ {
		fivedots(w, xp, y, cfg.dicewidth, cfg.dotsize, cfg.dotcolor)
		xp += cfg.dicespacing
	}
	rem := d.value % cfg.diceunit
	dice(w, xp, y, cfg.dicewidth, cfg.dotsize, rem, "red")
	legend(w, cfg)

	// nudge the value optimally next to the last block
	var ns float64
	if cfg.diceunit > 5 {
		rem /= 5
	}
	switch rem {
	case 0:
		ns = cfg.dicespacing * -0.4
	case 1, 2:
		ns = cfg.dicespacing * 0.1
	case 3, 4:
		ns = cfg.dicespacing / 2
	}
	text(w, strconv.Itoa(d.value), xp+ns, ly, cfg.valuesize)
}

// dicechart reads data and makes the chart
func dicechart(w io.Writer, r io.Reader, cfg config) {
	data := readData(r)
	beginDeck(w)
	if len(cfg.title) > 0 {
		ctext(w, cfg.title, 50, cfg.top+(cfg.textsize*4), cfg.textsize*1.5)
	}
	y := cfg.top
	for _, d := range data {
		dicerow(w, d, cfg, y)
		y -= cfg.vskip
	}
	endDeck(w)
}

// dice makes a one, two, three, four, or five dot die.
func dice(w io.Writer, x, y, r, size float64, n int, color string) {
	x1, y1 := polar(x, y, r, 135, cw, ch) // top left
	x2, y2 := polar(x, y, r, 225, cw, ch) // bottom left
	x3, y3 := polar(x, y, r, 45, cw, ch)  // top right
	x4, y4 := polar(x, y, r, 315, cw, ch) // bottom right
	nd := n
	if n > 5 {
		nd = n / 5
	}
	switch nd {
	case 1:
		circle(w, x1, y1, size, color)
	case 2:
		circle(w, x1, y1, size, color)
		circle(w, x2, y2, size, color)
	case 3:
		circle(w, x1, y1, size, color)
		circle(w, x2, y2, size, color)
		circle(w, x3, y3, size, color)
	case 4:
		circle(w, x1, y1, size, color)
		circle(w, x2, y2, size, color)
		circle(w, x3, y3, size, color)
		circle(w, x4, y4, size, color)
	case 5:
		circle(w, x1, y1, size, color)
		circle(w, x2, y2, size, color)
		circle(w, x3, y3, size, color)
		circle(w, x4, y4, size, color)
		circle(w, x, y, size, color)
	}
}

// legend makes dice / unit legend
func legend(w io.Writer, cfg config) {
	ly := legendy - cfg.dotsize
	fivedots(w, cfg.datax, legendy, cfg.dicewidth/2, cfg.dotsize/2, cfg.dotcolor)
	text(w, strconv.Itoa(cfg.diceunit)+" items", cfg.datax+cfg.dicewidth, ly, cfg.textsize*0.7)
}

// fivedots makes a full 5-dot die
func fivedots(w io.Writer, x, y, r, size float64, color string) {
	x1, y1 := polar(x, y, r, 135, cw, ch) // top left
	x2, y2 := polar(x, y, r, 225, cw, ch) // bottom left
	x3, y3 := polar(x, y, r, 45, cw, ch)  // top right
	x4, y4 := polar(x, y, r, 315, cw, ch) // bottom right
	circle(w, x1, y1, size, color)
	circle(w, x2, y2, size, color)
	circle(w, x3, y3, size, color)
	circle(w, x4, y4, size, color)
	circle(w, x, y, size, color)
}

// setup processes command line flags and sets where data is read from
func setup() (config, io.Reader, error) {
	var cfg config
	flag.IntVar(&cfg.diceunit, "unit", diceunit, "dice unit")
	flag.Float64Var(&cfg.cw, "width", cw, "canvas width")
	flag.Float64Var(&cfg.ch, "height", ch, "canvas height")
	flag.Float64Var(&cfg.top, "top", top, "top of the chart")
	flag.Float64Var(&cfg.vskip, "vskip", vskip, "vertical skip")
	flag.Float64Var(&cfg.textsize, "textsize", textsize, "canvas width")
	flag.Float64Var(&cfg.valuesize, "valsize", valuesize, "canvas width")
	flag.Float64Var(&cfg.labelx, "lx", labelx, "label left position")
	flag.Float64Var(&cfg.datax, "dx", labelx+25, "data left position")
	flag.Float64Var(&cfg.dicewidth, "dw", dicewidth, "dice width")
	flag.Float64Var(&cfg.dicespacing, "ds", dicespacing, "dice spacing")
	flag.Float64Var(&cfg.dotsize, "dotsize", dotsize, "dot size")
	flag.StringVar(&cfg.dotcolor, "color", dotcolor, "dotcolor")
	flag.StringVar(&cfg.title, "title", "", "chart title")
	flag.Parse()

	var err error
	r := os.Stdin
	if len(flag.Args()) > 0 {
		r, err = os.Open(flag.Arg(0))
	}
	return cfg, r, err
}

func main() {
	cfg, r, err := setup()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	dicechart(os.Stdout, r, cfg)
}
