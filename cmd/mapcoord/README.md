# mapcoord -- convert lat/long pairs to deck/decksh markup

mapcoord reads space separated lat/long pairs from stdin or specified files, and emits deck/decksh markup denoting the path to stdout.
Typically other programs will generate the input, for example ```fitscsvcoord``` reads CSV files with FITS data:

```
$ fitscsvcoord path.csv | mapcoord [options] > path.dsh
```

## Options

```
$ mapcoord --help
Usage of mapcoord:
  -bbox string
    	bounding box color ("" no box)
  -bgcolor string
    	background color
  -color string
    	line color (default "black")
  -fulldeck
    	make a full deck (default true)
  -info
    	only report bounding box, and center
  -latmax float
    	latitude x maxmum (default 90)
  -latmin float
    	latitude x minimum (default -90)
  -linewidth float
    	line width (default 0.1)
  -longmax float
    	longitude y maximum (default 180)
  -longmin float
    	longitude y minimum (default -180)
  -shape string
    	polygon, polyline (default "polyline")
  -style string
    	deck, decksh, plain (default "deck")
  -xmax float
    	canvas x maxmum (default 95)
  -xmin float
    	canvas x minimum (default 5)
  -ymax float
    	canvas y maximum (default 95)
  -ymin float
    	canvas y minimum (default 5)
```

