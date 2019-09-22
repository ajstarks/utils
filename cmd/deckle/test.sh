#!/bin/sh
(
	printf "<deck><slide>"
	deckle -raw=t -x 10 -y 10 -h -1.5 -w 50   -type lh -n 20 -color red
	deckle -raw=t -x 10 -y 10 -h 50   -w -1.5 -type lv -n 20 -color blue
	deckle -raw=t -x 10 -y 60 -h 1.5  -w 50   -type lh -n 20 -color red
	deckle -raw=t -x 60 -y 10 -h 50   -w 1.5  -type lv -n 20 -color blue
	printf "</slide>"
	printf "<slide>"
	deckle -raw=t -x 10 -y 10 -h -1.5 -w 50   -type fh -n 20 -color red
	deckle -raw=t -x 10 -y 10 -h 50   -w -1.5 -type fv -n 20 -color blue
	deckle -raw=t -x 10 -y 60 -h 1.5  -w 50   -type fh -n 20 -color red
	deckle -raw=t -x 60 -y 10 -h 50   -w 1.5  -type fv -n 20 -color blue
	printf "</slide>"
	printf "</deck>"
) | pdfdeck $* -stdout - > r.pdf

(
	printf "deck\n"
	printf "slide\n"
	deckle -raw=f -x 10 -y 10 -h -1.5 -w 50   -type lh -n 20 -color red
	deckle -raw=f -x 10 -y 10 -h 50   -w -1.5 -type lv -n 20 -color blue
	deckle -raw=f -x 10 -y 60 -h 1.5  -w 50   -type lh -n 20 -color red
	deckle -raw=f -x 60 -y 10 -h 50   -w 1.5  -type lv -n 20 -color blue
	printf "eslide\n"
	printf "slide\n"
	deckle -raw=f -x 10 -y 10 -h -1.5 -w 50   -type fh -n 20 -color red
	deckle -raw=f -x 10 -y 10 -h 50   -w -1.5 -type fv -n 20 -color blue
	deckle -raw=f -x 10 -y 60 -h 1.5  -w 50   -type fh -n 20 -color red
	deckle -raw=f -x 60 -y 10 -h 50   -w 1.5  -type fv -n 20 -color blue
	printf "eslide\n"
	printf "edeck\n"
) | decksh | pdfdeck $* -stdout - > d.pdf