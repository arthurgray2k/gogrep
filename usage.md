# gogrep Usage Guide

`gogrep` is a powerful, customizable text search tool built in Go. It supports standard grep-like capabilities, extensive file metadata output, data export (JSON, CSV, Markdown, TXT), and complex hierarchical JSON configurations.

## Building / Installation

To build the executable:
```bash
go build -o gogrep.exe ./cmd/gogrep
```

---

## Command-Line Flags

| Flag | Description |
|---|---|
| `--pattern` | The text or regex pattern to search for. |
| `--dir` | The root directory to search in (default is `.`). |
| `--i` | Perform a case-insensitive search. |
| `--F` | Perform a fixed string match (disables Regex). |
| `--v` | Invert the match (returns lines that *do not* match). |
| `--hidden` | Include hidden files and directories in the search. |
| `--size` | Display the file size in the output. |
| `--date` | Display the file creation date in the output. |
| `--type` | Display the file type (MIME) in the output. |
| `--export-format` | The format for exporting results: `txt` (default), `json`, `csv`, or `markdown`. |
| `--export-file` | The path to save the exported data. If omitted, results print to `stdout`. |
| `--config` | Path to a JSON configuration file for advanced hierarchical rules. |

---

## Basic Examples

**Simple search in current directory:**
```bash
./gogrep.exe --pattern "func main"
```

**Case-insensitive fixed string search in a specific directory:**
```bash
./gogrep.exe --pattern "TODO" --i --F --dir ./src
```

**Exporting results to a CSV file with metadata:**
```bash
./gogrep.exe --pattern "error" --size --date --export-format csv --export-file errors.csv
```

---

## Advanced JSON Configuration

For complex projects, you can use a JSON configuration file to define rules for specific folders. Rules naturally inherit from their parent folders unless explicitly overridden.

**example_rules.json**
```json
[
  {
    "path": ".",
    "pattern": "panic",
    "include_hidden": false,
    "show_file_size": true,
    "show_create_date": true,
    "show_file_type": true
  },
  {
    "path": "./vendor",
    "include_hidden": true,
    "ignore_case": true
  }
]
```

**Using the configuration file:**
```bash
./gogrep.exe --config example_rules.json --dir .
```

*In the example above, the `./vendor` directory inherits the `pattern`, `show_file_size`, and other metadata flags from the root (`.`), but it enables case-insensitivity and searching in hidden files just for itself.*
