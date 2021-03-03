// fanchart -- make a fanchart like Dubois plate 27, reading from a CSV data,
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

// Measure describes the data set
type Measure struct {
	name  string
	value float64
	color string
}

// Dataset is a labeled data set
type Dataset struct {
	name     string
	measures []Measure
}

const (
	midx        = 50.0            // middle of the canvas
	ty          = 95.0            // title y coordinate
	arcsize     = 30.0            // size of the wedges
	wrapwidth   = 12.0            // when to wrap legend titles
	titlesize   = 3.0             // text size of titles
	notesize    = titlesize * 0.6 // text size of footnotes
	topbegAngle = 145.0           // top beginning angle
	botbegAngle = 215.0           // bottom beginning angle
	fanspan     = 110.0           // span size of the top and bottom of the fan
)

// title makes a title
func title(s string) {
	ctext(s, midx, ty, titlesize)
}

// arc draws a filled arc
func arc(cx, cy, a1, a2, size float64, color string) {
	fmt.Printf("arc %.2f %.2f %.2f %.2f %.2f %.2f %.2f %q\n",
		cx, cy, size, size, a1, a2, size, color)
}

// circle makes a filled circle
func circle(x, y, r float64, color string) {
	fmt.Printf("circle %v %v %v %q\n", x, y, r, color)
}

// textblock makes a block of text
func textblock(s string, x, y, size, width float64) {
	fmt.Printf("textblock \"%s\" %v %v %v %v\n", s, x, y, width, size)
}

// text renders text at specified location and size
func text(s string, x, y, size float64) {
	fmt.Printf("text \"%s\" %v %v %v\n", s, x, y, size)
}

// ctext makes centered text
func ctext(s string, x, y, size float64) {
	fmt.Printf("ctext \"%s\" %v %v %v\n", s, x, y, size)
}

// legend makes a balanced left and right hand legend
func legend(data []Measure, rows int, ts float64) {
	right := len(data) % rows
	left := len(data) - right
	x := 5.0
	y := 60.0
	r := ts + 1.0
	//tw := wrapwidth
	// left legend
	for i := 0; i < left; i++ {
		label := data[i].name
		circle(x, y, r, data[i].color)
		legendlabel(label, x+3, y, ts)
		y -= 10.0
	}
	// right legend
	x = 100 - x
	y = 60
	for i := left; i < len(data); i++ {
		label := data[i].name
		circle(x, y, r, data[i].color)
		legendlabel(label, x-20, y, ts)
		y -= 10.0
	}
}

// legendlabel lays out the legend labels
func legendlabel(s string, x, y, ts float64) {
	w := strings.Split(s, `\n`)
	lw := len(w)
	if lw == 1 {
		text(s, x, y-(ts/3), ts)
	} else {
		y = y + (ts * (float64(lw / 3)))
		for i := 0; i < lw; i++ {
			text(w[i], x, y, ts)
			y -= (ts * 1.5)
		}
	}
}

// arclabel
func arclabel(cx, cy, a1, a2, asize, value, cw, ch float64) {
	v := strconv.FormatFloat(value, 'f', 1, 64)
	diff := a2 - a1
	lx, ly := polar(cx, cy, asize*0.75, a1+(diff*0.5), cw, ch)
	fmt.Printf("ctext \"%s%%\" %.3f %.3f 1.5\n", v, lx, ly)
}

// polar to Cartesian coordinates, corrected for aspect ratio
func polar(cx, cy, r, theta, cw, ch float64) (float64, float64) {
	ry := r * (cw / ch)
	t := theta * (math.Pi / 180)
	return cx + (r * math.Cos(t)), cy + (ry * math.Sin(t))
}

// topfan makes the top of the fan chart
func topfan(dataset Dataset, cx, cy, asize, start, fansize, cw, ch float64) {
	lx, ly := polar(cx, cy, asize+1, 90, cw, ch)
	ctext(dataset.name, lx, ly, 2)
	data := dataset.measures
	for _, d := range data {
		m := (d.value / 100) * fansize
		a1 := start - m
		a2 := start
		arc(cx, cy, a1, a2, asize, d.color)
		arclabel(cx, cy, a1, a2, asize, d.value, cw, ch)
		start = a1
	}
}

// botfan makes the bottom of the fan chart
func botfan(dataset Dataset, cx, cy, asize, start, fansize, cw, ch float64) {
	lx, ly := polar(cx, cy, asize+2, 270, cw, ch)
	ctext(dataset.name, lx, ly, 2)
	data := dataset.measures
	for i := len(data) - 1; i >= 0; i-- {
		d := data[i]
		m := (d.value / 100) * fansize
		a1 := start + m
		a2 := start
		arc(cx, cy, a2, a1, asize, d.color)
		arclabel(cx, cy, a1, a2, asize, d.value, cw, ch)
		start = a1
	}
}

// readData reads a CSV file containing top and bottom fan data
// File layout:
//
// column headers
// title,footnotes
// section name
// item,value,color
// ...
// bottom section name
// item,value,color
// ...
func readData(filename string) (Dataset, Dataset, error) {
	var topdata, botdata Dataset
	var td, bd Measure
	var tds, bds []Measure
	r, err := os.Open(filename)
	if err != nil {
		return topdata, botdata, err
	}
	input := csv.NewReader(r)
	n := 0
	topcount := 0
	botcount := 0
	setnum := 0
	for {
		record, err := input.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(record) != 3 {
			return topdata, botdata, err
		}
		n++
		// skip header
		if n == 1 {
			continue
		}
		// title is next
		if n == 2 {
			title(record[0])
			if len(record[1]) > 0 {
				note(record[1])
			}
			continue
		}
		// check to see if we are in the top (setnum=1) or bottom set (setnum=2)
		if isheader(record) {
			setnum++
		}
		switch setnum {
		case 1:
			if isheader(record) { // set header
				topdata.name = record[0]
			} else { // load set data
				td.name = record[0]
				td.value, _ = strconv.ParseFloat(record[1], 64)
				td.color = record[2]
				tds = append(tds, td)
				topcount++
			}
		case 2:
			if isheader(record) { // set header
				botdata.name = record[0]
			} else { //load  set data
				bd.name = record[0]
				bd.value, _ = strconv.ParseFloat(record[1], 64)
				bd.color = record[2]
				bds = append(bds, bd)
				botcount++
			}
		}
	}
	if topcount != botcount {
		fmt.Fprintf(os.Stderr,
			"The number of top items, %d is not the same as the bottom: %d\n", topcount, botcount)
	}
	topdata.measures = tds
	botdata.measures = bds
	return topdata, botdata, nil
}

// newset determines if a new set of data has begun in the input
func isheader(s []string) bool {
	return len(s[1]) == 0 && len(s[2]) == 0
}

// beginDeck makes the markup to begin a deck
func beginDeck(w, h float64) {
	fmt.Printf("deck\ncanvas %v %v\n", w, h)
}

// beginSlide makes the markup to begin a slide
func beginSlide() {
	fmt.Println("slide")
}

// endSlide makes the markup to end a slide
func endSlide() {
	fmt.Println("eslide")
}

// beginDeck makes the markup to begin a deck
func endDeck() {
	fmt.Println("edeck")
}

// note makes a footnote
func note(s string) {
	ctext(s, 50, 3, notesize)
}

func main() {
	var canvasWidth, canvasHeight, arcsize float64
	flag.Float64Var(&canvasHeight, "h", 612, "canvas height") // canvas height
	flag.Float64Var(&canvasWidth, "w", 792, "canvas width")   // canvas width
	flag.Float64Var(&arcsize, "size", 30, "fansize")          // size of the fan
	flag.Parse()

	beginDeck(canvasWidth, canvasHeight)
	for _, f := range flag.Args() {
		beginSlide()
		topdata, botdata, err := readData(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		topfan(topdata, midx, 50, arcsize, topbegAngle, fanspan, canvasWidth, canvasHeight)
		botfan(botdata, midx, 50, arcsize, botbegAngle, fanspan, canvasWidth, canvasHeight)
		legend(topdata.measures, 3, 1.5)
		endSlide()
	}
	endDeck()
}
