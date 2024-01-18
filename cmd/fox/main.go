// fox -- in the stype of "Fox I" by Anni Albers
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

const minbound = 10
const maxbound = 95
const minstep = 2.0
const maxstep = 20.0
const defaultstep = 5.0
const shadowshift = 0.5
const rangefmt = "%v,%v,%v"

var palette = map[string][]string{
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

// random returns a random number between a range
func random(min, max float64) float64 {
	return vmap(rand.Float64(), 0, 1, min, max)
}

// vmap maps one interval to another
func vmap(value, low1, high1, low2, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
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

// parserange returns the string "v1,v2.v3" as v1, v2, v3
func parserange(s string) (float64, float64, float64) {
	v := strings.Split(s, ",")
	if len(v) == 3 {
		min, err := strconv.ParseFloat(v[0], 64)
		if err != nil {
			min = minbound
		}
		max, err := strconv.ParseFloat(v[1], 64)
		if err != nil {
			max = maxbound
		}
		step, err := strconv.ParseFloat(v[2], 64)
		if err != nil {
			step = defaultstep
		}
		return min, max, step
	}
	return minbound, maxbound, defaultstep
}

// triangle makes a colored triangle pointing to the specified direction ((u)p, (d)own, (l)eft, (r)ight)
func triangle(x, y, width, height float64, color string, opacity float64, hue1, hue2 float64, direction string) {
	var xp, yp [3]float64
	w2 := width / 2
	h2 := height / 2
	switch direction {
	case "u": // up
		xp[0], xp[1], xp[2] = x, x-w2, x+w2
		yp[0], yp[1], yp[2] = y+h2, y-h2, y-h2
	case "d": // down
		xp[0], xp[1], xp[2] = x, x-w2, x+w2
		yp[0], yp[1], yp[2] = y-h2, y+h2, y+h2
	case "l": // left
		xp[0], xp[1], xp[2] = x-w2, x+w2, x+w2
		yp[0], yp[1], yp[2] = y, y+h2, y-h2
	case "r": // right
		xp[0], xp[1], xp[2] = x+w2, x-w2, x-w2
		yp[0], yp[1], yp[2] = y, y+h2, y-h2
	}

	if hue1 > -1 && hue2 > -1 { // use hue
		color = fmt.Sprintf("hsv(%v,100,100)", random(hue1, hue2))
	}
	if c, ok := palette[color]; ok { // use a palette
		color = c[rand.Intn(len(c))]
	}
	fmt.Printf("<polygon xc=\"%v %v %v\" yc=\"%v %v %v\" color=%q opacity=\"%.3f\"/>\n", xp[0], xp[1], xp[2], yp[0], yp[1], yp[2], color, opacity)
}

// usage prints usage info
func usage() {
	defrange := fmt.Sprintf(rangefmt, minbound, maxbound, defaultstep)
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "Option      Default     Description\n")
	fmt.Fprintf(os.Stderr, "..........................................................\n")
	fmt.Fprintf(os.Stderr, "-help       false       show usage\n")
	fmt.Fprintf(os.Stderr, "-width      1000        canvas width\n")
	fmt.Fprintf(os.Stderr, "-height     1000        canvas height\n")
	fmt.Fprintf(os.Stderr, "-w          "+defrange+"     percent begin,end,step for the width\n")
	fmt.Fprintf(os.Stderr, "-h          "+defrange+"     percent begin,end,step for the height\n")
	fmt.Fprintf(os.Stderr, "-bgcolor    white       background color\n")
	fmt.Fprintf(os.Stderr, "-color      gray        color name, h1:h2, or palette:\n\n")
	fmt.Fprintln(os.Stderr, "Palette Name            Colors\n..........................................................")
	for p, k := range palette {
		fmt.Fprintf(os.Stderr, "%-20s\t%v\n", p, k)
	}
	os.Exit(1)
}

// slide generation functions
func beginDeck()              { fmt.Println("<deck>") }
func endDeck()                { fmt.Println("</deck>") }
func beginSlide(color string) { fmt.Printf("<slide bg=%q>\n", color) }
func endSlide()               { fmt.Println("</slide>") }

func main() {
	// options
	var showhelp bool
	var bgcolor, color, xconfig, yconfig string
	var shadowop float64
	defrange := fmt.Sprintf(rangefmt, minbound, maxbound, defaultstep)
	flag.BoolVar(&showhelp, "help", false, "show usage")
	flag.Float64Var(&shadowop, "shadow", 40, "shadow opacity (0 for no shadow shape)")
	flag.StringVar(&xconfig, "w", defrange, "horizontal config (min,max,step)")
	flag.StringVar(&yconfig, "h", defrange, "vertical config (min,max,step)")
	flag.StringVar(&bgcolor, "bgcolor", "white", "background color")
	flag.StringVar(&color, "color", "gray", "pen color; named color, palette, or h1:h2 for a random hue range hsv(h1:h2, 100, 100)")
	flag.Parse()
	if showhelp {
		usage()
	}
	bx, ex, xstep := parserange(xconfig)
	by, ey, ystep := parserange(yconfig)
	h1, h2 := parseHues(color)
	directions := []string{"u", "d", "l", "r"}
	nd := len(directions)

	// generation
	beginDeck()
	beginSlide(bgcolor)
	for y := by; y < ey; y += ystep {
		for x := bx; x < ex; x += xstep {
			w := random(minstep, xstep)
			h := random(minstep, ystep)
			triangle(x, y, w, h, color, 100, h1, h2, directions[rand.Intn(nd)])
			if shadowop > 0 {
				triangle(x+shadowshift, y-shadowshift, w, h, color, shadowop, h1, h2, directions[rand.Intn(nd)])
			}
		}
	}
	endSlide()
	endDeck()

}
