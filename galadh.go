package galadh

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type glyphSet struct {
	pipe string
	last string
	item string
}

var (
	asciiGlyphSet   = glyphSet{pipe: "|", last: "`--", item: "|--"}
	unicodeGlyphSet = glyphSet{pipe: "│", last: "└──", item: "├──"}
)

// Galadh is simple clone of the posix tree command.
type Galadh struct {
	// Listing options
	allFiles       bool
	dirOnly        bool
	excludePattern string
	ignoreCase     bool
	includePattern string
	level          int
	matchDirs      bool
	printFullPath  bool
	// File optionsShowTree
	humanReadable       bool
	printSize           bool
	replaceNonPrintable bool
	outputFile          bool
	// Graphic options
	useASCII  bool
	useColors bool
	// Reporting
	cntFiles int
	cntDirs  int

	printer treePrinter
	w       io.Writer
}

// PrintTree prints a tree according the configured options for
// the passed directory.
func (g Galadh) PrintTree(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("Galadh: path '%v' does not exist", path)
		}

		if os.IsPermission(err) {
			return fmt.Errorf("Galadh: access denied for path '%v'", path)
		}

		return fmt.Errorf("Galadh: unknown error occured: '%w'", err)
	}

	if !fi.IsDir() {
		return fmt.Errorf("Galadh: '%v' is not a directory", path)
	}

	g.cntDirs, g.cntFiles = 0, 0

	fmt.Fprintln(g.w, path)
	g.printDir(path, g.level)
	return nil
}

func (g *Galadh) printDir(path string, lvl int) {

	_ := g.getFiles(fi)
}

func (g Galadh) getFiles(path string) (files []os.FileInfo, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		return []os.FileInfo{}, nil
	}

	fs, err := file.Readdir(0)
	if err != nil {
		return nil, err
	}

	// Todo: implement
	for _, f := range fs {
		fullPath := filepath.Join(path, f.Name())
		if !g.allFiles && isHidden()

	}

	return files, nil
}

type treePrinter struct {
	glyphs   glyphSet
	indents  []string
	lastItem bool
	w        io.Writer
}

func (t *treePrinter) indent() {
	if t.lastItem {
		t.indents = append(t.indents, strings.Repeat(" ", 5))
		return
	}

	t.indents = append(t.indents, t.glyphs.pipe+strings.Repeat(" ", 4))
}

func (t *treePrinter) unindent() {
	if len(t.indents) > 0 {
		t.indents = t.indents[:len(t.indents)-1]
	}
}

func (t *treePrinter) printItem(label string, last bool) {
	for i := range t.indents {
		fmt.Fprint(t.w, t.indents[i])
	}

	if last {
		fmt.Fprint(t.w, t.glyphs.last)
		t.lastItem = true
	} else {
		fmt.Fprint(t.w, t.glyphs.item)
		t.lastItem = false
	}

	fmt.Fprintln(t.w, " ", label)
}

func PrintTest() {
	printer := treePrinter{glyphs: unicodeGlyphSet, w: os.Stdout}
	printer.printItem("Hello", false)
	printer.printItem("World", false)
	printer.indent()
	printer.printItem("LotR", false)
	printer.indent()
	printer.printItem("Legolas", false)
	printer.printItem("Gandalf", false)
	printer.printItem("Frodo", true)
	printer.unindent()
	printer.printItem("Swtor", true)
	printer.indent()
	printer.printItem("Han", false)
	printer.printItem("Solo", false)
	printer.printItem("Leia", true)
	printer.unindent()
	printer.unindent()
	printer.printItem("Bye", true)
}
