package main

import (
	"fmt"
	"os"

	"github.com/go-sharp/galadh"
	flag "github.com/jessevdk/go-flags"
)

// Version
var gitCommit string
var gitTag string

var options cmdOptions
var parser = flag.NewParser(&options, flag.HelpFlag|flag.PassDoubleDash)

func main() {
	_, err := parser.Parse()
	if err != nil {
		if e, ok := err.(*flag.Error); ok && e.Type == flag.ErrHelp {
			fmt.Printf("galadh \u00A9 go-sharp\ncommit : %v\nversion: %v\n\n", gitCommit, gitTag)
			fmt.Println(err)
		} else {
			fmt.Println("Galadh: failed to parse arguments:", err)
		}
		os.Exit(1)
	}

	maxLevel := -1
	if options.Level != nil && *options.Level >= 0 {
		maxLevel = *options.Level
	}

	opts := []galadh.Options{
		galadh.UseAllFiles(options.AllFiles),
		galadh.UseDirOnly(options.DirOnly),
		galadh.UseExcludePattern(options.ExcludePattern),
		galadh.UseIncludePattern(options.IncludePattern),
		galadh.UseIgnoreCase(options.IgnoreCase),
		galadh.UseMaxLevel(maxLevel),
		galadh.UseMatchDirs(options.MatchDirs),
		galadh.UsePrintFullPath(options.PrintFullPath),
		galadh.UseHumanReadableSize(options.HumanReadable),
		galadh.UseSize(options.PrintSize),
		galadh.UseReplaceNonPrintableChars(options.ReplaceNonPrintable),
		galadh.UseASCIIGlyphs(options.ASCII),
		galadh.UseColors(!options.DisableColors),
	}

	if options.FileOutput != "" {
		f, err := os.OpenFile(options.FileOutput, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
		if err != nil {
			fmt.Println("Galadh: failed to open file:", err)
			os.Exit(1)
		}
		defer f.Close()
		opts = append(opts, galadh.UseOutput(f))
	}

	g := galadh.NewGaladh(opts...)
	if options.Args.Path == "" {
		options.Args.Path = "."
	}

	if err := g.PrintTree(options.Args.Path); err != nil {
		fmt.Println("Galadh:", err)
		os.Exit(2)
	}
}

type cmdOptions struct {
	// Listing Options
	AllFiles       bool   `short:"a" long:"all" description:"All files are listed."`
	DirOnly        bool   `short:"d" long:"dirs-only" description:"List directories only."`
	ExcludePattern string `short:"I" long:"exclude-pattern" description:"Do not list files that match the given pattern."`
	IncludePattern string `short:"P" long:"include-pattern" description:"List only those files that match the pattern given."`
	IgnoreCase     bool   `long:"ignore-case" description:"Ignore case when pattern matching."`
	Level          *int   `short:"L" long:"level" description:"Descend only level directories deep."`
	MatchDirs      bool   `long:"matchdirs" description:"Include directory names in -P pattern matching."`
	PrintFullPath  bool   `short:"f" long:"fullpath" description:"Print the full path prefix for each file."`
	// File Options
	HumanReadable       bool   `short:"H" long:"human-readable" description:"Print the size in a more human readable way."`
	PrintSize           bool   `short:"s" long:"size" description:"Print the size in bytes of each file."`
	ReplaceNonPrintable bool   `short:"q" long:"replace-nonprintable" description:"Print non-printable characters as '?'."`
	FileOutput          string `short:"o" long:"output" description:"Output to file instead of stdout."`
	ASCII               bool   `short:"A" long:"ascii" description:"Print ASCII only indentation lines."`
	DisableColors       bool   `short:"n" long:"no-colors" description:"Turn colorization off."`
	// Rest parameters
	Args struct {
		Path string `positional-arg-name:"PATH"`
	} `positional-args:"yes"`
}
