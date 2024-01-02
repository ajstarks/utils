# desordres

desordres make visuals in the style of Des Ordres by Vera Molnár, using decksh markup
For example:

make a 1000x1000 pdf using default parameters

```
desordres | decksh | pdfdeck -stdout -pagesize 1000x1000 - > gray.pdf
```
<img src="gray.png" width=500 height=500/>


make a 1000x1000 PNG file with random colors HSV(20-60, 100, 100), 14 tiles/row

```
desordres -tiles 14 -color '20:60'  -bgcolor=black | decksh > f.xml
pngdeck -pagesize 1000x1000 f.xml
```
<img src="hot-14-20-60.png" width=500 height=500/>

## options
```
Usage of desordres:
  -bgcolor string
    	background color (default "white")
  -color string
    	pen color (default "gray")
  -maxlw float
    	maximum line thickness (default 1)
  -tiles float
    	tiles/row (default 10)
```

