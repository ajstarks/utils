// desordres -- tile blocks of lines as in Vera Molnar's Des Ordres, using deck markup
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

var palette = spalette{
	"kirokaze-gameboy":       {"#332c50", "#46878f", "#94e344", "#e2f3e4"},
	"ice-cream-gb":           {"#7c3f58", "#eb6b6f", "#f9a875", "#fff6d3"},
	"2-bit-demichrome":       {"#211e20", "#555568", "#a0a08b", "#e9efec"},
	"mist-gb":                {"#2d1b00", "#1e606e", "#5ab9a8", "#c4f0c2"},
	"rustic-gb":              {"#2c2137", "#764462", "#edb4a1", "#a96868"},
	"2-bit-grayscale":        {"#000000", "#676767", "#b6b6b6", "#ffffff"},
	"hollow":                 {"#0f0f1b", "#565a75", "#c6b7be", "#fafbf6"},
	"ayy4":                   {"#00303b", "#ff7777", "#ffce96", "#f1f2da"},
	"nintendo-gameboy-bgb":   {"#081820", "#346856", "#88c070", "#e0f8d0"},
	"red-brick":              {"#eff9d6", "#ba5044", "#7a1c4b", "#1b0326"},
	"nostalgia":              {"#d0d058", "#a0a840", "#708028", "#405010"},
	"spacehaze":              {"#f8e3c4", "#cc3495", "#6b1fb1", "#0b0630"},
	"moonlight-gb":           {"#0f052d", "#203671", "#36868f", "#5fc75d"},
	"links-awakening-sgb":    {"#5a3921", "#6b8c42", "#7bc67b", "#ffffb5"},
	"arq4":                   {"#ffffff", "#6772a9", "#3a3277", "#000000"},
	"blk-aqu4":               {"#002b59", "#005f8c", "#00b9be", "#9ff4e5"},
	"pokemon-sgb":            {"#181010", "#84739c", "#f7b58c", "#ffefff"},
	"nintendo-super-gameboy": {"#331e50", "#a63725", "#d68e49", "#f7e7c6"},
	"blu-scribbles":          {"#051833", "#0a4f66", "#0f998e", "#12cc7f"},
	"kankei4":                {"#ffffff", "#f42e1f", "#2f256b", "#060608"},
	"dark-mode":              {"#212121", "#454545", "#787878", "#a8a5a5"},
	"ajstarks":               {"#aa0000", "#aaaaaa", "#000000", "#ffffff"},
	"pen-n-paper":            {"#e4dbba", "#a4929a", "#4f3a54", "#260d1c"},
	"autumn-decay":           {"#313638", "#574729", "#975330", "#c57938", "#ffad3b", "#ffe596"},
	"polished-gold":          {"#000000", "#361c1b", "#754232", "#cd894a", "#e6b983", "#fff8bc", "#ffffff", "#2d2433", "#4f4254", "#b092a7"},
	"funk-it-up":             {"#e4ffff", "#e63410", "#a23737", "#ffec40", "#81913b", "#26f675", "#4c714e", "#40ebda", "#394e4e", "#0a0a0a"},
}

const (
	linefmt   = "<line xp1=\"%v\" yp1=\"%v\" xp2=\"%v\" yp2=\"%v\" sp=\"%v\" color=%q/>\n"
	squarefmt = "<rect xp=\"%v\" yp=\"%v\" wp=\"%v\" hr=\"100\" color=%q/>\n"
)

// random returns a random number between a range
func random(min, max float64) float64 {
	return vmap(rand.Float64(), 0, 1, min, max)
}

// vmap maps one interval to another
func vmap(value, low1, high1, low2, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

// csquare makes a square with lines, using a specified width and color
// if a hue range is set, the color is randomly selected in that range,
// if a palette is specified, use a random color from it,
// otherwise, the named color is used.
func csquare(x, y, size, maxlw, h1, h2 float64, color string) {

	if c, ok := palette[color]; ok { // use a palette
		color = c[rand.Intn(len(c))]
	}
	if h1 > -1 && h2 > -1 { // hue range set
		color = fmt.Sprintf("hsv(%v,100,100)", random(h1, h2))
	}
	// define the corners
	hs := size / 2
	tlx, tly := x-hs, y+hs
	trx, try := x+hs, y+hs
	blx, bly := x-hs, y-hs
	brx, bry := x+hs, y-hs

	lw := random(0.1, maxlw)
	// make the boundaries
	hline(tlx, tly, size, lw, color)
	hline(blx, bly, size, lw, color)
	vline(blx, bly, size, lw, color)
	vline(brx, bry, size, lw, color)
	// make the corners
	square(tlx, tly, lw, color)
	square(blx, bly, lw, color)
	square(brx, bry, lw, color)
	square(trx, try, lw, color)
}

// square makes a square
func square(x, y, size float64, color string) {
	fmt.Printf(squarefmt, x, y, size, color)
}

// hline makes a horizontal line
func hline(x, y, size, lw float64, color string) {
	fmt.Printf(linefmt, x, y, x+size, y, lw, color)
}

// vline makes a vertical line
func vline(x, y, size, lw float64, color string) {
	fmt.Printf(linefmt, x, y, x, y+size, lw, color)
}

// desordres makes a series of concentric squares
func desordres(x, y, minsize, maxsize, maxlw, h1, h2 float64, color string) {
	step := random(1, 5)
	for v := minsize; v < maxsize; v += step {
		csquare(x, y, v, maxlw, h1, h2, color)
	}
}

// parseHues parses a color string: if the string is of the form "h1:h2",
// where h1, and h2 are numbers between 0 and 360, they are a range of hues.
// Otherwise, set to -1 for invalid entries (use named colors instead)
func parseHues(color string) (float64, float64) {
	var h1, h2 float64 = -1.0, -1.0
	hb := strings.Split(color, ":")
	if len(hb) == 2 {
		var err error
		h1, err = strconv.ParseFloat(hb[0], 64)
		if err != nil {
			h1 = -1
		}
		h2, err = strconv.ParseFloat(hb[1], 64)
		if err != nil {
			h2 = -1
		}
	}
	return h1, h2
}

func usage() {
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "Option      Default    Description\n")
	fmt.Fprintf(os.Stderr, ".....................................................\n")
	fmt.Fprintf(os.Stderr, "-tiles      10          number of tiles/row\n")
	fmt.Fprintf(os.Stderr, "-maxlw      1           maximim line thickness\n")
	fmt.Fprintf(os.Stderr, "-bgcolor    white       background color\n")
	fmt.Fprintf(os.Stderr, "-p          \"\"          palette file\n")
	fmt.Fprintf(os.Stderr, "-color      gray        color name, h1:h2, or palette:\n\n")
	for p, k := range palette {
		fmt.Fprintf(os.Stderr, "%-20s\t%v\n", p, k)
	}
	os.Exit(1)
}

func loadpalette(pfile string) error {
	if len(pfile) > 0 {
		r, err := os.Open(pfile)
		if err != nil {
			return err
		}
		palette, err = ReadString(r)
		if err != nil {
			return err
		}
	}
	return nil
}

// slide generation functions
func beginDeck()              { fmt.Println("<deck>") }
func endDeck()                { fmt.Println("</deck>") }
func beginSlide(color string) { fmt.Printf("<slide bg=%q>\n", color) }
func endSlide()               { fmt.Println("</slide>") }

func main() {
	var tiles, maxlw float64
	var bgcolor, color, pfile string
	var showhelp bool

	flag.Float64Var(&tiles, "tiles", 10, "tiles/row")
	flag.Float64Var(&maxlw, "maxlw", 1, "maximum line thickness")
	flag.StringVar(&bgcolor, "bgcolor", "white", "background color")
	flag.StringVar(&color, "color", "gray", "pen color: (named color, hue range (h1:h2), or palette name")
	flag.StringVar(&pfile, "p", "", "palette file")
	flag.BoolVar(&showhelp, "help", false, "show usage")
	flag.Parse()
	h1, h2 := parseHues(color) // set hue range, or named color/palette

	err := loadpalette(pfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if showhelp {
		usage()
	}

	size := 100 / tiles     // size of each tile
	top := 100 - (size / 2) // top of the beginning row
	left := 100 - top       // left of the beginning row

	beginDeck()
	beginSlide(bgcolor)
	for y := top; y > 0; y -= size {
		for x := left; x < 100; x += size {
			desordres(x, y, 2, size, maxlw, h1, h2, color)
		}
	}
	endSlide()
	endDeck()
}
