package gogrep

import (
	"strings"
	"testing"
)

func TestConfigResolver(t *testing.T) {
	jsonStr := `[
		{"path": ".", "pattern": "main"},
		{"path": "src", "ignore_case": true},
		{"path": "src/vendor", "pattern": "vendor_main"}
	]`

	resolver, err := NewConfigResolver(strings.NewReader(jsonStr))
	if err != nil {
		t.Fatalf("Failed to create resolver: %v", err)
	}

	base := Config{Pattern: "default"}

	// Root path should resolve to "main"
	cfg1 := resolver.Resolve(".", base)
	if cfg1.Pattern != "main" {
		t.Errorf("Expected pattern 'main', got %s", cfg1.Pattern)
	}
	if cfg1.IgnoreCase {
		t.Errorf("Expected ignore_case false")
	}

	// src path should inherit pattern "main" and set ignore_case true
	cfg2 := resolver.Resolve("src", base)
	if cfg2.Pattern != "main" {
		t.Errorf("Expected inherited pattern 'main', got %s", cfg2.Pattern)
	}
	if !cfg2.IgnoreCase {
		t.Errorf("Expected ignore_case true")
	}

	// vendor should get its own pattern and inherit ignore_case from src
	// using filepath matching for src/vendor/test.go
	cfg3 := resolver.Resolve("src/vendor/test.go", base)
	if cfg3.Pattern != "vendor_main" {
		t.Errorf("Expected pattern 'vendor_main', got %s", cfg3.Pattern)
	}
	if !cfg3.IgnoreCase {
		t.Errorf("Expected inherited ignore_case true for vendor")
	}
}
