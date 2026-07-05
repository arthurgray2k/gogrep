package gogrep

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Match represents a single line match in a file.
type Match struct {
	Path       string
	LineNumber int
	Text       string
	Metadata   *FileMetadata // Present if requested
}

// Searcher handles the search operation across the filesystem.
type Searcher struct {
	BaseConfig Config
	Resolver   *ConfigResolver
}

// NewSearcher creates a new Searcher. If resolver is nil, it acts with just BaseConfig.
func NewSearcher(base Config, resolver *ConfigResolver) *Searcher {
	if resolver == nil {
		resolver = &ConfigResolver{}
	}
	return &Searcher{
		BaseConfig: base,
		Resolver:   resolver,
	}
}

// Run executes the search starting from rootDir and returns a slice of Matches.
func (s *Searcher) Run(rootDir string) ([]Match, error) {
	var results []Match

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Skip unreadable files/directories instead of completely failing
			return nil 
		}

		// Resolve config for this path
		cfg := s.Resolver.Resolve(path, s.BaseConfig)

		// Check hidden files/folders (ignore current dir "." root check)
		if !cfg.IncludeHidden && isHidden(d.Name()) && path != rootDir {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			return nil
		}

		// If a pattern is not set, we can't grep it.
		if cfg.Pattern == "" {
			return nil
		}

		matches, err := s.grepFile(path, d, cfg)
		if err != nil {
			// We skip files we can't read
			return nil
		}
		results = append(results, matches...)
		return nil
	})

	return results, err
}

func (s *Searcher) grepFile(path string, d fs.DirEntry, cfg Config) ([]Match, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var patternRegex *regexp.Regexp
	if !cfg.FixedString {
		expr := cfg.Pattern
		if cfg.IgnoreCase {
			expr = "(?i)" + expr
		}
		patternRegex, err = regexp.Compile(expr)
		if err != nil {
			return nil, fmt.Errorf("invalid regex pattern: %w", err)
		}
	}

	var metadata *FileMetadata
	if cfg.ShowFileSize || cfg.ShowCreateDate || cfg.ShowFileType {
		info, err := d.Info()
		if err == nil {
			m := GetMetadata(path, info)
			metadata = &m
		}
	}

	var fileMatches []Match
	scanner := bufio.NewScanner(file)
	lineNum := 1
	for scanner.Scan() {
		line := scanner.Text()
		matched := false

		if cfg.FixedString {
			if cfg.IgnoreCase {
				matched = strings.Contains(strings.ToLower(line), strings.ToLower(cfg.Pattern))
			} else {
				matched = strings.Contains(line, cfg.Pattern)
			}
		} else {
			matched = patternRegex.MatchString(line)
		}

		if cfg.InvertMatch {
			matched = !matched
		}

		if matched {
			fileMatches = append(fileMatches, Match{
				Path:       path,
				LineNumber: lineNum,
				Text:       line,
				Metadata:   metadata,
			})
		}
		lineNum++
	}

	return fileMatches, scanner.Err()
}

func isHidden(name string) bool {
	return strings.HasPrefix(name, ".") && name != "." && name != ".."
}
