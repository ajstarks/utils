// hsv2rgb -- convert hsv to rgb colors
package main

import (
	"flag"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func main() {
	var hue, saturation, value float64
	var name string
	flag.Float64Var(&hue, "h", 360, "hue")
	flag.Float64Var(&saturation, "s", 50, "saturation")
	flag.Float64Var(&value, "v", 50, "value")
	flag.StringVar(&name, "n", "", "named color like 'hsv(0,20,50)'")
	flag.Parse()

	var r, g, b int
	if len(name) > 0 {
		v := colorNumbers(name)
		if len(v) == 3 {
			hue, _ = strconv.ParseFloat(v[0], 64)
			saturation, _ = strconv.ParseFloat(v[1], 64)
			value, _ = strconv.ParseFloat(v[2], 64)
			r, g, b = hsv2rgb(hue, saturation, value)
		}
	} else {
		r, g, b = hsv2rgb(hue, saturation, value)
	}
	fmt.Printf("hsv(%g, %g, %g) => rgb(%d, %d, %d)\n", hue, saturation, value, r, g, b)
}

// colorNumbers returns a list of numbers from a comma separated list,
// in the form of xxx(n1, n2, n3), after removing tabs and spaces.
func colorNumbers(s string) []string {
	return strings.Split(strings.NewReplacer(" ", "", "\t", "").Replace(s[4:len(s)-1]), ",")
}

// hsv2rgb converts hsv(h (0-360), s (0-100), v (0-100)) to rgb
// reference: https://en.wikipedia.org/wiki/HSL_and_HSV#HSV_to_RGB
func hsv2rgb(h, s, v float64) (int, int, int) {
	s /= 100
	v /= 100
	if s > 1 || v > 1 {
		return 0, 0, 0
	}
	h = math.Mod(h, 360)
	c := v * s
	section := h / 60
	x := c * (1 - math.Abs(math.Mod(section, 2)-1))

	var r, g, b float64
	switch {
	case section >= 0 && section <= 1:
		r = c
		g = x
		b = 0
	case section > 1 && section <= 2:
		r = x
		g = c
		b = 0
	case section > 2 && section <= 3:
		r = 0
		g = c
		b = x
	case section > 3 && section <= 4:
		r = 0
		g = x
		b = c
	case section > 4 && section <= 5:
		r = x
		g = 0
		b = c
	case section > 5 && section <= 6:
		r = c
		g = 0
		b = x
	default:
		return 0, 0, 0
	}
	m := v - c
	r += m
	g += m
	b += m
	return int(r * 255), int(g * 255), int(b * 255)
}
