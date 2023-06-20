package main

import (
	"fmt"
	"math"
)

func circle(x, y, size float64, color string) {
	fmt.Printf("circle %.2f %.2f %.2f %q\n", x, y, size, color)
}

func cpolar(x, y, r, t, size float64, color string) {
	fmt.Printf("p=polar %.2f %.2f %.2f %.2f\ncircle p_x  p_y %.2f %q\n", x, y, r, t, size, color)
}

func polar(x, y, r, deg float64) (float64, float64) {
	theta := deg * (math.Pi / 180)
	px := x + (r * math.Cos(theta))
	py := y + (r * math.Sin(theta))
	return px, py
}

func begin(color string) {
	fmt.Printf("deck\ncanvas 1000 1000\nslide \"%s\"\n", color)
}

func end() {
	fmt.Println("eslide\nedeck")
}

func planet(x, y, size, radius, a1, a2, steps float64, color string) {
	for t := a1; t < a2; t += steps {
		cpolar(x, y, radius, t, size, color)
	}
}

func solar(x, y, csize, radius, psize, steps float64, ccolor, pcolor string) {
	circle(x, y, csize, ccolor)
	planet(x, y, psize, radius, 0, 360, steps, pcolor)
}

func d1() {
	begin("black")
	circle(50, 50, 20, "red")
	step := 360.0 / 15
	for t := 0.0; t <= 360; t += step {
		px, py := polar(50, 50, 30, t)
		solar(px, py, 5, 5, 1, 10, "red", "orange")
	}
	for t := step / 2; t <= 360; t += step {
		px, py := polar(50, 50, 40, t)
		solar(px, py, 5, 5, 1, 10, "red", "orange")
	}
	end()
}

func hsv(hue, sat, value int) string {
	return fmt.Sprintf("hsv(%d,%d,%d)", hue, sat, value)
}

func d2a(r float64) {
	/*
		cstep := 1.0
		c := 1.0
		step := 20.0
		halfstep := step / 2
		csize := r * 1.5
		hue := 10
		circle(50, 50, csize, "hsv(0,100,100)")
		planet(50, 50, c, r, 0, 360, step, hsv(hue, 100, 100))
		for i := 0; i < 6; i++ {
			r += float64(i + 2)
			c += cstep
			hue += 7
			if i%2 == 0 {
				planet(50, 50, c, r, halfstep, 360, step, hsv(hue, 100, 100))
			} else {
				if i == 5 {
					r += 1
					hue = 10
				}
				planet(50, 50, c, r, 0, 360, step, hsv(hue, 100, 100))
			}
		}
	*/
}

func d2b(r float64) {
	cstep := 1.0
	c := 1.0
	step := 20.0
	halfstep := step / 2
	csize := r * 1.5
	hue := 10

	circle(50, 50, csize, "hsv(0,100,100)")
	planet(50, 50, c, r, 0, 360, step, "hsv(10,100,100)")
	r += 2
	c += cstep
	hue += 7
	planet(50, 50, c, r, halfstep, 360, step, hsv(hue, 100, 100))
	r += 3
	c += cstep
	hue += 7
	planet(50, 50, c, r, 0, 360, step, hsv(hue, 100, 100))
	r += 4
	c += cstep
	hue += 7
	planet(50, 50, c, r, halfstep, 360, step, hsv(hue, 100, 100))
	r += 5
	c += cstep
	hue += 7
	planet(50, 50, c, r, 0, 360, step, hsv(hue, 100, 100))
	r += 6
	c += cstep
	hue += 7
	planet(50, 50, c, r, halfstep, 360, step, hsv(hue, 100, 100))
	r += 8
	hue = 0
	planet(50, 50, 7, r, 0, 360, step, hsv(hue, 100, 100))
}

func d2() {
	begin("black")
	r := 10.0
	cstep := 1.0
	c := 1.0
	step := 20.0
	halfstep := step / 2
	csize := r * 1.5
	hue := 10

	circle(50, 50, csize, "hsv(0,100,100)")
	planet(50, 50, c, r, 0, 360, step, "hsv(10,100,100)")
	r += 2
	c += cstep
	hue += 7
	planet(50, 50, c, r, halfstep, 360, step, hsv(hue, 100, 100))
	r += 3
	c += cstep
	hue += 7
	planet(50, 50, c, r, 0, 360, step, hsv(hue, 100, 100))
	r += 4
	c += cstep
	hue += 7
	planet(50, 50, c, r, halfstep, 360, step, hsv(hue, 100, 100))
	r += 5
	c += cstep
	hue += 7
	planet(50, 50, c, r, 0, 360, step, hsv(hue, 100, 100))
	r += 6
	c += cstep
	hue += 7
	planet(50, 50, c, r, halfstep, 360, step, hsv(hue, 100, 100))
	r += 8
	hue = 0
	planet(50, 50, 7, r, 0, 360, step, hsv(hue, 100, 100))

	for t := 0.0; t <= 360; t += step {
		px, py := polar(50, 50, r, t)
		planet(px, py, 1, 5, 0, 360, 30, hsv(38, 100, 100))
	}
	end()
}

func main() {
	d2()
}
