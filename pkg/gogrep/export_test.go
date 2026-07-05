package gogrep

import (
	"bytes"
	"strings"
	"testing"
)

func TestExportJSON(t *testing.T) {
	matches := []Match{
		{Path: "test.go", LineNumber: 1, Text: "func main()"},
	}
	var buf bytes.Buffer
	err := Export(&buf, matches, FormatJSON)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}
	if !strings.Contains(buf.String(), `"Path": "test.go"`) {
		t.Errorf("Expected JSON to contain Path, got %s", buf.String())
	}
}

func TestExportCSV(t *testing.T) {
	matches := []Match{
		{Path: "test.go", LineNumber: 1, Text: "func main()"},
	}
	var buf bytes.Buffer
	err := Export(&buf, matches, FormatCSV)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}
	if !strings.Contains(buf.String(), "test.go,1,func main(),,,") {
		t.Errorf("Expected CSV formatting, got %s", buf.String())
	}
}

func TestExportMarkdown(t *testing.T) {
	matches := []Match{
		{Path: "test.go", LineNumber: 1, Text: "func main() { | }"},
	}
	var buf bytes.Buffer
	err := Export(&buf, matches, FormatMD)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}
	// Verify escaping of pipes
	if !strings.Contains(buf.String(), "func main() { \\| }") {
		t.Errorf("Expected escaped pipes in markdown, got %s", buf.String())
	}
}

func TestExportTXT(t *testing.T) {
	matches := []Match{
		{Path: "test.go", LineNumber: 1, Text: "func main()"},
	}
	var buf bytes.Buffer
	err := Export(&buf, matches, FormatTXT)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}
	if !strings.Contains(buf.String(), "test.go:1: func main()") {
		t.Errorf("Expected TXT formatting, got %s", buf.String())
	}
}

func TestExportUnsupported(t *testing.T) {
	var buf bytes.Buffer
	err := Export(&buf, []Match{}, ExportFormat("unsupported"))
	if err == nil {
		t.Errorf("Expected error for unsupported format, got nil")
	}
}
