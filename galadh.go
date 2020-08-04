package galadh

import (
	"fmt"
	"os"
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

type printer interface {
	printRootDir(dir string)
	printItem(fi os.FileInfo, last bool)
	indent()
	unindent()
}

type treePrinter struct {
	glyphs   glyphSet
	indents  []string
	lastItem bool
}

func (t *treePrinter) indent() {

}

func (t *treePrinter) printRootDir(dir string) {
	fmt.Printf("+ %v\n", dir)
}

func (t *treePrinter) printItem(fi os.FileInfo, last bool) {
	for i := range t.indents {
		fmt.Print(i)
	}

	if last {
		fmt.Print(t.glyphs.last)
		t.lastItem = true
	} else {
		fmt.Print(t.glyphs.item)
		t.lastItem = false
	}

	fmt.Println(" ", fi.Name())
}

func PrintTest() {
	fmt.Println(asciiGlyphSet.pipe)
	fmt.Println(asciiGlyphSet.item)
	fmt.Println(asciiGlyphSet.last)
	fmt.Println(unicodeGlyphSet.pipe)
	fmt.Println(unicodeGlyphSet.item)
	fmt.Println(unicodeGlyphSet.last)
}
