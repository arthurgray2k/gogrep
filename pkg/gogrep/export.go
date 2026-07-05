package gogrep

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// ExportFormat defines the format for exporting search results.
type ExportFormat string

const (
	FormatJSON ExportFormat = "json"
	FormatCSV  ExportFormat = "csv"
	FormatMD   ExportFormat = "markdown"
	FormatTXT  ExportFormat = "txt"
)

// Export writes the provided matches to w in the requested format.
func Export(w io.Writer, matches []Match, format ExportFormat) error {
	switch format {
	case FormatJSON:
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		return encoder.Encode(matches)
	case FormatCSV:
		writer := csv.NewWriter(w)
		// Header
		if err := writer.Write([]string{"Path", "LineNumber", "Text", "FileType", "FileSize", "CreateDate"}); err != nil {
			return err
		}
		for _, m := range matches {
			var ft, fs, cd string
			if m.Metadata != nil {
				ft = m.Metadata.Type
				fs = fmt.Sprintf("%d", m.Metadata.Size)
				cd = m.Metadata.CreateDate.Format("2006-01-02 15:04:05")
			}
			if err := writer.Write([]string{m.Path, fmt.Sprintf("%d", m.LineNumber), m.Text, ft, fs, cd}); err != nil {
				return err
			}
		}
		writer.Flush()
		return writer.Error()
	case FormatMD:
		fmt.Fprintln(w, "| Path | Line | Text | Type | Size | Created |")
		fmt.Fprintln(w, "|---|---|---|---|---|---|")
		for _, m := range matches {
			var ft, fs, cd string
			if m.Metadata != nil {
				ft = m.Metadata.Type
				fs = fmt.Sprintf("%d", m.Metadata.Size)
				cd = m.Metadata.CreateDate.Format("2006-01-02 15:04:05")
			}
			// Escape pipes in the text so it doesn't break the markdown table
			safeText := strings.ReplaceAll(m.Text, "|", "\\|")
			fmt.Fprintf(w, "| %s | %d | %s | %s | %s | %s |\n", m.Path, m.LineNumber, safeText, ft, fs, cd)
		}
		return nil
	case FormatTXT:
		for _, m := range matches {
			metaStr := ""
			if m.Metadata != nil {
				metaStr = fmt.Sprintf(" [%s | %d bytes | %s]", m.Metadata.Type, m.Metadata.Size, m.Metadata.CreateDate.Format("2006-01-02 15:04:05"))
			}
			fmt.Fprintf(w, "%s:%d%s: %s\n", m.Path, m.LineNumber, metaStr, m.Text)
		}
		return nil
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}
