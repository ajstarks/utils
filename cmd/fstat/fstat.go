// fstat -- file status
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"sort"
)

const (
	dirfmt  = "%-15s %20s %15d\t%s\n"
	timefmt = "2006-01-02 15:04:05"
)

// flags for sorting and printing
type dirflags struct {
	na      bool // name ascending
	sa      bool // size ascending
	nd      bool // name descending
	sd      bool // size descending
	older   bool // date oldest first
	newer   bool // date newest first
	showdot bool // show dotfiles
}

// status gets file status
func status(filename string) (os.FileInfo, bool, error) {
	f, err := os.Stat(filename)
	if err != nil {
		return nil, false, err
	}
	return f, f.IsDir(), err
}

// nameSortAsc sorts directory entries by name (ascending)
func nameSortAsc(files []fs.DirEntry) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
}

// nameSortDec sorts directory entries by name (ascending)
func nameSortDec(files []fs.DirEntry) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() > files[j].Name()
	})
}

// sizeSortAsc sort directory entries by file sizes (ascendiing)
func sizeSortAsc(files []fs.DirEntry) {
	sort.Slice(files, func(i, j int) bool {
		fi, erri := files[i].Info()
		fj, errj := files[j].Info()
		if erri != nil || errj != nil {
			return false
		}
		ni := fi.Size()
		nj := fj.Size()
		return ni < nj
	})
}

// sizeSortDec sort directory entries by file sizes (ascendiing)
func sizeSortDec(files []fs.DirEntry) {
	sort.Slice(files, func(i, j int) bool {
		fi, erri := files[i].Info()
		fj, errj := files[j].Info()
		if erri != nil || errj != nil {
			return false
		}
		ni := fi.Size()
		nj := fj.Size()
		return ni > nj
	})
}

// timeSortOlder sorts directory entries by time older to newer
func timeSortOlder(files []fs.DirEntry) {
	sort.Slice(files, func(i, j int) bool {
		fi, erri := files[i].Info()
		fj, errj := files[j].Info()
		if erri != nil || errj != nil {
			return false
		}
		ti := fi.ModTime()
		tj := fj.ModTime()
		return ti.Before(tj)
	})
}

// timeSortNewer sorts directory entries by time newer to older
func timeSortNewer(files []fs.DirEntry) {
	sort.Slice(files, func(i, j int) bool {
		fi, erri := files[i].Info()
		fj, errj := files[j].Info()
		if erri != nil || errj != nil {
			return false
		}
		ti := fi.ModTime()
		tj := fj.ModTime()
		return ti.After(tj)
	})
}

// dirstat shows directory infomation
func dirstat(dirname string, sf dirflags) {
	f, err := os.Open(dirname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	di, err := f.ReadDir(0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	// set the sort option
	if sf.na {
		nameSortAsc(di)
	}
	if sf.nd {
		nameSortDec(di)
	}
	if sf.sa {
		sizeSortAsc(di)
	}
	if sf.sd {
		sizeSortDec(di)
	}
	if sf.older {
		timeSortOlder(di)
	}
	if sf.newer {
		timeSortNewer(di)
	}

	// print the entries
	for _, d := range di {
		fi, err := d.Info()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		if !sf.showdot && fi.Name()[0] == '.' {
			continue
		}
		printstat(fi)
	}
}

// printstat shows file status
func printstat(f os.FileInfo) {
	fmt.Printf(dirfmt, f.Mode(), f.ModTime().Format(timefmt), f.Size(), f.Name())
}

func main() {
	var df dirflags
	flag.BoolVar(&df.showdot, "a", false, "show dot files")
	flag.BoolVar(&df.na, "na", true, "sort by name ascending")
	flag.BoolVar(&df.nd, "nd", false, "sort by name descending")
	flag.BoolVar(&df.sa, "sa", false, "sort by size ascending")
	flag.BoolVar(&df.sd, "sd", false, "sort by size descending")
	flag.BoolVar(&df.older, "old", false, "sort by age oldest first")
	flag.BoolVar(&df.newer, "new", false, "sort by age newest first")
	flag.Parse()
	args := flag.Args()

	// if other options are set turn off the default
	if df.nd || df.sa || df.sd || df.older || df.newer {
		df.na = false
	}
	// if the default is set turn off other options
	if df.na {
		df.nd, df.sa, df.sd, df.older, df.newer = false, false, false, false, false
	}
	// make size and time sorting option are exclusive; either ascending OR descending
	if df.sa && df.sd {
		fmt.Fprintln(os.Stderr, "pick one: size ascending or descending")
		os.Exit(1)
	}
	if df.older && df.newer {
		fmt.Fprintln(os.Stderr, "pick one: old or new")
		os.Exit(1)
	}
	// if no args, show the current directory
	if len(args) == 0 {
		dirstat(".", df)
		return
	}
	// for every argument, print directory or file info
	for i, filename := range args {
		s, isdir, err := status(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		if isdir {
			if i > 0 {
				fmt.Printf("%s:\n", filename)
			}
			dirstat(filename, df)
		} else {
			printstat(s)
		}
	}
}
