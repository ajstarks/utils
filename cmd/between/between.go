// between: compute the hours between two dates
package main

import (
	"fmt"
	"os"
	"time"
)

const isofmt = "2006-01-02"
const usage = "specify start time and end time in the YYYY-MM-DD format; for example:\n\nbetween 1970-01-01 2009-10-11"

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}
	begintime := os.Args[1]
	endtime := os.Args[2]

	t0, err := time.Parse(isofmt, begintime)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s is not a valid time\n", begintime)
		os.Exit(2)
	}
	t1, err := time.Parse(isofmt, endtime)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s is not a valid time\n", endtime)
		os.Exit(3)
	}
	between := t1.Sub(t0).Hours()
	fmt.Printf("%s %s %g\n", begintime, endtime, between)
}
