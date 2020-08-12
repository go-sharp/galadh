package galadh

import (
	"fmt"
	"io"
	"os"
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
