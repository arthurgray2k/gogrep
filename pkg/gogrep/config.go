package gogrep

import (
	"encoding/json"
	"io"
	"path/filepath"
	"sort"
	"strings"
)

// Config is the fully resolved configuration for a search operation.
// This is used directly by our library.
type Config struct {
	Pattern        string
	IncludeHidden  bool
	IgnoreCase     bool
	InvertMatch    bool
	FixedString    bool
	ShowFileSize   bool
	ShowCreateDate bool
	ShowFileType   bool
}

// TargetConfig is the raw JSON representation, allowing inheritance via pointers.
type TargetConfig struct {
	Path           string  `json:"path"`
	Pattern        *string `json:"pattern,omitempty"`
	IncludeHidden  *bool   `json:"include_hidden,omitempty"`
	IgnoreCase     *bool   `json:"ignore_case,omitempty"`
	InvertMatch    *bool   `json:"invert_match,omitempty"`
	FixedString    *bool   `json:"fixed_string,omitempty"`
	ShowFileSize   *bool   `json:"show_file_size,omitempty"`
	ShowCreateDate *bool   `json:"show_create_date,omitempty"`
	ShowFileType   *bool   `json:"show_file_type,omitempty"`
}

// ConfigResolver helps figure out the Config for any given path.
type ConfigResolver struct {
	rules []TargetConfig
}

// NewConfigResolver creates a resolver from a JSON stream.
func NewConfigResolver(r io.Reader) (*ConfigResolver, error) {
	var rules []TargetConfig
	if err := json.NewDecoder(r).Decode(&rules); err != nil {
		return nil, err
	}
	
	// Normalize paths in rules
	for i := range rules {
		rules[i].Path = filepath.Clean(rules[i].Path)
	}
	
	return &ConfigResolver{rules: rules}, nil
}

// Resolve returns the final Config for a targetPath.
// It starts with a base fallback config (usually populated from CLI arguments),
// and then overlays JSON configurations hierarchically.
func (cr *ConfigResolver) Resolve(targetPath string, fallback Config) Config {
	targetPath = filepath.Clean(targetPath)
	cfg := fallback

	// Find all rules that apply to this path (i.e. rule path is a prefix of targetPath)
	var applicableRules []TargetConfig
	for _, rule := range cr.rules {
		// If the rule is "." or a direct prefix (accounting for path separators)
		if rule.Path == "." || rule.Path == targetPath || strings.HasPrefix(targetPath, rule.Path+string(filepath.Separator)) {
			applicableRules = append(applicableRules, rule)
		}
	}

	// Sort applicable rules by path length (shortest first) so parents apply before children
	sort.Slice(applicableRules, func(i, j int) bool {
		return len(applicableRules[i].Path) < len(applicableRules[j].Path)
	})

	// Apply overrides
	for _, rule := range applicableRules {
		if rule.Pattern != nil {
			cfg.Pattern = *rule.Pattern
		}
		if rule.IncludeHidden != nil {
			cfg.IncludeHidden = *rule.IncludeHidden
		}
		if rule.IgnoreCase != nil {
			cfg.IgnoreCase = *rule.IgnoreCase
		}
		if rule.InvertMatch != nil {
			cfg.InvertMatch = *rule.InvertMatch
		}
		if rule.FixedString != nil {
			cfg.FixedString = *rule.FixedString
		}
		if rule.ShowFileSize != nil {
			cfg.ShowFileSize = *rule.ShowFileSize
		}
		if rule.ShowCreateDate != nil {
			cfg.ShowCreateDate = *rule.ShowCreateDate
		}
		if rule.ShowFileType != nil {
			cfg.ShowFileType = *rule.ShowFileType
		}
	}

	return cfg
}
