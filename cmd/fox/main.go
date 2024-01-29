// fox -- in the style of "Fox I" by Anni Albers
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
const defxs = 0.5
const defys = -0.5
const defop = 40
const rangefmt = "%v,%v,%v"

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

// triangle makes a colored triangle pointing to the specified direction
// ((u)p, (d)own, (l)eft, (r)ight, ne, nw, se, sw)
func triangle(x, y, width, height float64, color string, opacity float64, hue1, hue2 float64, direction string) {
	var xp0, xp1, xp2, yp0, yp1, yp2 float64
	w2 := width / 2
	h2 := height / 2
	switch direction {
	case "n": // up
		xp0, xp1, xp2 = x, x-w2, x+w2
		yp0, yp1, yp2 = y+h2, y-h2, y-h2
	case "s": // down
		xp0, xp1, xp2 = x, x-w2, x+w2
		yp0, yp1, yp2 = y-h2, y+h2, y+h2
	case "e": // left
		xp0, xp1, xp2 = x-w2, x+w2, x+w2
		yp0, yp1, yp2 = y, y+h2, y-h2
	case "w": // right
		xp0, xp1, xp2 = x+w2, x-w2, x-w2
		yp0, yp1, yp2 = y, y+h2, y-h2
	case "ne": // northeast
		xp0, xp1, xp2 = x-w2, x-w2, x+w2
		yp0, yp1, yp2 = y-h2, y+h2, y+h2
	case "nw": // northwest
		xp0, xp1, xp2 = x-w2, x+w2, x-w2
		yp0, yp1, yp2 = y-h2, y+h2, y+h2
	case "sw": // southwest
		xp0, xp1, xp2 = x+w2, x-w2, x-w2
		yp0, yp1, yp2 = y-h2, y-h2, y+h2
	case "se": // southeast
		xp0, xp1, xp2 = x-w2, x+w2, x+w2
		yp0, yp1, yp2 = y-h2, y-h2, y+h2
	}
	if hue1 > -1 && hue2 > -1 { // use hue
		color = fmt.Sprintf("hsv(%v,100,100)", random(hue1, hue2))
	}
	if c, ok := palette[color]; ok { // use a palette
		color = c[rand.Intn(len(c))]
	}
	fmt.Printf("<polygon xc=\"%v %v %v\" yc=\"%v %v %v\" color=%q opacity=\"%.3f\"/>\n", xp0, xp1, xp2, yp0, yp1, yp2, color, opacity)
}

// usage prints usage info
func usage() {
	defrange := fmt.Sprintf(rangefmt, minbound, maxbound, defaultstep)
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "Option    Default               Description\n")
	fmt.Fprintf(os.Stderr, "..................................................................\n")
	fmt.Fprintf(os.Stderr, "-help     false                 show usage\n")
	fmt.Fprintf(os.Stderr, "-w        "+defrange+"               percent begin,end,step for the width\n")
	fmt.Fprintf(os.Stderr, "-h        "+defrange+"               percent begin,end,step for the height\n")
	fmt.Fprintf(os.Stderr, "-shadow   40                    shadow opacity,xoffset,ysoffset\n")
	fmt.Fprintf(os.Stderr, "-d        \"n s e w nw sw ne se\" shape directions\n")
	fmt.Fprintf(os.Stderr, "-xshift   0.5                   shadow x shift\n")
	fmt.Fprintf(os.Stderr, "-yshift   -0.5                  shadow y shift\n")
	fmt.Fprintf(os.Stderr, "-bgcolo   white                 background color\n")
	fmt.Fprintf(os.Stderr, "-p        \"\"                    palette file\n")
	fmt.Fprintf(os.Stderr, "-color    gray                  color name, hue range (h1:h2), or palette:\n\n")
	fmt.Fprintln(os.Stderr, "Palette Name                    Colors\n..........................................................")
	for p, k := range palette {
		fmt.Fprintf(os.Stderr, "%-25s\t%v\n", p, k)
	}
	os.Exit(1)
}

// slide generation functions
func beginDeck()              { fmt.Println("<deck>") }
func endDeck()                { fmt.Println("</deck>") }
func beginSlide(color string) { fmt.Printf("<slide bg=%q>\n", color) }
func endSlide()               { fmt.Println("</slide>") }

func userpalette(pfile string) {
	if len(pfile) > 0 {
		var err error
		palette, err = LoadPalette(pfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	}
}

func setdir(s string) []string {
	d := strings.Fields(s)
	return d
}

func main() {
	// options
	var showhelp bool
	var bgcolor, color, xconfig, yconfig, pfile, dirs string
	var shadowop, xshift, yshift float64
	defrange := fmt.Sprintf(rangefmt, minbound, maxbound, defaultstep)
	flag.BoolVar(&showhelp, "help", false, "show usage")
	flag.Float64Var(&shadowop, "shadow", 40, "shadow opacity (0 for no shadow shape)")
	flag.Float64Var(&xshift, "xshift", 0.5, "shadow x shift")
	flag.Float64Var(&yshift, "yshift", -0.5, "shadow y shift")
	flag.StringVar(&xconfig, "w", defrange, "horizontal config (min,max,step)")
	flag.StringVar(&yconfig, "h", defrange, "vertical config (min,max,step)")
	flag.StringVar(&bgcolor, "bgcolor", "white", "background color")
	flag.StringVar(&dirs, "d", "n s e w sw se nw ne", "directions")
	flag.StringVar(&pfile, "p", "", "palette file")
	flag.StringVar(&color, "color", "gray", "pen color; named color, palette, or h1:h2 for a random hue range hsv(h1:h2, 100, 100)")
	flag.Parse()

	userpalette(pfile)

	if showhelp {
		usage()
	}

	directions := setdir(dirs)
	h1, h2 := parseHues(color)
	bx, ex, xstep := parserange(xconfig)
	by, ey, ystep := parserange(yconfig)
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
				triangle(x+xshift, y+yshift, w, h, color, shadowop, h1, h2, directions[rand.Intn(nd)])
			}
		}
	}
	endSlide()
	endDeck()

}
