/*
ctime - time a command

Run the specified command, show the execution time in seconds to standard output.
Any output is ignored.

The -t option shows the specified string before the time display,
otherwise the first argument of the command is shown.

$ ctime sleep 10
sleep	10.00685

$ ctime -t "Go to sleep" sleep 5
Go to sleep 5.00442
*/

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func work(tag string, s []string) {
	if len(s) < 1 {
		return
	}
	if tag == "" {
		tag = s[0]
	}
	b := time.Now()
	err := exec.Command(s[0], s[1:]...).Run()
	e := time.Now()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	fmt.Printf("%v\t%.5f\n", tag, e.Sub(b).Seconds())
}

func main() {
	var tag = flag.String("t", "", "tag for the command")
	flag.Parse()
	work(*tag, flag.Args())	
}
