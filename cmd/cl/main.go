// cl: show a file with numbered lines, showing the whole file, a range, or a single line.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// linerange returns the begin and end using a "-" as delimiter
func linerange(s string) (int, int) {
	b := 0
	e := 0
	var err error
	l := strings.Split(s, "-")
	switch len(l) {
	case 1: // a single number (for example: "10")
		b, err = strconv.Atoi(l[0])
		if err != nil {
			return 0, 0
		}
	case 2: // two numbers (for example "10-20")
		b, err = strconv.Atoi(l[0])
		if err != nil {
			return 0, 0
		}
		e, err = strconv.Atoi(l[1])
		if err != nil {
			return b, 0
		}
	default:
		return 0, 0
	}
	return b, e
}

func countLines(r io.Reader, lr string) {
	begin, end := linerange(lr)
	scanner := bufio.NewScanner(r)
	for n := 1; scanner.Scan(); n++ {
		t := scanner.Text()

		// show all lines
		if begin == 0 && end == 0 {
			fmt.Printf("%d: %s\n", n, t)
			continue
		}
		// show a single specified line
		if n == begin && end == 0 {
			fmt.Printf("%d: %s\n", n, t)
			break
		}
		// show a range of lines
		if n >= begin && n <= end {
			fmt.Printf("%d: %s\n", n, t)
		}
	}
}

func main() {
	var linerange string
	flag.StringVar(&linerange, "n", "all", "line range (begin-end)")
	flag.Parse()
	countLines(os.Stdin, linerange)
}
