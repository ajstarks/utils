// c19chart -- chart covid-19 data
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ajstarks/dchart2"
	"github.com/ajstarks/deck/generate"
)

const (
	c19URL      = "https://coronavirus.projectpage.app/.json?period=0"
	titlefmt    = "COVID-19 Global Status: "
	fatalfmt    = "Fatality Rate: "
	c19Filename = "c19.csv"
)

type yrange struct {
	min, max, step float64
}

// C19 is the json data returned by https://coronavirus.projectpage.app/.json?period=0
type C19 struct {
	Dates            []string `json:"dates"`
	Deaths           []int    `json:"deaths"`
	Confirmed        []int    `json:"confirmed"`
	AllTimeDeaths    int      `json:"alltimeDeaths"`
	AllTimeConfirmed int      `json:"allTimeConfirmed"`
}

// makedata reads from the API, makes the CSV
func makedata() error {
	if t := fileage(); t < (8 * time.Hour) {
		fmt.Fprintf(os.Stderr, "using the data file that is %v old\n", t)
		return nil
	}
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(c19URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response: %v from %s", resp.Status, c19URL)
	}
	var data C19
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}
	w, err := os.Create(c19Filename)
	if err != nil {
		return err
	}
	fmt.Fprintln(w, "date,deaths,confirmed")
	for i := 0; i < len(data.Dates); i++ {
		fmt.Fprintf(w, "\"%s\",%d,%d\n", data.Dates[i], data.Deaths[i], data.Confirmed[i])
	}
	return w.Close()
}

// readChartData reads chart data into chartboxes
func readChartData(s string) (dchart2.ChartBox, dchart2.ChartBox, error) {
	var cc dchart2.ChartBox
	var dc dchart2.ChartBox

	r, err := os.Open(s)
	if err != nil {
		return cc, dc, err
	}
	cc, err = dchart2.ReadCSV(r, "date,confirmed")
	if err != nil {
		return cc, dc, err
	}
	r.Seek(0, 0)
	dc, err = dchart2.ReadCSV(r, "date,deaths")
	if err != nil {
		return cc, dc, err
	}
	return cc, dc, r.Close()
}

// ftoa converts a floating point value to string
func ftoa(x float64, n int) string {
	return strconv.FormatFloat(x, 'f', n, 64)
}

// https://stackoverflow.com/questions/13020308/how-to-fmt-printf-an-integer-with-thousands-comma
func thousands(n float64, sep rune) string {
	s := ftoa(n, 0)
	startOffset := 0
	var b bytes.Buffer

	if n < 0 {
		startOffset = 1
		b.WriteByte('-')
	}
	l := len(s)
	commaIndex := 3 - ((l - startOffset) % 3)
	if commaIndex == 3 {
		commaIndex = 0
	}
	for i := startOffset; i < l; i++ {
		if commaIndex == 3 {
			b.WriteRune(sep)
			commaIndex = 0
		}
		commaIndex++
		b.WriteByte(s[i])
	}
	return b.String()
}

// c19curve shows the covid-19 curve
func c19curve(deck *generate.Deck, chart dchart2.ChartBox, label, color string, yr yrange, h float64) {
	left := chart.Left + 5
	ly := chart.Top - 10
	chart.Bottom = chart.Top - h
	dl := len(chart.Data)
	v := chart.Data[dl-1].Value
	pv := chart.Data[dl-2].Value

	pctchange := ((v - pv) / pv) * 100
	deck.Text(left, ly+5, label, "sans", 2.5, chart.DataColor)
	deck.Text(left, ly, thousands(v, ','), "sans", 4, color)
	deck.Text(left, ly-4, ftoa(pctchange, 3)+"% change", "sans", 2, chart.LabelColor)
	chart.DataColor = color
	chart.Frame(deck, 5)
	chart.XLabel(deck, 5)
	chart.DataFormat = "%0.f"
	chart.YAxis(deck, yr.min, yr.max, yr.step, false)
	chart.Line(deck, 0.2)
	chart.Opacity = 40
	chart.Area(deck)
}

// summarychart combines cases and deaths charts
func summarychart(deck *generate.Deck, cc, dc dchart2.ChartBox, casecolor, deathcolor string, yr yrange, h float64) {
	cc.Bottom = cc.Top - h
	cc.DataColor = casecolor
	cc.XLabel(deck, 5)
	cc.DataFormat = "%0.f"
	cc.YAxis(deck, yr.min, yr.max, yr.step, false)
	cc.Line(deck, 0.2)
	cc.Frame(deck, 5)
	cc.Opacity = 40
	cc.Area(deck)
	dc.Top = cc.Top
	dc.Maxvalue = cc.Maxvalue
	dc.DataColor = deathcolor
	dc.Bottom = cc.Bottom
	dc.Line(deck, 0.2)
	dc.Opacity = 40
	dc.Area(deck)
	deck.Text(cc.Right-40, 15, "Cases", "sans", 2, casecolor)
	deck.Text(cc.Right-10, 10, "Deaths", "sans", 2, deathcolor)
}

// labels makes chart labels
func labels(deck *generate.Deck, cc, dc dchart2.ChartBox, y float64, color string) {
	last := len(cc.Data) - 1
	frate := (dc.Data[last].Value / cc.Data[last].Value) * 100
	deck.Text(cc.Left, y, titlefmt+cc.Data[last].Label, "sans", 3.5, "")
	deck.TextEnd(cc.Right, y, fatalfmt+ftoa(frate, 2)+"%", "sans", 2, color)
}

func yrangeparse(s string) yrange {
	var min, max, step float64
	var yr yrange
	n, err := fmt.Sscanf(s, "%v,%v,%v", &min, &max, &step)
	if err != nil || n != 3 {
		return yr
	}
	yr.min = min
	yr.max = max
	yr.step = step
	return yr
}

func fileage() time.Duration {
	f, err := os.Stat(c19Filename)
	if err != nil {
		return 8 * time.Hour
	}
	return time.Since(f.ModTime())
}

func main() {
	var cyrs, dyrs string
	flag.StringVar(&cyrs, "cyr", "0,3.2e6,5e5", "case y range")
	flag.StringVar(&dyrs, "dyr", "0,2.2e5,5e4", "death y range")
	flag.Parse()

	ty := 92.0
	h := 20.0
	casecolor := "rgb(100,100,100)"
	deathcolor := "maroon"
	cyr := yrangeparse(cyrs)
	dyr := yrangeparse(dyrs)

	err := makedata()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	casesChart, deathsChart, err := readChartData(c19Filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}

	deck := generate.NewSlides(os.Stdout, 0, 0)
	deck.StartDeck()
	deck.StartSlide()
	labels(deck, casesChart, deathsChart, ty, deathcolor)
	casesChart.Top = 85
	c19curve(deck, casesChart, "Cases", casecolor, cyr, h)
	deathsChart.Top = 55
	c19curve(deck, deathsChart, "Deaths", deathcolor, dyr, h)
	casesChart.Top = 25
	summarychart(deck, casesChart, deathsChart, casecolor, deathcolor, cyr, h)
	deck.EndSlide()
	deck.EndDeck()
}
