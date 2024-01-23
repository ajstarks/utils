package main

import (
	"bufio"
	"fmt"
	"image/color"
	"io"
	"os"
	"strconv"
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
	scanner := bufio.NewScanner(r)
	p := make(spalette)
	for scanner.Scan() {
		args := strings.Fields(scanner.Text())
		l := len(args)
		if l < 2 {
			continue
		}
		name := args[0]
		p[name] = args[1:]
	}
	return p, scanner.Err()
}

func ReadRGB(r io.Reader) (rgbpalette, error) {

	palette, err := ReadString(r)
	if err != nil {
		return nil, err
	}

	rp := make(rgbpalette)
	for name, value := range palette {
		colors := make([]color.NRGBA, len(value))
		i := 0
		for _, c := range value {
			if len(c) != 7 {
				continue // must be #nnnnnn
			}
			x, err := strconv.ParseUint(c[1:], 16, 32)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				continue
			}
			r, g, b := rgb(uint32(x))
			colors[i] = color.NRGBA{R: r, G: g, B: b, A: 0xff}
			i++
		}
		rp[name] = colors
	}
	return rp, nil
}
