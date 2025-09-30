package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Downloader handles ISO downloads
type Downloader struct {
	outputDir string
}

// NewDownloader creates a new downloader instance
func NewDownloader(outputDir string) *Downloader {
	return &Downloader{
		outputDir: outputDir,
	}
}

// DownloadISO downloads an ISO file
func (d *Downloader) DownloadISO(url, filename string) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(d.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create the full file path
	filepath := filepath.Join(d.outputDir, filename)

	// Download the file
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: %d %s", resp.StatusCode, resp.Status)
	}

	// Create the output file
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	// Copy the response body to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("✅ Downloaded: %s\n", filepath)
	return nil
}

// GetOutputDir prompts for output directory
func GetOutputDir() (string, error) {
	// For now, return a default directory
	// In a real implementation, you'd use a TUI directory picker
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	defaultDir := filepath.Join(homeDir, "Downloads", "isos")

	// Create directory if it doesn't exist
	if err := os.MkdirAll(defaultDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create default directory: %w", err)
	}

	return defaultDir, nil
}

// ValidateURL checks if a URL is valid
func ValidateURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}
