// between: compute the time (hours, minutes, or seconds) between two dates
package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	const isofmt = "2006-01-02"
	var begintime, endtime, unit string
	flag.StringVar(&begintime, "begin", "2006-01-02", "begin time")
	flag.StringVar(&endtime, "end", "2009-11-10", "end time")
	flag.StringVar(&unit, "unit", "hour", "time unit (month, hour, minute, second)")
	flag.Parse()

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
	var between float64
	switch unit {
	case "hr", "hour", "h":
		between = t1.Sub(t0).Hours()
	case "min", "minute", "m":
		between = t1.Sub(t0).Minutes()
	case "sec", "second", "s":
		between = t1.Sub(t0).Seconds()
	}

	fmt.Printf("%s %s %.2f %s\n", begintime, endtime, between, unit)
}
