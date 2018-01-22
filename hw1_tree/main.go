package main

import (
	"fmt"
	"io"
	"os"
	"io/ioutil"
)

const (
	folderFmt        = "%s├───%s\n"
	folderFmtLast    = "%s└───%s\n"
	fileFmt          = "%s├───%s (%db)\n"
	fileFmtLast      = "%s└───%s (%db)\n"
	fileFmtEmpty     = "%s├───%s (empty)\n"
	fileFmtLastEmpty = "%s└───%s (empty)\n"
	updateStr        = "│\t"
	lastUpdateStr    = "\t"
)

type printFuncType func(formt Formatter, finfo os.FileInfo)

var Printer printFuncType

type Display struct {
	Path   string
	Indent string
}

type Formatter struct {
	last   bool
	output io.Writer
	disp   Display
}

func Filter(filt []os.FileInfo) []os.FileInfo {
	var filtered []os.FileInfo
	for _, f := range filt {
		if f.IsDir() {
			filtered = append(filtered, f)
		}
	}
	return filtered
}

func PrintFolder(formt Formatter, finfo os.FileInfo) {

	var display *Display
	display = &formt.disp

	var frmtStr, updStr string

	if formt.last {
		frmtStr, updStr = folderFmtLast, lastUpdateStr
	} else {
		frmtStr, updStr = folderFmt, updateStr
	}

	if finfo.IsDir() {
		fmt.Fprintf(formt.output, frmtStr, display.Indent, finfo.Name())
		updt := *display
		updt.update(Display{finfo.Name(), updStr})
		Walk(updt, formt.output, false)
	}
}

func PrintFIle(formt Formatter, finfo os.FileInfo) {
	var display *Display
	display = &formt.disp

	var frmtStr, updStr string

	if formt.last {
		frmtStr, updStr = folderFmtLast, lastUpdateStr
	} else {
		frmtStr, updStr = folderFmt, updateStr
	}

	if finfo.IsDir() {
		fmt.Fprintf(formt.output, frmtStr, display.Indent, finfo.Name())
		updt := *display
		updt.update(Display{finfo.Name(), updStr})
		Walk(updt, formt.output, true)
	} else {
		if formt.last {
			if finfo.Size() == 0 {
				frmtStr, updStr = fileFmtLastEmpty, lastUpdateStr
			} else {
				frmtStr, updStr = fileFmtLast, lastUpdateStr
			}
		} else {
			if finfo.Size() == 0 {
				frmtStr, updStr = fileFmtEmpty, updateStr
			} else {
				frmtStr, updStr = fileFmt, updateStr
			}
		}
		if finfo.Size() == 0 {
			fmt.Fprintf(formt.output, frmtStr, display.Indent, finfo.Name())
		} else {
			fmt.Fprintf(formt.output, frmtStr, display.Indent, finfo.Name(), finfo.Size())
		}
	}
}

func (d *Display) update(updt Display) {
	d.Path += `/`
	d.Path += updt.Path
	d.Indent += updt.Indent
}

func Walk(display Display, out io.Writer, f bool) (err error) {
	entries, err := ioutil.ReadDir(display.Path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du1: %s %v\n", display.Path, err)
		return err
	}

	if f == false {
		entries = Filter(entries)
	}

	if len(entries) > 0 {
		for i := 0; i < len(entries)-1; i++ {
			Printer(Formatter{false, out, display}, entries[i])
		}
		i := len(entries) - 1
		Printer(Formatter{true, out, display}, entries[i])
	}
	return

}

func dirTree(out io.Writer, path string, printFiles bool) (err error) {
	if printFiles {
		Printer = PrintFIle
	} else {
		Printer = PrintFolder
	}
	err = Walk(Display{path, ""}, out, printFiles)
	return
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
