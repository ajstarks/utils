package readpalette

import (
	"bufio"
	"fmt"
	"image/color"
	"io"
	"strings"
)

type spalette map[string][]string
type rgbpalette map[string][]color.NRGBA

func rgb(x uint32) (uint8, uint8, uint8) {
	r := x & 0xff0000 >> 16
	g := x & 0x00ff00 >> 8
	b := x & 0x0000ff
	return uint8(r), uint8(g), uint8(b)
}

func ReadString(r io.Reader) (spalette, error) {
	var p spalette
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		args := strings.Fields(scanner.Text())
		l := len(args)
		if l < 2 {
			continue
		}
		name := args[0]
		value := make([]string, l-1)
		copy(value, args[1:])
		p[name] = value
	}
	return p, scanner.Err()
}

func ReadRGB(r io.Reader) (rgbpalette, error) {
	palette, err := ReadString(r)
	if err != nil {
		return nil, err
	}

	var rp rgbpalette
	var x uint32
	for name, value := range palette {
		colors := make([]color.NRGBA, len(value))
		for i, c := range value {
			fmt.Sscanf(c[1:], "%x", &x)
			r, g, b := rgb(x)
			colors[i] = color.NRGBA{R: r, G: g, B: b, A: 0xff}
		}
		rp[name] = colors
	}
	return rp, nil
}
