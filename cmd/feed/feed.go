// feed: process "Friday Feed" files, output as plain text, deck, html, RTF, or JSON
package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ajstarks/deck/generate"
)

// Feed has a title, date, with a series of entries
type Feed struct {
	Title   string  `xml:"title"`
	Date    string  `xml:"date"`
	Entries []entry `xml:"entry"`
}

// entry consists of a title, quote and link
type entry struct {
	Title string   `xml:"title"`
	Quote string   `xml:"quote"`
	Link  []string `xml:"link"`
}

var (
	outfmt      = flag.String("f", "deck", "output format (deck, rtf, json, html, or plain)")
	leftmargin  = flag.Float64("left", 15.0, "left margin")
	rightmargin = flag.Float64("right", 60.0, "right margin")
	fontsize    = flag.Float64("fs", 1.8, "font size")
)

const (
	ecolor    = "black"            // entry color
	tcolor    = "rgb(127,0,0)"     // title color
	qcolor    = "rgb(127,127,127)" // quote color
	htmltop   = `<html><head><style type="text/css">body {font-size:14pt;font-family: Calibri, Arial, sans-serif; margin-left:10%; margin-right:10%;} h1 {font-size:200%;} h2 {margin-top:20pt;font-size: 120%;} p {margin-top: 4pt; margin-bottom:3pt;}</style></head><body>`
	htmldate  = "<h1>%s</h1><p>%s</p>"
	htmltitle = "<h2>%s</h2>\n"
	htmlquote = "<p>%s</p>\n"
	htmlend   = "</body>\n</html>"
	htmllink  = "<a href=\"%s\">%s</a>\n"
	rtfhead   = "{\\rtf1\\ansi\\ansicpg1252{\\fonttbl\\f0\\fnil\\fcharser0 Calibri;}{\\colortbl ;\\red0\\green0\\blue238;}"
	tfmt      = "\\f0\\b\\fs28 %s\\\n\n"
	qfmt      = "\\i\\b0 %s\\\n"
	hfmt      = "\\i0{\\field{\\*\\fldinst HYPERLINK \"%s\"}{\\fldrslt {\\ul\\cf1%s}}}\\\n\\ulnone\n\\\n\\\n"
)

// map utf-8 to windows-1252 notation
var unicodemap = strings.NewReplacer(
	"\u2018", "\\'91", "\u2019", "\\'92", "\u201c", "\\'93", "\u201d", "\\'94", "\u2022", "\\'95", "\u2013", "\\'96",
	"\u2014", "\\'97", "\u2122", "\\'99", "\u20ac", "\\'80", "\u2026", "\\'85", "\u00b6", "\\'b6", "\u00a7", "\\'a7",
	"\u00a9", "\\'a9", "\u00ae", "\\'ae", "\u00b0", "\\'b0", "\u0192", "\\'c0", "\u0193", "\\'c1", "\u0194", "\\'c2",
	"\u0195", "\\'c3", "\u0196", "\\'c4", "\u0197", "\\'c5", "\u0198", "\\'c6", "\u0199", "\\'c7", "\u0200", "\\'c8",
	"\u0201", "\\'c9", "\u0202", "\\'ca", "\u0203", "\\'cb", "\u0204", "\\'cc", "\u0204", "\\'cd", "\u0206", "\\'ce",
	"\u0207", "\\'cf", "\u0208", "\\'d0", "\u0209", "\\'d1", "\u0210", "\\'d2", "\u0211", "\\'d3", "\u0212", "\\'d4",
	"\u0213", "\\'d5", "\u0214", "\\'d6", "\u0215", "\\'d7", "\u0216", "\\'d8", "\u0217", "\\'d9", "\u0218", "\\'da",
	"\u0219", "\\'db", "\u0220", "\\'dc", "\u0221", "\\'dd", "\u0222", "\\'de", "\u0223", "\\'df", "\u0224", "\\'e0",
	"\u0225", "\\'e1", "\u0226", "\\'e2", "\u0227", "\\'e3", "\u0228", "\\'e4", "\u0229", "\\'e5", "\u0230", "\\'e6",
	"\u0231", "\\'e7", "\u0232", "\\'e8", "\u0233", "\\'e9", "\u0234", "\\'ea", "\u0235", "\\'eb", "\u0236", "\\'ec",
	"\u0237", "\\'ed", "\u0238", "\\'ee", "\u0239", "\\'ef", "\u0240", "\\'f0", "\u0241", "\\'f1", "\u0242", "\\'f2",
	"\u0243", "\\'f3", "\u0244", "\\'f4", "\u0245", "\\'f5", "\u0246", "\\'f6", "\u0247", "\\'f7", "\u0248", "\\'f8",
	"\u0249", "\\'f9", "\u0250", "\\'fa", "\u0251", "\\'fb", "\u0252", "\\'fc", "\u0253", "\\'fd", "\u0254", "\\'fe",
	"\u0255", "\\'ff",
)

// unitranslate converts unicode to RTF escapes
func unitranslate(s string) string {
	return unicodemap.Replace(s)
}

var xmlmap = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;")

func xmltranslate(s string) string {
	return xmlmap.Replace(s)
}

// genplain outputs the feed markup as plain text
func genplain(w io.Writer, f Feed) {
	fmt.Fprintf(w, "%s: %s\n\n", f.Title, f.Date)
	for _, e := range f.Entries {
		fmt.Fprintf(w, "%s\n\n", e.Title)
		fmt.Fprintf(w, "%s\n", e.Quote)
		for _, l := range e.Link {
			fmt.Fprintf(w, "%s\n", l)
		}
		fmt.Fprintf(w, "\n\n")
	}
}

// genrtf converts the feed to RTF
func genrtf(w io.Writer, f Feed) {
	fmt.Fprintf(w, rtfhead)
	for _, e := range f.Entries {
		fmt.Fprintf(w, tfmt, unitranslate(e.Title))
		fmt.Fprintf(w, qfmt, unitranslate(e.Quote))
		for _, l := range e.Link {
			fmt.Fprintf(w, hfmt, l, l)
		}
	}
	fmt.Fprintln(w, "}")
}

// genhtml outputs the feed markup as HTML
func genhtml(w io.Writer, f Feed) {
	fmt.Fprintln(w, htmltop)
	fmt.Fprintf(w, htmldate, f.Title, f.Date)
	for _, e := range f.Entries {
		fmt.Fprintf(w, htmltitle, e.Title)
		fmt.Fprintf(w, htmlquote, e.Quote)
		for _, l := range e.Link {
			fmt.Fprintf(w, htmllink, l, l)
		}
	}
	fmt.Fprintln(w, htmlend)
}

// genjson outputs the feed markup as JSON
func genjson(w io.Writer, f Feed) {
	b, err := json.Marshal(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	fmt.Fprintf(w, "%v\n", string(b))
}

// gendeck outputs the feed markup as deck markup
func gendeck(d *generate.Deck, f Feed) {
	// set text locations
	top := 90.0
	x := *leftmargin
	y := top
	fs := *fontsize
	right := *rightmargin
	bottom := 15.0

	d.StartSlide()
	d.Text(x, y, f.Title, "sans", fs*1.5, tcolor)
	d.TextEnd(right+x+fs, y, f.Date, "sans", fs, tcolor)

	y -= fs * 4.0
	for _, e := range f.Entries {
		// check for slide overflow
		if y < bottom {
			d.EndSlide()
			d.StartSlide()
			y = top
		}
		// do title with the first link
		// subsequent links are smaller
		for il, l := range e.Link {
			if il == 0 {
				d.TextLink(x, y, xmltranslate(e.Title), l, "sans", fs, ecolor)
			} else {
				y -= fs
				d.TextLink(x, y, "See also:"+xmltranslate(l), l, "sans", fs/2, "rgb(127,0,0)")
			}
		}
		y -= (fs * 1.5)
		d.TextBlock(x, y, xmltranslate(e.Quote), "serif", fs*0.8, right, qcolor)
		y -= (quoteskip(e.Quote) * (fs)) + (fs * 3.2)
	}
	d.EndSlide()
}

// compute spacing based on the size of the quote
func quoteskip(s string) float64 {
	fl := float64(len(s)) / 100.0
	i := float64(int(fl))
	d := fl - i
	if d > 0.2 {
		fl = fl + 1
	}
	return float64(int(fl))
}

// process each specifed file
func main() {
	flag.Parse()
	var d *generate.Deck
	for i, filename := range flag.Args() {
		f, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		var data Feed
		err = xml.NewDecoder(f).Decode(&data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		switch *outfmt {
		case "rtf":
			genrtf(os.Stdout, data)
		case "plain":
			genplain(os.Stdout, data)
		case "html":
			genhtml(os.Stdout, data)
		case "json":
			genjson(os.Stdout, data)
		case "deck":
			if i == 0 {
				d = generate.NewSlides(os.Stdout, 0, 0)
				d.StartDeck()
			}
			gendeck(d, data)
			if i == len(flag.Args())-1 {
				d.EndDeck()
			}
		}
		f.Close()
	}
}
