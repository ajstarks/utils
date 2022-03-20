// fanchart -- make a fanchart like Dubois plate 27, reading from a CSV data
// generates deck markup
// usage: fanchart file | deckrenderer
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
	midx          = 50.0            // middle of the canvas
	midy          = 50.0            // middle of the canvas
	ty            = 95.0            // title y coordinate
	arcsize       = 30.0            // size of the wedges
	titlesize     = 3.0             // text size of titles
	catsize       = 2.0             // size of data category label
	legendsize    = 1.8             // legend label size
	labelsize     = 1.5             // data label size
	notesize      = titlesize * 0.6 // text size of footnotes
	topbegAngle   = 145.0           // top beginning angle
	botbegAngle   = 215.0           // bottom beginning angle
	fanspan       = 110.0           // span size of the top and bottom of the fan
	leftbegAngle  = 135.0           // left beginning angle
	rightbegAngle = 315.0           // right beginning angle
	wingspan      = 90.0            // span size of the left and right wings
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

// title makes a title
func title(s string) {
	ctext(s, midx, ty, titlesize)
}

// arc draws a filled arc
func arc(cx, cy, a1, a2, size float64, color string) {
	fmt.Printf(
		"<arc xp=\"%.2f\" yp=\"%.2f\" wp=\"%.2f\" hp=\"%.2f\" a1=\"%.2f\" a2=\"%.2f\" sp=\"%.2f\" color=%q/>\n",
		cx, cy, size, size, a1, a2, size, color)
}

// circle makes a filled circle
func circle(x, y, r float64, color string) {
	fmt.Printf("<ellipse xp=\"%v\" yp=\"%v\" wp=\"%v\" hr=\"100\" color=%q/>\n", x, y, r, color)
}

// text renders text at specified location and size
func text(s string, x, y, size float64) {
	fmt.Printf("<text xp=\"%v\" yp=\"%v\" sp=\"%v\">%s</text>\n", x, y, size, xmlesc(s))
}

// etext renders text at specified location and size, end justified
func etext(s string, x, y, size float64) {
	fmt.Printf("<text align=\"e\" xp=\"%v\" yp=\"%v\" sp=\"%v\">%s</text>\n", x, y, size, xmlesc(s))
}

// ctext makes centered text
func ctext(s string, x, y, size float64) {
	fmt.Printf("<text align=\"c\" xp=\"%v\" yp=\"%v\"sp=\"%v\">%s</text>\n", x, y, size, xmlesc(s))
}

// beginDeck makes the markup to begin a deck
func beginDeck(w, h float64) {
	fmt.Printf("<deck>\n<canvas width=\"%v\" height=\"%v\"/>\n", w, h)
}

// beginSlide makes the markup to begin a slide
func beginSlide(bgcolor, textcolor string) {
	fmt.Printf("<slide bg=%q fg=%q>\n", bgcolor, textcolor)
}

// endSlide makes the markup to end a slide
func endSlide() {
	fmt.Println("</slide>")
}

// beginDeck makes the markup to begin a deck
func endDeck() {
	fmt.Println("</deck>")
}

// legend makes a balanced left and right hand legend
func legend(data []Measure, orientation string, ts float64) {
	var x, y, xoffset float64
	l := len(data)
	h := l / 2
	rem := l % 2
	hr := h + rem

	r := ts + 1.0
	leading := ts * 6

	switch orientation {
	case "tb":
		x = 5.0
		y = 60.0
	case "lr":
		x = midx - 10
		y = 87.0
	}
	// left/top legend
	xoffset = 3
	for i := 0; i < hr; i++ {
		label := data[i].name
		circle(x, y, r, data[i].color)
		legendlabel(label, x+xoffset, y, ts)
		y -= leading
	}
	// right/bottom legend
	switch orientation {
	case "tb":
		x = 100 - x
		y = 60
		xoffset = -20.0
	case "lr":
		y = 25.0
	}
	for i := hr; i < len(data); i++ {
		label := data[i].name
		circle(x, y, r, data[i].color)
		legendlabel(label, x+xoffset, y, ts)
		y -= leading
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
			y -= (ts * 1.8)
		}
	}
}

// arclabel labels the data items
func arclabel(cx, cy, a1, a2, asize, value, cw, ch float64) {
	v := strconv.FormatFloat(value, 'f', 1, 64)
	diff := a2 - a1
	lx, ly := polar(cx, cy, asize*0.9, a1+(diff*0.5), cw, ch)
	ctext(v+"%", lx, ly, labelsize)
}

// polar to Cartesian coordinates, corrected for aspect ratio
func polar(cx, cy, r, theta, cw, ch float64) (float64, float64) {
	ry := r * (cw / ch)
	t := theta * (math.Pi / 180)
	return cx + (r * math.Cos(t)), cy + (ry * math.Sin(t))
}

// wedge makes data wedges
func wedge(data Dataset, cx, cy, begAngle, asize, cw, ch float64) {
	start := begAngle
	for _, d := range data.measures {
		m := (d.value / 100) * wingspan
		a1 := start
		a2 := start + m
		arc(cx, cy, a1, a2, asize, d.color)
		arclabel(cx, cy, a1, a2, asize, d.value, cw, ch)
		start = a2
	}
}

// wings makes left and right data "wings"
func wings(top, bot Dataset, cx, cy, asize, cw, ch float64) {
	var lx, ly float64
	lx, ly = polar(cx, cy, asize+1, 180, cw, ch)
	etext(top.name, lx, ly, legendsize)
	wedge(top, cx, cy, leftbegAngle, asize, cw, ch)
	lx, ly = polar(cx, cy, asize+1, 0, cw, ch)
	text(bot.name, lx, ly, legendsize)
	wedge(bot, cx, cy, rightbegAngle, asize, cw, ch)
}

// fan makes the top and bottom fan
func fan(top, bot Dataset, cx, cy, asize, cw, ch float64) {
	var lx, ly, start float64
	// the top of the fan chart
	lx, ly = polar(cx, cy, asize+1, 90, cw, ch)
	ctext(top.name, lx, ly, catsize)
	start = topbegAngle
	for _, d := range top.measures {
		m := (d.value / 100) * fanspan
		a1 := start - m
		a2 := start
		arc(cx, cy, a1, a2, asize, d.color)
		arclabel(cx, cy, a1, a2, asize, d.value, cw, ch)
		start = a1
	}
	// bottom of the fan chart
	lx, ly = polar(cx, cy, asize+2, 270, cw, ch)
	ctext(bot.name, lx, ly, catsize)
	start = botbegAngle
	for i := len(bot.measures) - 1; i >= 0; i-- {
		d := bot.measures[i]
		m := (d.value / 100) * fanspan
		a1 := start + m
		a2 := start
		arc(cx, cy, a2, a1, asize, d.color)
		arclabel(cx, cy, a1, a2, asize, d.value, cw, ch)
		start = a1
	}
}

// readData reads a CSV file containing top and bottom fan data
// File layout:
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
			"The number of top items, %d is not the same as the bottom: %d\n",
			topcount, botcount)
	}
	topdata.measures = tds
	botdata.measures = bds
	return topdata, botdata, nil
}

// newset determines if a new set of data has begun in the input
func isheader(s []string) bool {
	return len(s[1]) == 0 && len(s[2]) == 0
}

// note makes a footnote
func note(s string) {
	ctext(s, 50, 3, notesize)
}

func main() {
	var canvasWidth, canvasHeight, arcsize float64
	var orientation, textcolor, bgcolor string

	flag.Float64Var(&canvasHeight, "h", 612, "canvas height") // canvas height
	flag.Float64Var(&canvasWidth, "w", 792, "canvas width")   // canvas width
	flag.Float64Var(&arcsize, "size", 30, "fan/wing size")    // size of the fan
	flag.StringVar(&orientation, "dir", "tb", "orientation (tb=Top/Bottom, lr=Left/Right)")
	flag.StringVar(&bgcolor, "bgcolor", "white", "background color")
	flag.StringVar(&textcolor, "textcolor", "black", "text color")

	flag.Parse()

	beginDeck(canvasWidth, canvasHeight)
	for _, f := range flag.Args() {
		beginSlide(bgcolor, textcolor)
		data1, data2, err := readData(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		if orientation == "tb" {
			fan(data1, data2, midx, midy, arcsize, canvasWidth, canvasHeight)
		} else {
			wings(data1, data2, midx, midy, arcsize, canvasWidth, canvasHeight)
		}
		legend(data1.measures, orientation, labelsize)
		endSlide()
	}
	endDeck()
}
