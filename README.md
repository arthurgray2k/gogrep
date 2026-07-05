# gogrep

`gogrep` is a highly customizable, high-performance text search tool and library written in Go. 

Unlike standard `grep`, `gogrep` is explicitly designed to be easily embeddable into other Go applications while also serving as a robust standalone command-line tool. It introduces advanced hierarchical rule-based configuration, allowing developers to define highly granular, folder-specific search behaviors within a single repository.

## Key Features

- **Dual-Purpose Architecture:** Built with cleanly separated packages (`pkg/gogrep`), making it just as easy to import into your own Go backend as it is to run from the terminal.
- **Hierarchical Configuration:** Provide a JSON configuration to dictate folder-specific rules. For example, you can ignore hidden files globally, but include them specifically for a `./vendor` or `.git` subdirectory. Rules dynamically inherit from their parent directories.
- **Rich Metadata Extraction:** View not only the matching text but also intrinsic file properties such as MIME types, file sizes, and native file creation dates (with built-in OS-specific syscalls for Windows).
- **Multiple Export Formats:** Export your search results natively to structured JSON, cleanly formatted CSV, Markdown tables, or standard terminal text outputs.
- **Zero Dependencies:** Built strictly using Go's standard library to guarantee portability, tight security, and a minimal footprint.

## Documentation

See the [Usage Guide](usage.md) for detailed instructions on CLI flags, JSON configuration formatting, and code examples.

## Installation

To build the standalone command-line binary:
```bash
go build -o gogrep.exe ./cmd/gogrep
```
