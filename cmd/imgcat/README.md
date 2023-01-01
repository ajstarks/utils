# imgcat -- make a multipage image catalog using deck markup

imgcat generates deck markup for an image catalog, given a list of supported images (PNG, JPEG, GIF).  The generated markup is usually fed to pdfdeck for rendering.

## usage

```
Usage of imgcat:
  -a int
    	all n
  -bg string
    	background color (default "white")
  -bottom float
    	bottom margin (default 5)
  -h int
    	canvas height (default 720)
  -l int
    	landscape n
  -left float
    	left margin (default 5)
  -p int
    	portrait n
  -right float
    	right margin (default 5)
  -showname
    	show name
  -top float
    	top margin (default 5)
  -w int
    	canvas width (default 1280)
```

## examples

make a 3up (3 image/page) catalog for both portrait and landscape images
```
$ imgcat -a 3 *.jpg | pdfdeck -stdout - > 3up-catalog.pdf 
```

make a 2 image/page catalog for portrait images

```
$ imgcat -p 2 *.jpg |  pdfdeck -stdout - > 2-up.pdf
```

make a 3 up catalog of landscape images; include the filenames

```
$ imgcat -showname -l 3 *.png |   pdfdeck -stdout - > catalog.pdf
```