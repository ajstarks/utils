// spl -- image catalogs
package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

const (
	bslide   = "<slide bg=\"%s\">\n"
	eslide   = "</slide>"
	imfmt    = "<image xp=\"%.3f\" yp=\"%.3f\" width=\"%d\" height=\"%d\" scale=\"%.2g\" name=\"%s\"/>\n"
	namefmt  = "<text xp=\"%.3f\" yp=\"%.3f\" sp=\"2\" align=\"center\">%s</text>\n"
	snamefmt = "<text xp=\"50\" yp=\"5\" sp=\"3\" align=\"center\">%s</text>\n"
	simgfmt  = "<image xp=\"%.3f\" yp=\"%.3f\" width=\"%d\" height=\"%d\" name=\"%s\"/>\n"
	sdeck    = "<deck>"
	edeck    = "</deck>"
)

type Picture struct {
	x, y          float64
	width, height int
	name          string
	orientation   string
}

type Pictures []Picture

type Canvas struct {
	width, height            int
	left, right, top, bottom float64
	bgcolor                  string
	showname                 bool
}

func marginw(c Canvas, p Picture) (int, int) {
	aspect := float64(p.height) / float64(p.width)
	pw := float64(c.width) * (100.0 - (c.left + c.right)) / 100.0
	if int(pw) > p.width {
		return p.width, p.height
	}
	return int(pw), int(aspect * pw)
}

func marginh(c Canvas, p Picture) (int, int) {
	aspect := float64(p.height) / float64(p.width)
	ph := float64(c.height) * (100.0 - (c.top + c.bottom)) / 100.0
	if int(ph) > p.height {
		return p.width, p.height
	}
	return int(ph / aspect), int(ph)
}

func landlayout(c Canvas, p []Picture) {
	switch len(p) {
	case 1:
		p[0].x, p[0].y = 50, 50
		placepicsh(c, p[0:1], 90)
	case 2:
		p[0].x, p[1].x = 25, 75
		p[0].y, p[1].y = 50, 50
		placepicsh(c, p[0:2], 45)
	case 3:
		p[0].x, p[1].x, p[2].x = 17, 50, 83
		p[0].y, p[1].y, p[2].y = 50, 50, 50
		placepicsh(c, p[0:3], 30)
	}
}

func portlayout(c Canvas, p []Picture) {
	switch len(p) {
	case 1:
		p[0].x, p[0].y = 50, 50
		placepicsw(c, p[0:1], 90)
	case 2:
		p[0].x, p[1].x = 25, 75
		p[0].y, p[1].y = 50, 50
		placepicsw(c, p[0:2], 45)
	case 3:
		p[0].x, p[1].x, p[2].x = 17, 50, 83
		p[0].y, p[1].y, p[2].y = 50, 50, 50
		placepicsw(c, p[0:3], 30)
	}
}

func placepicsh(c Canvas, pics []Picture, targetpct float64) {
	fmt.Printf(bslide, c.bgcolor)
	tp := (targetpct / 100) * float64(c.height)
	for _, p := range pics {
		scalepct := (tp / float64(p.height)) * 100
		fmt.Printf(imfmt, p.x, p.y, p.width, p.height, scalepct, p.name)
		if c.showname {
			fmt.Printf(namefmt, p.x, 5.0, p.name)
		}
	}
	fmt.Println(eslide)
}

func placepicsw(c Canvas, pics []Picture, targetpct float64) {
	fmt.Printf(bslide, c.bgcolor)
	tp := (targetpct / 100) * float64(c.width)
	for _, p := range pics {
		scalepct := (tp / float64(p.width)) * 100
		fmt.Printf(imfmt, p.x, p.y, p.width, p.height, scalepct, p.name)
		if c.showname {
			fmt.Printf(namefmt, p.x, 5.0, p.name)
		}
	}
	fmt.Println(eslide)
}

func piclist(filelist []string) []Picture {
	var pic Picture
	p := []Picture{}
	for _, imagefile := range filelist {
		imf, err := os.Open(imagefile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		im, _, err := image.DecodeConfig(imf)
		if err != nil {
			imf.Close()
			continue
		}
		pic.width = im.Width
		pic.height = im.Height
		pic.name = imagefile
		p = append(p, pic)
		imf.Close()
	}
	return p
}

func ll(c Canvas, pics []Picture, n int) {
	lands := []Picture{}
	e := []Picture{}

	nl := 0
	for _, p := range pics {
		if p.width > p.height {
			nl++
			lands = append(lands, p)
			if nl%n == 0 {
				landlayout(c, lands)
				lands = e
			}
		}
	}
}

func lp(c Canvas, pics []Picture, n int) {
	ports := []Picture{}
	e := []Picture{}

	np := 0
	for _, p := range pics {
		if p.width < p.height {
			np++
			ports = append(ports, p)
			if np%n == 0 {
				portlayout(c, ports)
				ports = e
			}
		}
	}
}

func single(c Canvas, pics []Picture) {
	for i := 0; i < len(pics); i++ {
		if pics[i].width >= pics[i].height {
			landlayout(c, pics[i:i+1])
		} else {
			portlayout(c, pics[i:i+1])
		}
	}
}

func msingle(c Canvas, pics []Picture) {

	var pw, ph int
	for _, p := range pics {
		p.x, p.y = 50, 50
		if p.width > p.height {
			pw, ph = marginw(c, p)
		} else {
			pw, ph = marginh(c, p)
		}
		fmt.Printf(bslide, c.bgcolor)
		if c.showname {
			fmt.Printf(snamefmt, p.name)
		}
		fmt.Printf(simgfmt, p.x, p.y, pw, ph, p.name)
		fmt.Printf(eslide)
	}
}

func main() {
	cw := flag.Int("w", 1280, "canvas width")
	ch := flag.Int("h", 720, "canvas height")
	tm := flag.Float64("top", 5, "top margin")
	bm := flag.Float64("bottom", 5, "bottom margin")
	lm := flag.Float64("left", 5, "left margin")
	rm := flag.Float64("right", 5, "right margin")
	port := flag.Int("p", 0, "portrait n")
	land := flag.Int("l", 0, "landscape n")
	all := flag.Int("a", 0, "all n")
	showname := flag.Bool("showname", false, "show name")
	bgcolor := flag.String("bg", "white", "background color")
	flag.Parse()

	pics := piclist(flag.Args())
	c := Canvas{width: *cw, height: *ch, left: *lm, right: *rm, top: *tm, bottom: *bm, bgcolor: *bgcolor, showname: *showname}
	fmt.Println(sdeck)
	// fmt.Printf("<canvas width=\"%d\" height=\"%d\"/>\n", c.width, c.height)
	switch {
	case *port > 0:
		lp(c, pics, *port)
	case *land > 0:
		ll(c, pics, *land)
	case *all > 0:
		ll(c, pics, *all)
		lp(c, pics, *all)
	default:
		msingle(c, pics)
	}
	fmt.Println(edeck)
}
