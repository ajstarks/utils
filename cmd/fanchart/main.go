// fanchart -- make a fanchart like Dubois plate 27 from CSV data,
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
	midx        = 50.0
	ty          = 95.0
	arcsize     = 30.0
	wrapwidth   = 12.0
	titlesize   = 3.0
	notesize    = titlesize * 0.6
	topbegAngle = 145.0
	botbegAngle = 215.0
	fanspan     = 110.0
)

// title makes a title
func title(s string) {
	ctext(s, midx, ty, titlesize)
}

// arc draws a filled arc
func arc(cx, cy, a1, a2, size float64, color string) {
	fmt.Printf("arc %.2f %.2f %.2f %.2f %.2f %.2f %.2f %q\n", cx, cy, size, size, a1, a2, size, color)
}

// circle makes a filled circle
func circle(x, y, r float64, color string) {
	fmt.Printf("circle %v %v %v %q\n", x, y, r, color)
}

// textblock makes a block of text
func textblock(s string, x, y, size, width float64) {
	fmt.Printf("textblock \"%s\" %v %v %v %v\n", s, x, y, width, size)
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
	tw := wrapwidth
	// left legend
	for i := 0; i < left; i++ {
		label := data[i].name
		circle(x, y, r, data[i].color)
		textblock(label, x+3, textshift(label, y, ts), ts, tw)
		y -= 10.0
	}
	// right legend
	x = 100 - x
	y = 60
	for i := left; i < len(data); i++ {
		label := data[i].name
		circle(x, y, r, data[i].color)
		textblock(label, x-20, textshift(label, y, ts), ts, tw)
		y -= 10.0
	}
}

// textshift aligns text to a point depending on the size of the text
func textshift(s string, y, ts float64) float64 {
	var ty float64
	w := strings.Split(s, " ")
	if len(w) > 2 {
		ty = y + (ts / 2)
	} else {
		ty = y - (ts / 3)
	}
	return ty
}

// arclabel
func arclabel(cx, cy, a1, a2, asize, value, cw, ch float64) {
	v := strconv.FormatFloat(value, 'f', 1, 64)
	diff := a2 - a1
	lx, ly := polar(cx, cy, asize*0.75, a1+(diff*0.5), cw, ch)
	fmt.Printf("ctext \"%s%%\" %.3f %.3f 1.5\n", v, lx, ly)
}

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
		if newset(record, n) {
			setnum++
		}
		switch setnum {
		case 1:
			if len(record[1]) == 0 && len(record[2]) == 0 { // set header
				topdata.name = record[0]
			} else { // load set data
				td.name = record[0]
				td.value, _ = strconv.ParseFloat(record[1], 64)
				td.color = record[2]
				tds = append(tds, td)
				topcount++
			}
		case 2:
			if len(record[1]) == 0 && len(record[2]) == 0 { // set header
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

func newset(s []string, n int) bool {
	return len(s[1]) == 0 && len(s[2]) == 0
}

func beginDeck(w, h float64) {
	fmt.Printf("deck\ncanvas %v %v\n", w, h)
}

func beginSlide() {
	fmt.Println("slide")
}

func endSlide() {
	fmt.Println("eslide")
}

func endDeck() {
	fmt.Println("edeck")
}

func note(s string) {
	ctext(s, 50, 10, notesize)
}
func main() {
	var canvasWidth, canvasHeight, arcsize float64
	flag.Float64Var(&canvasHeight, "h", 612, "canvas height")
	flag.Float64Var(&canvasWidth, "w", 792, "canvas width")
	flag.Float64Var(&arcsize, "size", 35, "fansize")
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
