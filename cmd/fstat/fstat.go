package main

import (
	"fmt"
	"os"
)

func status(filename string) (os.FileInfo, bool, error) {
	f, err := os.Stat(filename)
	if err != nil {
		return nil, false, err
	}
	return f, f.IsDir(), err
}

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
		printstat(n)
	}
}

func printstat(f os.FileInfo) {
	fmt.Printf(
		"%s\t%s\t%10d\t%s\n",
		f.Mode(),
		f.ModTime().Format("2006-01-02 15:04:05"),
		f.Size(),
		f.Name())
}

func main() {
	if len(os.Args) <= 1 {
		dirstat(".")
	} else {
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
}
