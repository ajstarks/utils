// imgps -- show GPS coordinates contained in EXIF data
package main

import (
	"fmt"
	"os"

	"github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
)

// process reads GPS info from a filename
func process(filename string) error {
	rawExif, err := exif.SearchFileAndExtractExif(filename)
	if err != nil {
		return err
	}
	im, err := exifcommon.NewIfdMappingWithStandard()
	if err != nil {
		return err
	}
	ti := exif.NewTagIndex()
	_, index, err := exif.Collect(im, ti, rawExif)
	if err != nil {
		return err
	}
	ifd, err := index.RootIfd.ChildWithIfdPath(exifcommon.IfdGpsInfoStandardIfdIdentity)
	if err != nil {
		return err
	}
	gi, err := ifd.GpsInfo()
	if err != nil {
		return err
	}
	lat, lon := ddLatLon(gi)
	fmt.Printf("%s %.8f %.8f\n", filename, lat, lon)
	return nil
}

// ddLatLong returns lat/long in decimal degrees
func ddLatLon(g *exif.GpsInfo) (float64, float64) {
	glat := g.Latitude
	glon := g.Longitude

	ddlat := glat.Degrees + (glat.Minutes / 60) + (glat.Seconds / 3600)
	ddlon := glon.Degrees + (glon.Minutes / 60) + (glon.Seconds / 3600)

	if glat.Orientation == 'S' {
		ddlat = -ddlat
	}
	if glon.Orientation == 'W' {
		ddlon = -ddlon
	}
	return ddlat, ddlon
}

func main() {
	for _, f := range os.Args {
		err := process(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", f, err)
			continue
		}
	}
}
