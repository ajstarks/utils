# geodeck -- convert lat/long pairs to deck/decksh markup

geodeck reads space separated decimal lat/long pairs from stdin or specified files, and emits deck/decksh markup representing the path to stdout.
Typically other programs will generate the input, for example the ```fitscsvcoord``` command [reads CSV files with FIT data](https://developer.garmin.com/fit/fitcsvtool/).
Note that ```fitscsvcoord``` converts from "semicircle" units to decimal latitude and longitude.

```
java -jar $JARLOC/FitCSVTool.jar -iso8601 path.fit |
grep position_lat | 
csvread  -plain=t 4 7 10  | 
awk 'length($1) == 20 {printf "%.6f %.6f\n",  $2 * (180/2^31) , $3 * (180/2^31)}'

```

```
$ fitscsvcoord path.csv > path.coord
$ geodeck [options] path.coord > path.dsh
```

The ```--info``` option reports information on the center and bounding box of the coordinates without deck generation.
The reported options may be used in subsequent calls to geodeck or used in other tools like [```create-static-map```](https://github.com/flopp/go-staticmaps/tree/master/create-static-map)

```
$ geodeck --info path.coord
--center=40.6291415,-74.4224255 -bbox="40.636468,-74.4292|40.621815,-74.415651" --longmin=-74.4292 --longmax=-74.415651 --latmin=40.621815 --latmax=40.636468
```

## Options

```
$ geodeck --help
Usage of geodeck:
  -bbox string
    	bounding box color ("" no box)
  -bgcolor string
    	background color
  -color string
    	line color (default "black")
  -fulldeck
    	make a full deck (default true)
  -info
    	only report bounding box, center, and extremes
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

