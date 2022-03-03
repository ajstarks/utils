package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/ajstarks/deck/generate"
)

type place struct {
	name     string
	distance float64
}

type distanceTable struct {
	name string
	dist []place
}

func main() {
	data, err := readtable(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading the input")
		os.Exit(1)
	}
	makeslide(os.Stdout, data, "Morris County, New Jersey")
}

// makeside makes the slide deck
func makeslide(w io.Writer, data []distanceTable, title string) {
	deck := generate.NewSlides(w, 0, 0)
	deck.StartDeck()
	deck.StartSlide()
	deck.Text(40, 95, "TABLE OF DISTANCES", "sans", 3, "gray")
	deck.Text(40, 89, title, "sans", 3.5, "")
	deck.TextBlock(40, 85, "Showing the distance in miles and tenths, in an air line from one place to another", "serif", 1.5, 50, "")
	distable(deck, data, 1, 90, 1.1)
	deck.EndSlide()
	deck.EndDeck()
}

// readtable reads in distance table data
// name1
// <tab>place1:distance
// <tab>place2:distance
// ...
func readtable(r io.Reader) ([]distanceTable, error) {
	var table []distanceTable
	var t distanceTable
	var p place
	var places []place
	scanner := bufio.NewScanner(r)
	n := -1
	for scanner.Scan() {
		text := scanner.Text()
		// single name
		if !strings.Contains(text, "\t") {
			t.name = text
			table = append(table, t)
			places = make([]place, 0)
			n++
			continue
		}
		// <tab>name:distance
		if strings.Contains(text, "\t") {
			i := strings.Index(text, ":")
			if i > 0 && len(text) > 3 {
				d, _ := strconv.ParseFloat(text[i+1:], 64)
				p.name = text[1:i]
				p.distance = d
				places = append(places, p)
				table[n].dist = places
			}
		}
	}
	err := scanner.Err()
	return table, err
}

// dumptable prints out the distance table to an io.Writer
func dumptable(w io.Writer, table []distanceTable) {
	for _, t := range table {
		fmt.Fprintf(w, "%s\n", t.name)
		for _, d := range t.dist {
			fmt.Fprintf(w, "\t%s:%.2f\n", d.name, d.distance)
		}
	}
}

// distable makes a distance table using deck markup
func distable(deck *generate.Deck, table []distanceTable, left, top, size float64) {
	distleft := left + (size * 10)
	x := distleft
	y := top
	vspacing := size * 2.4
	hspacing := size * 2.4

	// vertical column headings
	for _, t := range table {
		deck.TextRotate(x, y-vspacing, t.name, "", "serif", 90, size, "")
		deck.Line(x-size-0.2, y-1, x-size-0.2, 2, 0.05, "gray")
		x += hspacing
		y -= vspacing
	}
	// horizontal headings, data
	x = left
	y = top - vspacing
	for _, t := range table {
		// place names
		deck.Text(x, y, t.name, "serif", size, "")
		dx := distleft
		dy := y
		// distances for each place
		for _, d := range t.dist {
			td := strconv.FormatFloat(d.distance, 'f', 1, 64)
			deck.TextMid(dx, dy, td, "mono", size, "")
			dx += hspacing
		}
		deck.Line(distleft-size, y-1, dx+size+0.3, y-1, 0.05, "gray")
		y -= vspacing
	}
}
