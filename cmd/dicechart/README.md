# dicechart

Make dicecharts in the style of [Monroe Work's "Negro Year Book"](https://nightingaledvs.com/monroe-nathan-work-education-in-the-negro-year-book/), using deck markup

![dicechart](dicechart.png)

## Usage

```dicechart [options] [file.csv]```

where file.csv contains "label",value pairs

for example:

```
"MARYLAND",57
"NORTH CAROLINA",50
"GEORGIA",48
"VIRGINIA",47
"TEXAS",47
"FLORIDA",43
"ALABAMA",27
"SOUTH CAROLINA",26
"LOUISIANA",23
```
If no file is specified, input is from the standard input.  Output is always to standard output.
Typical usage is to use pdfdeck to render the chart in pdf:

```
$ dicechart data.csv > chart.xml
$ pdfdeck chart.xml
```

or in a pipeline:

```
$ dicechart data.csv | pdfdeck -stdout - > chart.pdf
```

Command options are:
```

  -color string
    	dotcolor (default "black")
  -dotsize float
    	dot size (default 1)
  -ds float
    	dice spacing (default 5)
  -dw float
    	dice width (default 1.5)
  -dx float
    	data left position (default 35)
  -height float
    	canvas height (default 612)
  -lx float
    	label left position (default 10)
  -textsize float
    	canvas width (default 2)
  -title string
    	chart title
  -top float
    	top of the chart (default 85)
  -unit
      dice unit (default 5)
  -valsize float
    	canvas width (default 2)
  -vskip float
    	vertical skip (default 7)
  -width float
    	canvas width (default 792)
```
