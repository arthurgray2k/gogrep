package gogrep

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSearcherRun(t *testing.T) {
	tempDir := t.TempDir()

	// Create some files
	err := os.WriteFile(filepath.Join(tempDir, "test1.txt"), []byte("hello world\nfoo bar"), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	err = os.Mkdir(filepath.Join(tempDir, ".hidden"), 0755)
	if err != nil {
		t.Fatalf("Failed to create hidden dir: %v", err)
	}
	err = os.WriteFile(filepath.Join(tempDir, ".hidden", "test2.txt"), []byte("hello secret"), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	baseCfg := Config{
		Pattern:        "hello",
		FixedString:    true,
		ShowFileSize:   true,
		ShowCreateDate: true,
		ShowFileType:   true,
	}

	searcher := NewSearcher(baseCfg, nil)
	matches, err := searcher.Run(tempDir)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(matches) != 1 {
		t.Errorf("Expected 1 match, got %d", len(matches))
	}
	if len(matches) > 0 && matches[0].Text != "hello world" {
		t.Errorf("Expected match 'hello world', got %s", matches[0].Text)
	}

	// Test with IncludeHidden
	baseCfg.IncludeHidden = true
	searcher = NewSearcher(baseCfg, nil)
	matches, err = searcher.Run(tempDir)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(matches) != 2 {
		t.Errorf("Expected 2 matches including hidden, got %d", len(matches))
	}
}

func TestRegexSearch(t *testing.T) {
	tempDir := t.TempDir()
	os.WriteFile(filepath.Join(tempDir, "regex.txt"), []byte("Func main()"), 0644)

	baseCfg := Config{
		Pattern:    "func",
		IgnoreCase: true,
	}

	searcher := NewSearcher(baseCfg, nil)
	matches, _ := searcher.Run(tempDir)
	if len(matches) != 1 {
		t.Errorf("Expected 1 match, got %d", len(matches))
	}
}
