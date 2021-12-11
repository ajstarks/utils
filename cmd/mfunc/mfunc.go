// mfunc: math functions
package main

import (
	"flag"
	"fmt"
	"math"
)

type tfunc struct {
	label    string
	function func(x float64) float64
}

func main() {
	funcname := flag.String("f", "sine", "function name")
	min := flag.Float64("min", 0.0, "minimum")
	max := flag.Float64("max", math.Pi*2, "maximum")
	incr := flag.Float64("incr", 0.1, "increment")
	xfmt := flag.String("xfmt", "%.2f", "x format")
	yfmt := flag.String("yfmt", "%.4f", "y format")
	flag.Parse()
	var f tfunc
	switch *funcname {
	case "sine", "sin":
		f = tfunc{"y=sin(x)", math.Sin}
	case "cosine", "cos":
		f = tfunc{"y=cos(x)", math.Cos}
	case "sqrt":
		f = tfunc{"y=sqrt(x)", math.Sqrt}
	case "log":
		f = tfunc{"y=log(x)", math.Log}
	case "log10":
		f = tfunc{"y=log10(x)", math.Log10}
	case "log2":
		f = tfunc{"y=log2(x)", math.Log2}
	case "tan":
		f = tfunc{"y=tan(x)", math.Tan}
	case "exp":
		f = tfunc{"y=exp(x)", math.Tan}
	case "sincos":
		f = tfunc{"y=sin(x) * cos(x)",
			func(x float64) float64 { return math.Sin(x) * math.Cos(x) }}
	default:
		f = tfunc{"y=1", func(float64) float64 { return 1 }}
	}
	fmt.Printf("# %s\n", f.label)
	format := *xfmt + "\t" + *yfmt + "\n"
	for x := *min; x <= *max; x += *incr {
		fmt.Printf(format, x, f.function(x))
	}
}
