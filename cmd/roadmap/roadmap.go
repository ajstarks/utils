// roadmap: create roadmaps from a XML description
package main

import (
	"encoding/hex"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/ajstarks/gensvg"
)

// Geometry defines the dimensions of objects
type Geometry struct {
	x, y, w, h float64
}

// Roadmap describes the structure of the roadmap
type Roadmap struct {
	Title      string     `xml:"title,attr"`
	Begin      float64    `xml:"begin,attr"`
	End        float64    `xml:"end,attr"`
	Scale      float64    `xml:"scale,attr"`
	Catpercent float64    `xml:"catpercent,attr"`
	Vspace     float64    `xml:"vspace,attr"`
	Itemheight float64    `xml:"itemheight,attr"`
	Fontname   string     `xml:"fontname,attr"`
	Shape      string     `xml:"shape,attr"`
	Category   []Category `xml:"category"`
}

// Category defines categories within roadmaps
type Category struct {
	Name       string  `xml:"name,attr"`
	Color      string  `xml:"color,attr"`
	Vspace     string  `xml:"vspace,attr"`
	Itemheight string  `xml:"itemheight,attr"`
	Shape      string  `xml:"shape,attr"`
	Bline      string  `xml:"bline,attr"`
	Catdesc    Catdesc `xml:"catdesc"`
	Item       []Item  `xml:"item"`
}

// Catdesc defines category description items
type Catdesc struct {
	Cditem []cditem `xml:"cditem"`
}

type cditem struct {
	Cdtext string `xml:",chardata"`
}

// Dep defines dependencies within items
type Dep struct {
	Dest string  `xml:"dest,attr"`
	Desc string  `xml:",chardata"`
	BPct float64 `xml:"bpct,attr"`
	EPct float64 `xml:"epct,attr"`
}

// Item defines items within categories
type Item struct {
	Id        string `xml:"id,attr"`
	Begin     string `xml:"begin,attr"`
	Duration  string `xml:"duration,attr"`
	Color     string `xml:"color,attr"`
	Milestone string `xml:"milestone,attr"`
	Shape     string `xml:"shape,attr"`
	Align     string `xml:"align,attr"`
	Vspace    string `xml:"vspace,attr"`
	Bline     string `xml:"bline,attr"`
	Text      string `xml:",chardata"`
	Dep       []Dep  `xml:"dep"`
	Desc      string `xml:"desc"`
	X         float64
	Y         float64
	W         float64
	H         float64
}

// command line variables
var (
	width       = flag.Float64("w", 1024, "width")
	height      = flag.Float64("h", 768, "height")
	lmargin     = flag.Float64("margin", 10, "margin")
	twrap       = flag.Float64("wrap", 20, "text wrap")
	tfs         = flag.Float64("tfs", 24, "title font size (px)")
	cfs         = flag.Float64("cfs", 14, "category font size (px)")
	ifs         = flag.Float64("ifs", 12, "item fontsize (px)")
	lalign      = flag.String("align", "end", "label alignment")
	topborder   = flag.Bool("tb", false, "top border")
	botborder   = flag.Bool("bb", false, "bottom border")
	leftborder  = flag.Bool("lb", true, "left border")
	rightborder = flag.Bool("rb", false, "right border")
	catborder   = flag.Bool("cb", false, "category border")
	boldcat     = flag.Bool("b", false, "bold categories")
	descend     = flag.Bool("de", true, "description at the end of the item")
	bgcolor     = flag.String("bg", "white", "background color")
	curves      = flag.String("curves", "0,0", "curve line connections")
	csv         = flag.String("csv", "", "write CSV to specified file")
	descolor    = flag.String("dc", "red", "description color")
	concolor    = flag.String("cc", "red", "connection color")
)

const (
	borderfmt    = "stroke:#BBBBBBCC;stroke-width:0.75px"
	twrapfmt     = "font-style:italic;text-anchor:%s;fill-opacity:%.2f;fill:%s;font-family:%s;font-size:%fpx"
	depfmt       = "stroke-width:2;stroke:%s;stroke-opacity:0.6;stroke-dasharray:2 2;fill:none"
	ccfmt        = "stroke:none;fill:%s;fill-opacity:0.3"
	connectfmt   = "stroke:none;text-anchor:middle;font-style:italic;fill:%s;font-size:60%%"
	catdescfmt   = "text-anchor:start;fill:red;font-size:%.2fpx"
	itemtextfmt  = "text-anchor:%s;fill:%s;font-size:%.2fpx"
	itemlinefmt  = "stroke:%s;stroke-width:%.2fpx"
	hexfillfmt   = "fill:%s;fill-opacity:%.2f"
	boldfmt      = "font-weight:bold"
	italicfmt    = "font-style:italic"
	strokefmt    = "stroke:%s;%s"
	categoryfmt  = "text-anchor:start;font-size:%.2fpx"
	catgstylefmt = "text-anchor:middle;font-family:"
)

// roadmap reads and processes a roadmap XML file
func roadmap(location string, canvas *gensvg.SVG) {
	var f *os.File
	var err error
	f = os.Stdin
	if len(location) > 0 {
		f, err = os.Open(location)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}
	}
	readrm(f, canvas)
	f.Close()
}

// readrm reads in the roadmap struct
func readrm(r io.Reader, canvas *gensvg.SVG) {
	var rm Roadmap
	err := xml.NewDecoder(r).Decode(&rm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	drawrm(rm, canvas)
	if len(*csv) > 0 {
		csvfile, err := os.Create(*csv)
		if err != err {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}
		rmcsv(rm, csvfile)
		csvfile.Close()
	}
}

// rmcsv renders the roadmap data as CSV
func rmcsv(r Roadmap, w io.Writer) {
	fmt.Fprintln(w, "\"Category/Item\",Begin,Duration,Id")
	for _, cat := range r.Category {
		if cat.Vspace == "0" {
			continue
		}
		fmt.Fprintf(w, ",,,\n\"%s\",,,\n", cat.Name)
		for _, item := range cat.Item {
			fmt.Fprintf(w, "\"%s\",%s,%s,%s\n", item.Text, item.Begin, item.Duration, item.Id)
		}
	}
}

// readrm reads in the roadmap struct
func drawrm(r Roadmap, canvas *gensvg.SVG) {
	var (
		itemshape = "r"
		itemalign = "middle"
		itemcolor = "#BBBBBB88"
		fontname  = "Calibri,sans-serif"
		itemMargin, rightMargin, itemheight, itemyear, itempart,
		itemvspace, catheight, cvspace, top, catx float64
		bline = false
	)

	// Global roadmap attributes
	beginyear := r.Begin
	endyear := r.End
	yearscale := r.Scale
	rvspace := r.Vspace
	ritemheight := r.Itemheight
	if len(r.Fontname) > 0 {
		fontname = r.Fontname
	}

	itemMargin = *width * (r.Catpercent / 100.0)
	rightMargin = *width * 0.98
	tloc := 30.0    // int(float(*height) * 0.05)
	top = tloc + 10 // int(float(*height) * 0.10)
	y := top
	milestone := false

	if *lalign == "end" {
		catx = itemMargin - *lmargin
	} else {
		catx = *lmargin
	}

	canvas.Start(*width, *height)
	canvas.Title(r.Title)
	canvas.Rect(0, 0, *width, *height, "fill:"+*bgcolor)

	canvas.Gstyle(catgstylefmt + fontname)
	canvas.Text(itemMargin, tloc, r.Title, fmt.Sprintf(categoryfmt, *tfs))

	// Process Categories
	for cc, cat := range r.Category {
		if len(cat.Itemheight) == 0 {
			catheight = ritemheight
		} else {
			catheight, _ = strconv.ParseFloat(cat.Itemheight, 64)
		}

		var ycatlabel float64
		if len(cat.Name) > 0 {
			label := strings.Split(cat.Name, "\\n")
			ll := len(label)
			if *boldcat {
				canvas.Gstyle(boldfmt)
			} else {
				canvas.Gstyle(italicfmt)
			}

			if ll <= 1 {
				ycatlabel = y + (catheight / 2) + *cfs/4
			} else {
				ycatlabel = y + ((float64(ll) * *cfs) / 2)
			}
			canvas.Textlines(catx, ycatlabel, label, *cfs, *cfs+2, "black", *lalign)
			canvas.Gend()
		}
		if cat.Bline == "on" {
			canvas.Line(itemMargin, y, rightMargin, y, borderfmt)
		}

		// Process Category descriptions
		yd := ycatlabel + *cfs
		canvas.Gstyle(fmt.Sprintf(catdescfmt, *ifs))
		for _, cdi := range cat.Catdesc.Cditem {
			canvas.Text(*lmargin, yd, cdi.Cdtext)
			yd += *cfs + 2
		}
		canvas.Gend()

		if len(cat.Vspace) == 0 {
			cvspace = rvspace
		} else {
			cvspace, _ = strconv.ParseFloat(cat.Vspace, 64)
		}

		// Process Items within categories
		for ii, item := range cat.Item {
			dt := strings.SplitN(item.Begin, "/", 2)
			if len(dt) == 2 {
				itemyear, _ = strconv.ParseFloat(dt[0], 64)
				itempart, _ = strconv.ParseFloat(dt[1], 64)
			} else {
				continue
			}

			itemduration, _ := strconv.ParseFloat(item.Duration, 64)

			if item.Bline == "on" {
				bline = true
			} else {
				bline = false
			}

			if item.Align == "" {
				itemalign = "middle"
			} else {
				itemalign = item.Align
			}

			if item.Milestone == "on" {
				milestone = true
				itemalign = "end"
			} else {
				milestone = false
			}

			if len(item.Color) > 0 {
				itemcolor = item.Color
			} else {
				if len(cat.Color) > 0 {
					itemcolor = cat.Color
				} else {
					itemcolor = "#BBBBBB88"
				}
			}

			if yearscale <= 0 {
				yearscale = 1.0
			}

			itemfraction := (itempart - 1) / yearscale
			itemx := fmap(itemyear+itemfraction, beginyear, endyear, itemMargin, rightMargin)
			itemw := fmap(itemduration, 0, (endyear-beginyear)*yearscale, 0, rightMargin-itemMargin)

			if len(cat.Shape) == 0 {
				itemshape = r.Shape
			} else {
				itemshape = cat.Shape
			}
			if len(item.Shape) > 0 {
				itemshape = item.Shape
			}
			if len(cat.Itemheight) == 0 {
				itemheight = ritemheight
			} else {
				itemheight = catheight
			}
			if len(cat.Vspace) == 0 {
				itemvspace = 0.0
			} else {
				itemvspace = cvspace
			}
			if len(item.Vspace) > 0 {
				itemvspace, _ = strconv.ParseFloat(item.Vspace, 64)
			}
			drawitem(item.Text, itemx, y, itemw, itemheight, itemshape, itemcolor, itemalign, milestone, bline, canvas)

			if len(item.Desc) > 0 {
				if *descend {
					textwrap(itemx+itemw, y+*ifs, *twrap, *ifs, *ifs+2, item.Desc, fontname, "start", *descolor, 1.0, canvas)
				} else {
					textwrap(itemx-5, y+*ifs, *twrap, *ifs, *ifs+2, item.Desc, fontname, "end", *descolor, 1.0, canvas)
				}
			}

			geo := &r.Category[cc].Item[ii]
			geo.X = itemx
			geo.Y = y
			geo.W = itemw
			geo.H = itemheight

			if ii < len(cat.Item)-1 {
				y += itemvspace
			}
		}

		if len(cat.Vspace) == 0 || itemvspace == 0 {
			y += catheight
		} else {
			y += rvspace
		}

		if *catborder && cc > 1 {
			canvas.Line(itemMargin, y-rvspace, rightMargin, y-rvspace, borderfmt)
		}

	}

	// Process dependencies
	canvas.Gstyle(fmt.Sprintf(depfmt, *concolor))
	for _, c := range r.Category {
		for _, i := range c.Item {
			for _, d := range i.Dep {
				connect(i, d, r.Category, canvas)
			}
		}
	}
	canvas.Gend()

	canvas.Gend()

	if *leftborder {
		canvas.Line(itemMargin, top, itemMargin, *height, borderfmt)
	}
	if *topborder {
		canvas.Line(itemMargin, top, rightMargin, top, borderfmt)
	}
	if *botborder {
		canvas.Line(itemMargin, *height, rightMargin, *height, borderfmt)
	}
	if *rightborder {
		canvas.Line(rightMargin, top, rightMargin, *height, borderfmt)
	}
	canvas.End()
}

// connect matches destinations to make connections
func connect(item Item, d Dep, cats []Category, canvas *gensvg.SVG) {
	var curvex, curvey float64
	fmt.Sscanf(*curves, "%d,%d", &curvex, &curvey)

	for _, c := range cats {
		for _, i := range c.Item {
			if (d.Dest == i.Id) && len(d.Dest) > 0 {
				bx := item.X + item.W*d.BPct
				by := item.Y + item.H/2
				ex := i.X + i.W*d.EPct
				ey := i.Y + i.H/2
				cx := ex + curvex
				cy := ey + curvey
				canvas.Qbez(bx, by, cx, cy, ex, ey)
				canvas.Circle(ex, ey, 4, fmt.Sprintf(ccfmt, *concolor))
				if len(d.Desc) > 0 {
					canvas.Text(ex, ey+item.H+5, d.Desc,
						fmt.Sprintf(connectfmt, *descolor))
				}
			}
		}
	}
}

// whitespace determines if a rune is whitespace
func whitespace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\t'
}

// textwrap draws text at location, wrapping at the specified width
func textwrap(x, y, w, fs, leading float64, s, font, align, color string, opacity float64, canvas *gensvg.SVG) {
	canvas.Gstyle(fmt.Sprintf(twrapfmt, align, opacity, color, font, fs))
	words := strings.FieldsFunc(s, whitespace)
	xp := x
	yp := y
	var line string
	for _, s := range words {
		line += s + " "
		if float64(len(line)) > w {
			canvas.Text(xp, yp, line)
			yp += leading
			line = ""
		}
	}
	if len(line) > 0 {
		canvas.Text(xp, yp, line)
	}
	canvas.Gend()
}

// drawitem renders roadmap items
func drawitem(s string, x, y, w, h float64, shape, color, align string, milestone, bline bool, canvas *gensvg.SVG) {
	var textfill string
	fc := fmt.Sprintf(strokefmt, *bgcolor, hexstyle(color))

	if bline {
		canvas.Line(x+w, y, x+w, *height, borderfmt)
	}

	if len(s) > 0 {
		switch shape {
		case "r":
			canvas.Rect(x, y, w, h, fc)
		case "rr":
			canvas.Roundrect(x, y, w, h, 5, 5, fc)
		case "e":
			canvas.Ellipse(x+(w/2), y+(h/2), w/2, h/2, fc)
		case "c":
			canvas.Circle(x+(w/2), y+(h/2), h/2, fc)
		case "a":
			arrow(x, y, w, h, h/2, fc, canvas)
		case "l":
			yl := (y + (h / 2)) + ((*ifs / 4) - (*ifs / 3))
			canvas.Line(x, yl, x+w, yl, fmt.Sprintf(itemlinefmt, color, *ifs))
		default:
			canvas.Rect(x, y, w, h, fc)
		}
	} else {
		canvas.Rect(x, y, w, h, fc)
	}

	red, green, blue, alpha := colorcomp(color)
	_, _, v := rgbtohsb(red, green, blue)

	if v <= 100.0 && v > 70.0 || alpha < 127 {
		textfill = "black"
	} else {
		textfill = "white"
	}

	if milestone {
		s += " \u2605"
	}
	tx := x
	switch align {
	case "start":
		tx += 1.0 // 5
	case "middle":
		tx += (w / 2)
	case "end":
		tx += (w)
	default:
		tx += (w / 2)
	}
	canvas.Text(tx, (y+(h/2))+(*ifs/4), s, fmt.Sprintf(itemtextfmt, align, textfill, *ifs))

}

// arrow makes an arrow shape
func arrow(x float64, y float64, w float64, h float64, ah float64, color string, canvas *gensvg.SVG) {
	end := x + w
	bot := y + h
	ap := end - ah
	var xp = []float64{x, ap, end, ap, x}
	var yp = []float64{y, y, y + (h / 2), bot, bot}
	canvas.Polyline(xp, yp, color)
}

// rgbtohsb converts an RGB triple to HSB
func rgbtohsb(red float64, green float64, blue float64) (hue, sat, bright float64) {
	hue = 0.0
	minRGB := math.Min(math.Min(red, green), blue)
	maxRGB := math.Max(math.Max(red, green), blue)
	delta := maxRGB - minRGB
	bright = maxRGB
	if maxRGB != 0.0 {
		sat = 255.0 * delta / maxRGB
	} else {
		sat = 0.0
	}
	if sat != 0.0 {
		if red == maxRGB {
			hue = (green - blue) / delta
		}
		if green == maxRGB {
			hue = 2.0 + (blue-red)/delta
		}
		if blue == maxRGB {
			hue = 4.0 + (red-green)/delta
		}
	} else {
		hue = -1.0
	}
	hue = hue * 60
	if hue < 0.0 {
		hue = hue + 360.0
	}
	sat = sat * 100.0 / 255.0
	bright = bright * 100.0 / 255.0
	return hue, sat, bright
}

// hexstyle returns the fill string from a hex representation
func hexstyle(s string) string {
	if len(s) == 9 {
		o, err := hex.DecodeString(s[7:9])
		if err == nil {
			op := float64(o[0]) / 255.0
			return fmt.Sprintf(hexfillfmt, s[0:7], op)
		}
	}
	return "fill:" + s
}

// colorcomp converts a hex string to red, green, blue
func colorcomp(s string) (red, green, blue, alpha float64) {
	if len(s) < 7 && s[0:0] != "#" {
		return red, green, blue, alpha
	}
	var av = make([]byte, 1)
	rv, _ := hex.DecodeString(s[1:3])
	gv, _ := hex.DecodeString(s[3:5])
	bv, _ := hex.DecodeString(s[5:7])

	red = float64(rv[0])
	green = float64(gv[0])
	blue = float64(bv[0])
	alpha = 255.0

	if len(s) == 9 {
		av, _ = hex.DecodeString(s[7:9])
		alpha = float64(av[0])
	}
	return red, green, blue, alpha
}

// wordstack takes a slice of string into a stack of words
func wordstack(x, y, fs float64, s []string, style string, canvas *gensvg.SVG) {
	ls := fs + (fs / 2)
	y -= ls
	for i := len(s); i > 0; i-- {
		canvas.Text(x, y, s[i-1], style)
		y -= ls
	}
}

// fmap maps ranges
func fmap(value, low1, high1, low2, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

func main() {
	canvas := gensvg.New(os.Stdout)
	flag.Parse()
	if len(flag.Args()) > 0 {
		for _, f := range flag.Args() {
			roadmap(f, canvas)
		}
	} else {
		roadmap("", canvas)
	}
}
