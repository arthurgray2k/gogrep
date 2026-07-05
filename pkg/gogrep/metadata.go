package gogrep

import (
	"mime"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

// FileMetadata holds the extra details we want to output for files.
type FileMetadata struct {
	Type       string
	Size       int64
	CreateDate time.Time
}

// GetMetadata extracts file information from os.FileInfo.
// It uses OS-specific sys-calls to fetch the creation time on Windows.
func GetMetadata(path string, info os.FileInfo) FileMetadata {
	meta := FileMetadata{
		Size: info.Size(),
		Type: mime.TypeByExtension(filepath.Ext(path)),
	}
	
	// If mime type is empty, default to something
	if meta.Type == "" {
		meta.Type = "application/octet-stream"
	}

	// Try to extract creation time (Birth time).
	// On Windows, Sys() returns *syscall.Win32FileAttributeData
	if stat, ok := info.Sys().(*syscall.Win32FileAttributeData); ok {
		meta.CreateDate = time.Unix(0, stat.CreationTime.Nanoseconds())
	} else {
		// Fallback for non-Windows where Sys() doesn't expose birth time uniformly.
		meta.CreateDate = info.ModTime()
	}

	return meta
}
