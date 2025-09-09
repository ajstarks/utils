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
	flag.StringVar(&begintime, "begin", "", "begin time")
	flag.StringVar(&endtime, "end", "", "end time")
	flag.StringVar(&unit, "unit", "hour", "time unit (month, hour, minute, second, ms)")
	flag.Parse()

	if begintime == "" || endtime == "" {
		fmt.Fprintf(os.Stderr, "usage: between -begin YYYY-MM-DD -end YYYY-MM-DD -unit (hour, minute, or second\n")
		os.Exit(1)
	}
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
	case "ms":
		between = float64(t1.Sub(t0).Milliseconds())
	default:
		fmt.Fprintf(os.Stderr, "%s is not a valid time unit (use one of hr, min, sec, ms)\n", unit)
		os.Exit(4)
	}
	fmt.Printf("%s %s %.2f %s\n", begintime, endtime, between, unit)
}
