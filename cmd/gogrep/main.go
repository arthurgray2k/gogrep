package main

import (
	"flag"
	"fmt"
	"os"

	"gogrep/pkg/gogrep"
)

func main() {
	var (
		pattern        = flag.String("pattern", "", "Pattern to search for (base configuration)")
		configPath     = flag.String("config", "", "Path to JSON config file")
		includeHidden  = flag.Bool("hidden", false, "Include hidden files/folders")
		ignoreCase     = flag.Bool("i", false, "Ignore case")
		invertMatch    = flag.Bool("v", false, "Invert match")
		fixedString    = flag.Bool("F", false, "Fixed string match (not regex)")
		showFileSize   = flag.Bool("size", false, "Show file size in output")
		showCreateDate = flag.Bool("date", false, "Show file create date in output")
		showFileType   = flag.Bool("type", false, "Show file type in output")
		exportFormat   = flag.String("export-format", "txt", "Export format: json, csv, markdown, txt")
		exportFile     = flag.String("export-file", "", "File to save the export to (default is stdout)")
		rootDir        = flag.String("dir", ".", "Root directory to search")
	)
	flag.Parse()

	// Setup Base Config from CLI arguments
	baseCfg := gogrep.Config{
		Pattern:        *pattern,
		IncludeHidden:  *includeHidden,
		IgnoreCase:     *ignoreCase,
		InvertMatch:    *invertMatch,
		FixedString:    *fixedString,
		ShowFileSize:   *showFileSize,
		ShowCreateDate: *showCreateDate,
		ShowFileType:   *showFileType,
	}

	// Setup JSON Resolver if provided
	var resolver *gogrep.ConfigResolver
	if *configPath != "" {
		f, err := os.Open(*configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening config file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()

		resolver, err = gogrep.NewConfigResolver(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing config file: %v\n", err)
			os.Exit(1)
		}
	}

	// Require a base pattern if no config is given
	if baseCfg.Pattern == "" && resolver == nil {
		fmt.Fprintf(os.Stderr, "Error: pattern is required (via -pattern or --config)\n")
		flag.Usage()
		os.Exit(1)
	}

	// Run Search
	searcher := gogrep.NewSearcher(baseCfg, resolver)
	matches, err := searcher.Run(*rootDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Search error: %v\n", err)
		os.Exit(1)
	}

	// Output results
	var out = os.Stdout
	if *exportFile != "" {
		f, err := os.Create(*exportFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating export file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		out = f
	}

	if err := gogrep.Export(out, matches, gogrep.ExportFormat(*exportFormat)); err != nil {
		fmt.Fprintf(os.Stderr, "Export error: %v\n", err)
		os.Exit(1)
	}
}
