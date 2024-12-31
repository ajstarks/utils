// fstat -- file status
package main

import (
	"fmt"
	"os"
)

// status gets file status
func status(filename string) (os.FileInfo, bool, error) {
	f, err := os.Stat(filename)
	if err != nil {
		return nil, false, err
	}
	return f, f.IsDir(), err
}

// dirstat shows directory infomation
func dirstat(name string) {
	f, err := os.Open(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	fi, err := f.Readdir(-1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	for _, n := range fi {
		if n.Name()[0] != '.' {
			printstat(n)
		}
	}
}

// printstat shows file status
func printstat(f os.FileInfo) {
	fmt.Printf(
		"%-15s %20s %15d\t%s\n",
		f.Mode(),
		f.ModTime().Format("2006-01-02 15:04:05"),
		f.Size(),
		f.Name())
}

func main() {
	// if no args, show the current directory
	if len(os.Args) <= 1 {
		dirstat(".")
		return
	}
	// for every argument, print directory or file info
	for i, filename := range os.Args[1:] {
		s, isdir, err := status(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		if isdir {
			if i > 0 {
				fmt.Printf("%s:\n", filename)
			}
			dirstat(filename)
		} else {
			printstat(s)
		}
	}
}
