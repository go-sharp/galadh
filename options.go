package galadh

import (
	"io"

	"github.com/go-sharp/color"
)

// Options is a function that configures the Galadh instance.
type Options func(g *Galadh)

// UseAllFiles configures whether to include also hidden files.
func UseAllFiles(b bool) Options {
	return func(g *Galadh) { g.allFiles = b }
}

// UseDirOnly configures whether
func UseDirOnly(b bool) Options {
	return func(g *Galadh) { g.dirOnly = b }
}

// UseExcludePattern configures which files should be excluded.
func UseExcludePattern(pattern string) Options {
	return func(g *Galadh) { g.excludePattern = pattern }
}

// UseIncludePattern configures which files should be included.
func UseIncludePattern(pattern string) Options {
	return func(g *Galadh) { g.includePattern = pattern }
}

// UseIgnoreCase configures whether pattern matching ignore casing.
func UseIgnoreCase(b bool) Options {
	return func(g *Galadh) { g.ignoreCase = b }
}

// UseMaxLevel configures how deep the recursion of the tree will be.
// If -1 is used, the recursion has no limits.
func UseMaxLevel(lvl int) Options {
	return func(g *Galadh) {
		if lvl < -1 {
			g.level = -1
			return
		}
		g.level = lvl
	}
}

// UseMatchDirs configures if the pattern matches also applies for directories.
// If a pattern matches for a directory, then all files and folders in that
// directory will be printed as well.
func UseMatchDirs(b bool) Options {
	return func(g *Galadh) { g.matchDirs = b }
}

// UsePrintFullPath configures whether the full path should be displayed.
func UsePrintFullPath(b bool) Options {
	return func(g *Galadh) { g.printFullPath = b }
}

// UseHumanReadableSize configures whether to display file size in bytes or human readable units.
func UseHumanReadableSize(b bool) Options {
	return func(g *Galadh) { g.humanReadable = b }
}

// UseSize configures whether to print size of a file or not.
func UseSize(b bool) Options {
	return func(g *Galadh) { g.printSize = b }
}

// UseReplaceNonPrintableChars to configure to replace non printable characters.
func UseReplaceNonPrintableChars(b bool) Options {
	return func(g *Galadh) { g.replaceNonPrintable = b }
}

// UseOutput configures to print output to the specified writer.
// This options disables colorization and can't be overriden by UseColors.
func UseOutput(writer io.Writer) Options {
	return func(g *Galadh) {
		g.disableColor = true
		g.w = writer
		color.NoColor = true
	}
}

// UseColors configures whether to use colorized output or not.
func UseColors(b bool) Options {
	return func(g *Galadh) {
		if g.disableColor {
			return
		}
		color.NoColor = !b
	}
}

// UseASCIIGlyphs configures whether to use ASCII glyphs or not.
func UseASCIIGlyphs(b bool) Options {
	return func(g *Galadh) {
		if b {
			g.glyphs = asciiGlyphSet
		} else {
			g.glyphs = unicodeGlyphSet
		}
	}
}
