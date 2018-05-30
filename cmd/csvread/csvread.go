package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
)

// getf turns a string slice of numbers into a slice of integers
func getf(s []string) []int {
	fn := []int{}
	for _, f := range s {
		n, err := strconv.Atoi(f)
		if err != nil {
			continue
		}
		fn = append(fn, n)
	}
	return fn
}

// output displays fields of data
func output(s []string, w *csv.Writer, plain bool) {
	if plain {
		nl := len(s) - 1
		for i := 0; i < nl; i++ {
			fmt.Printf("%s\t", s[i])
		}
		fmt.Println(s[nl])
	} else {
		w.Write(s)
	}
}

func main() {
	var plainout = flag.Bool("plain", true, "plain output")
	var headskip = flag.Bool("headskip", false, "skip the first record (header)")
	var err error
	var data []string
	flag.Parse()
	r := csv.NewReader(os.Stdin)
	w := csv.NewWriter(os.Stdout)
	r.LazyQuotes = true
	r.TrailingComma = true
	fields := getf(flag.Args())

	// loop over the input, making output
	for n :=0; ; n++ {
		data, err = r.Read()
		if err == io.EOF {
			break
		}
		if n == 0 && *headskip {
			continue
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		if len(fields) > 0 { // output selected fields
			selection := []string{}
			for _, n := range fields {
				if n >= 0 && n < len(data) {
					selection = append(selection, data[n])
				}
			}
			output(selection, w, *plainout)

		} else { // or output all fields
			output(data, w, *plainout)
		}
	}
	w.Flush()
}
