package galadh

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/go-sharp/color"
)

type glyphSet struct {
	pipe string
	last string
	item string
}

const (
	kbSize = 1 << 10
	mbSize = 1 << 20
	gbSize = 1 << 30
	tbSize = 1 << 40
	pbSize = 1 << 50
)

var (
	// Glyphsets
	asciiGlyphSet   = glyphSet{pipe: "|", last: "`--", item: "|--"}
	unicodeGlyphSet = glyphSet{pipe: "│", last: "└──", item: "├──"}

	// Color functions
	errColor    = color.New(color.FgRed).SprintFunc()
	dirColor    = color.New(color.FgBlue).SprintFunc()
	binColor    = color.New(color.FgMagenta).SprintFunc()
	hiddenColor = color.New(color.FgHiCyan).SprintFunc()
)

// NewGaladh returns a new Galadh tree viewer.
// It takes arbitrary options, but doesn't need any.
func NewGaladh(options ...Options) *Galadh {
	galadh := &Galadh{
		w:      os.Stdout,
		level:  -1,
		glyphs: unicodeGlyphSet,
	}

	// Applying all passed options
	for _, fn := range options {
		fn(galadh)
	}

	return galadh
}

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
	// File Options
	humanReadable       bool
	printSize           bool
	replaceNonPrintable bool
	disableColor        bool

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

	fmt.Fprintln(g.w, dirColor(path))
	g.printDir(path, 0)

	// Print statistics
	dirLabel := "directories"
	fileLabel := "files"
	if g.cntDirs == 1 {
		dirLabel = "directory"
	}
	if g.cntFiles == 1 {
		fileLabel = "file"
	}
	fmt.Fprintf(g.w, "\n %v %v, %v %v\n", g.cntDirs, dirLabel, g.cntFiles, fileLabel)

	return nil
}

func (g *Galadh) printDir(path string, lvl int) {
	// Fetch content of directory
	files, err := g.getFiles(path)
	if err != nil {
		g.printError(err.Error())
		return
	}

	// For the moment do not handle symlinks and hardlinks
	// differently, but in a future version it should be
	// configurable how to handle links.
	lastItem := len(files) - 1
	for i := range files {
		f := files[i]
		g.printItem(path, f, i == lastItem)
		// Recursively descending into a directory
		if f.IsDir() {
			g.cntDirs++
			if g.level == -1 || lvl < g.level {
				g.indent(i == lastItem)
				g.printDir(filepath.Join(path, f.Name()), lvl+1)
				g.unindent()
			}
		} else {
			g.cntFiles++
		}
	}
}

// getFiles gets the content of a directory and applies all
// configured filters.
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
		if g.excludePattern != "" && g.matchPattern(g.excludePattern, f.Name()) {
			continue
		}

		files = append(files, f)
	}

	// Let's sort our files slice, first directories and then files in alphabetic order.
	sort.SliceStable(files, func(i, j int) bool {
		fa := files[i]
		fb := files[j]

		if fa.IsDir() && !fb.IsDir() {
			return true
		}

		if !fa.IsDir() && fb.IsDir() {
			return false
		}

		return strings.ToLower(fa.Name()) < strings.ToLower(fb.Name())
	})

	return files, nil
}

func (g Galadh) matchPattern(pattern, name string) bool {
	if pattern == "" {
		return true
	}

	if g.ignoreCase {
		pattern = strings.ToLower(pattern)
		name = strings.ToLower(name)
	}

	match, err := filepath.Match(pattern, name)
	if err != nil {
		return false
	}

	return match
}

func (g *Galadh) indent(lastItem bool) {
	if lastItem {
		g.indents = append(g.indents, strings.Repeat(" ", 4))
		return
	}

	g.indents = append(g.indents, g.glyphs.pipe+strings.Repeat(" ", 3))
}

func (g *Galadh) unindent() {
	if len(g.indents) > 0 {
		g.indents = g.indents[:len(g.indents)-1]
	}
}

func (g Galadh) printError(str string) {
	for i := range g.indents {
		fmt.Fprint(g.w, g.indents[i])
	}

	fmt.Fprintln(g.w, errColor(str))
}

func (g Galadh) printItem(path string, file os.FileInfo, lastItem bool) {
	for i := range g.indents {
		fmt.Fprint(g.w, g.indents[i])
	}

	if lastItem {
		fmt.Fprint(g.w, g.glyphs.last)
	} else {
		fmt.Fprint(g.w, g.glyphs.item)
	}

	// First print the metadata
	fmt.Fprintf(g.w, " %v", g.getMetaData(path, file))

	label := file.Name()
	if g.printFullPath {
		label = filepath.Join(path, label)
	}

	// Replace non-printable characters if requested
	if g.replaceNonPrintable {
		label = replaceNonPrintable(label)
	}

	// Use color function according file type
	if file.IsDir() {
		fmt.Fprintln(g.w, dirColor(label))
	} else if isHidden(file) {
		fmt.Fprintln(g.w, hiddenColor(label))
	} else if isExecutable(file) {
		fmt.Fprintln(g.w, binColor(label))
	} else {
		fmt.Fprintln(g.w, label)
	}
}

func (g Galadh) getMetaData(path string, file os.FileInfo) string {
	var meta []string
	if g.printSize || g.humanReadable {
		var size string
		if g.humanReadable {
			switch sz := file.Size(); {
			case sz < kbSize:
				size = fmt.Sprintf("%v", sz)
			case sz < mbSize:
				size = fmt.Sprintf("%.2fKb", float64(sz)/kbSize)
			case sz < gbSize:
				size = fmt.Sprintf("%.2fMb", float64(sz)/mbSize)
			case sz < tbSize:
				size = fmt.Sprintf("%.2fGb", float64(sz)/gbSize)
			case sz < pbSize:
				size = fmt.Sprintf("%.2fTb", float64(sz)/tbSize)
			default:
				size = fmt.Sprintf("%.2fPb", float64(sz)/pbSize)
			}
			size = fillPrefix(9, size)

		} else {
			size = fmt.Sprintf("%v", file.Size())
			// Use only 11 digits, otherwise show in scientific notation
			if len(size) > 11 {
				size = fmt.Sprintf("%E", float64(file.Size()))
			}
			size = fillPrefix(11, size)
		}
		meta = append(meta, size)
	}

	if len(meta) > 0 {
		return "[" + strings.Join(meta, " ") + "] "
	}

	return ""
}

func fillPrefix(n int, str string) string {
	cnt := len([]rune(str))
	if cnt >= n {
		return str
	}

	return fmt.Sprintf("%v%v", strings.Repeat(" ", n-cnt), str)
}

func replaceNonPrintable(str string) string {
	runes := make([]rune, 0, len(str))
	for _, r := range str {
		if !unicode.IsPrint(r) {
			runes = append(runes, '?')
			continue
		}
		runes = append(runes, r)
	}

	return string(runes)
}
