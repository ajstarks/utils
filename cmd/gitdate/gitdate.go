// gitdate: visualize git commit history
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

const (
	gitime  = "Mon Jan 2 15:04:05 2006 -0700"
	isotime = "2006-01-02T15:04:05-07:00"
)

type config struct {
	title, btime, etime, color           string
	left, right, radius, ypoint, opacity float64
}

// seconds returns the number of seconds of the specified time
// since the Unix epoch (Jan 1, 1970, 00:00:00 UTC)
func seconds(s string) int64 {
	t, err := time.Parse(gitime, s)
	if err != nil {
		return -1
	}
	return t.Unix()
}

//vmap maps one interval to another
func vmap(value float64, low1 float64, high1 float64, low2 float64, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

// process reads a series of line containing timestamps
// in the ("Mon Jan 2 15:04:05 2006 -0700") format
// and maps each time time to a linear scale.
func process(w io.Writer, r io.Reader, c config) error {
	b, err := time.Parse(isotime, c.btime)
	if err != nil {
		return err
	}
	e, err := time.Parse(isotime, c.etime)
	if err != nil {
		return err
	}
	beg := b.Unix()
	end := e.Unix()

	labely := c.ypoint + 5
	fmt.Fprintf(w, "ctext %q %v %v %v\n", c.btime, c.left, labely, 1)
	fmt.Fprintf(w, "ctext %q %v %v %v\n", c.etime, c.right, labely, 1)
	fmt.Fprintf(w, "ctext %q %v %v 2\n", c.title, c.left+((c.right-c.left)/2), labely)
	fmt.Fprintf(w, "vline %v %v %v 0.1\n", c.left, c.ypoint, 4)
	fmt.Fprintf(w, "vline %v %v %v 0.1\n", c.right, c.ypoint, 4)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		t := scanner.Text()
		s := seconds(t)
		x := vmap(float64(s), float64(beg), float64(end), c.left, c.right)
		fmt.Fprintf(w, "circle %v %v %v %q %v\n", x, c.ypoint, c.radius, c.color, c.opacity)
	}
	return scanner.Err()
}

func main() {
	title := flag.String("title", "commit history", "title")
	btime := flag.String("begin", "1970-01-01T00:00:00+00:00", "begin time")
	etime := flag.String("end", time.Now().Format(isotime), "end time")
	ypoint := flag.Float64("y", 50, "y point")
	radius := flag.Float64("r", 2, "radius")
	color := flag.String("color", "black", "color")
	left := flag.Float64("left", 10, "left")
	right := flag.Float64("right", 90, "right")
	opacity := flag.Float64("opacity", 20, "opacity")
	flag.Parse()

	c := config{
		title:   *title,
		btime:   *btime,
		etime:   *etime,
		ypoint:  *ypoint,
		radius:  *radius,
		color:   *color,
		opacity: *opacity,
		left:    *left,
		right:   *right,
	}

	err := process(os.Stdout, os.Stdin, c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
