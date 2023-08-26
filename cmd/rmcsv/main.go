// rmcsv -- convert csv roadmap files to xml
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	rmfmt        = "<roadmap title=%q font=%q shape=%q begin=\"%d\" end=\"%d\" catpercent=\"%d\" scale=\"%d\" itemheight=\"%d\" vspace=\"%d\">\n"
	endrmfmt     = "\t</category>\n</roadmap>\n"
	yearcatfmt   = "\t<category color=\"#000000\" shape=\"r\" itemheight=\"30\" vspace=\"0\">\n"
	yearitemfmt  = "\t\t<item begin=\"%02d/01\" duration=\"12\" bline=\"on\">%d</item>\n"
	endcat       = "\t</category>\n"
	monthcatfmt  = "\t<category color=\"#bbbbbb\" shape=\"r\" itemheight=\"30\" vspace=\"0\">\n"
	monthitemfmt = "\t\t<item begin=\"%d/%02d\" duration=\"1\">%s</item>\n"
	catfmt       = "\t<category name=%q itemheight=\"%d\" vspace=\"%d\">\n"
	itemfmt1     = "\t\t<item id=%q begin=%q duration=%q><dep dest=%q/>%s</item>\n"
	itemfmt2     = "\t\t<item id=%q begin=%q duration=%q>%s</item>\n"
	itemfmt3     = "\t\t<item begin=%q duration=%q>%s</item>\n"
)

type rmconfig struct {
	title, shape, font                                string
	begin, end, catpercent, scale, vspace, itemheight int
	dh                                                bool
}

// main: process roadmap files on the command line,
// use stdin if no files specified.
func main() {
	var config rmconfig
	flag.StringVar(&config.title, "title", "Title", "Roadmap Title")
	flag.StringVar(&config.shape, "shape", "r", "item shape")
	flag.StringVar(&config.font, "font", "Calibri,sans-serif", "roadmap font")
	flag.IntVar(&config.begin, "begin", 2022, "begin year")
	flag.IntVar(&config.end, "end", 2023, "end year")
	flag.IntVar(&config.catpercent, "cp", 12, "category percent")
	flag.IntVar(&config.scale, "scale", 12, "scale")
	flag.IntVar(&config.vspace, "vspace", 35, "vspace")
	flag.IntVar(&config.itemheight, "itemh", 30, "itemheight")
	flag.BoolVar(&config.dh, "datehead", true, "include date header")

	flag.Parse()
	files := flag.Args()
	nf := len(files)
	if nf == 0 {
		process("", config, os.Stdout)
		return
	}
	for i := 0; i < nf; i++ {
		process(files[i], config, os.Stdout)
	}
}

// process reads and processes a roadmap XML file
func process(location string, config rmconfig, w io.Writer) {
	var r *os.File
	var err error
	r = os.Stdin
	if len(location) > 0 {
		r, err = os.Open(location)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}
	}
	csvtoxml(w, r, config)
	r.Close()
}

// csvtoxml reads the roadmap CSV, converting to XML
func csvtoxml(w io.Writer, r io.Reader, config rmconfig) {
	// roadmap root element
	fmt.Fprintf(w, rmfmt, xmlesc(config.title), xmlesc(config.font), xmlesc(config.shape),
		config.begin, config.end, config.catpercent, config.scale, config.itemheight, config.vspace)

	// read categories and items from csv, write XML
	input := csv.NewReader(r)
	nr := 0
	nc := 0
	if config.dh {
		dateheader(w, config.begin, config.end)
	}
	for {
		fields, csverr := input.Read()
		if csverr == io.EOF {
			break
		}
		if csverr != nil {
			fmt.Fprintf(os.Stderr, "%v %v\n", csverr, fields)
			continue
		}
		nr++
		// skip header
		if nr == 1 {
			continue
		}
		// skip invalid fields
		if len(fields) < 5 {
			continue
		}
		if len(fields[0]) == 0 {
			continue
		}
		// process categories
		if len(fields[0]) > 0 && len(fields[1]) == 0 && len(fields[2]) == 0 {
			nc++
			if nc > 1 {
				fmt.Fprint(w, endcat)
			}
			fmt.Fprintf(w, catfmt, xmlesc(fields[0]), config.itemheight, config.vspace)
			continue
		}
		// process items
		switch {
		case len(fields[3]) > 0 && len(fields[4]) > 0:
			fmt.Fprintf(w, itemfmt1, xmlesc(fields[3]), xmlesc(fields[1]), xmlesc(fields[2]), xmlesc(fields[4]), xmlesc(fields[0]))
		case len(fields[3]) > 0 && len(fields[4]) == 0:
			fmt.Fprintf(w, itemfmt2, xmlesc(fields[3]), xmlesc(fields[1]), xmlesc(fields[2]), xmlesc(fields[0]))
		default:
			fmt.Fprintf(w, itemfmt3, xmlesc(fields[1]), xmlesc(fields[2]), xmlesc(fields[0]))
		}
	}
	// end the XML
	fmt.Fprint(w, endrmfmt)
}

// xmlmap defines the XML substitutions
var xmlmap = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;")

// xmlesc escapes XML
func xmlesc(s string) string {
	return xmlmap.Replace(s)
}

// dateheader makes a date header
func dateheader(w io.Writer, begin, end int) {
	// years
	fmt.Fprint(w, yearcatfmt)
	for y := begin; y < end; y++ {
		fmt.Fprintf(w, yearitemfmt, y, y)
	}
	fmt.Fprint(w, endcat)

	// months
	fmt.Fprint(w, monthcatfmt)
	for y := begin; y < end; y++ {
		for i, m := range []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"} {
			fmt.Fprintf(w, monthitemfmt, y, i+1, m)
		}
	}
	fmt.Fprint(w, endcat)
}
