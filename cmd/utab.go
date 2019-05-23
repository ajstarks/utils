// utab -- print a unicode font glyph table
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jung-kurt/gofpdf"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "specify fontfile begin end\n")
		os.Exit(1)
	}
	fontname := os.Args[1]
	begin, _ := strconv.ParseInt(os.Args[2], 0, 32)
	end, _ := strconv.ParseInt(os.Args[3], 0, 32)

	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.AddUTF8Font("font", "", fontname)

	fontSize := 24.0
	left := 20.0
	top := 20.0
	right := 200.0
	bottom := 250.0
	footer := bottom + 20.0
	colsize := 18.0

	pdf.AddPage()
	x, y := left, top
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
	pdf.Text(left, footer, fontname)
	if err := pdf.OutputFileAndClose("utab.pdf"); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}
