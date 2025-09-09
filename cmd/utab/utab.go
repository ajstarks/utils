// utab -- print a unicode font glyph table
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"codeberg.org/go-pdf/fpdf"
)

func main() {
	la := len(os.Args)
	outfile := "utab.pdf"
	var fontname string
	var begin, end int64
	begin, end = 0, 255
	// fill in parameters from the command line
	if la > 1 { // utab file
		fontname = os.Args[1]
	}
	if la > 2 { // utab file begin
		begin, _ = strconv.ParseInt(os.Args[2], 0, 32)
	}
	if la > 3 { // utab file begin end
		end, _ = strconv.ParseInt(os.Args[3], 0, 32)
	}
	if la > 4 { // utab file begin end output
		outfile = os.Args[4]
	}
	// check for usage errors
	if begin >= end || len(fontname) == 0 || la <= 1 {
		fmt.Fprintf(os.Stderr, "Usage: utab fontname [begin (default=%d)] [end (default=%d)] [output default=%q]\n", begin, end, outfile)
		os.Exit(1)
	}
	// begin the document
	pdf := fpdf.New("P", "mm", "Letter", "")
	pdf.SetFontLocation(filepath.Dir(fontname))
	pdf.AddUTF8Font("font", "", filepath.Base(fontname))

	// set page parameters
	fontSize := 24.0
	left := 20.0
	top := 20.0
	right := 200.0
	bottom := 250.0
	footer := bottom + 20.0
	colsize := 18.0

	pdf.AddPage()
	x, y := left, top

	// for the specified range, make a font table, making new pages as needed.
	for i := begin; i <= end; i++ {
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("font", "", fontSize)
		pdf.Text(x, y, string(i))
		pdf.SetTextColor(127, 0, 0)
		pdf.SetFont("courier", "", 10)
		pdf.Text(x, y+5, fmt.Sprintf("%05x", i))
		x += fontSize * 1.2
		if x > right {
			x = left
			y += colsize
		}
		if y > bottom {
			pdf.Text(left, footer, fontname)
			pdf.AddPage()
			x, y = left, top
		}
	}
	// write the PDF
	pdf.Text(left, footer, fontname)
	if err := pdf.OutputFileAndClose(outfile); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
	os.Exit(0)
}
