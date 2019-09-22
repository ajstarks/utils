#!/bin/sh
(
	printf "<deck><slide>"
	deckle -x 10 -y 10 -h -1.5 -w 50   -type lh -n 20 -color red
	deckle -x 10 -y 10 -h 50   -w -1.5 -type lv -n 20 -color blue
	deckle -x 10 -y 60 -h 1.5  -w 50   -type lh -n 20 -color red
	deckle -x 60 -y 10 -h 50   -w 1.5  -type lv -n 20 -color blue
	printf "</slide></deck>"
) | pdf $*