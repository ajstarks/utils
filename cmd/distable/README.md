# distable: make a distance table

![distable](morris.png)

Derived from the design seen in the atlas of Morris county, New Jersey from 1868,
```distable``` makes a distance table using deck markup.

## Data format

distable takes as input a file with this format:

	Place1
	<tab>place2:distance
	<tab>place3:distance

For example:
```
Beavertown
Boonton
	Beavertown:5.30
Budds Lake
	Beavertown:21.30
	Boonton:16.00
Chatham
	Beavertown:12.20
	Boonton:10.40
	Budds Lake:18.70
Chester
	Beavertown:21.40
	Boonton:16.30
	Budds Lake:5.40
	Chatham:15.10
```

## Command options

Read from named files or standard input.
```

distable [options] file...

Options:
  -left float
    	left margin (default 1)
  -size float
    	text size (default 1.1)
  -dsize float
  	  distance text size (default 0.65*size)
  -subtitle string
    	subtitle (default "distance in miles")
  -title string
    	chart title (default "Distances")
  -top float
    	top (default 90)

```


## Running

	$ distable.go  morris.d | pdfdeck -stdout - > morris.pdf

