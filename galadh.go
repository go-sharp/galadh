package galadh

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-sharp/color"
)

type glyphSet struct {
	pipe string
	last string
	item string
}

var (
	// Glyphsets
	asciiGlyphSet   = glyphSet{pipe: "|", last: "`--", item: "|--"}
	unicodeGlyphSet = glyphSet{pipe: "│", last: "└──", item: "├──"}

	// Color functions
	errColor = color.New(color.FgRed).SprintFunc()
	dirColor = color.New(color.FgBlue).SprintFunc()
	binColor = color.New(color.FgMagenta).SprintFunc()
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

	glyphs  glyphSet
	indents []string
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
	// Fetch content of directory
	// files, err = g.getFiles(path)
	// if err != nil {
	// 	g.printer.printItem()
	// }
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
	if err != nil && err != io.EOF {
		return nil, err
	}

	// If we use the match dir options, check
	// if folder name matches and set var.
	dirPatternMatch := false
	if g.matchDirs {
		dirPatternMatch = g.matchPattern(g.includePattern, fi.Name())
	}

	// Processing files of this directory
	for _, f := range fs {
		// If option all files is not set
		// and it is a hidden file, we skip the file.
		if !g.allFiles && isHidden(f) {
			continue
		}

		// If only dir option is set and we have
		// a file, we skip it.
		if g.dirOnly && !f.IsDir() {
			continue
		}

		// If an include pattern is set, check
		// if file matches the pattern and continue
		// if not. If dir pattern matches include all files.
		if !(dirPatternMatch || g.matchPattern(g.includePattern, f.Name())) {
			continue
		}

		// Skip file if exclude pattern matches
		if g.matchPattern(g.excludePattern, f.Name()) {
			continue
		}

		files = append(files, f)
	}

	return files, nil
}

func (g Galadh) matchPattern(pattern, name string) bool {
	if pattern == "" {
		return true
	}

	match, err := filepath.Match(pattern, name)
	if err != nil {
		return false
	}

	return match
}

func (g *Galadh) indent(lastItem bool) {
	if lastItem {
		g.indents = append(g.indents, strings.Repeat(" ", 5))
		return
	}

	g.indents = append(g.indents, g.glyphs.pipe+strings.Repeat(" ", 4))
}

func (g *Galadh) unindent() {
	if len(g.indents) > 0 {
		g.indents = g.indents[:len(g.indents)-1]
	}
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

func (t *treePrinter) printItem(typ itemTyp, label string) {
	for i := range t.indents {
		fmt.Fprint(t.w, t.indents[i])
	}

	if typ&lastItem == lastItem {
		fmt.Fprint(t.w, t.glyphs.last)
		t.lastItem = true
	} else {
		fmt.Fprint(t.w, t.glyphs.item)
		t.lastItem = false
	}

	if typ&errorItem == errorItem {

	}

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
