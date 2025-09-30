package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	// Load configuration
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Debug: Print loaded families
	fmt.Printf("🔍 Loaded %d families:\n", len(config.Families))
	for name, family := range config.Families {
		fmt.Printf("  - %s: %d distros\n", name, len(family.Distros))
	}

	// Get output directory
	outputDir, err := GetOutputDir()
	if err != nil {
		return fmt.Errorf("failed to get output directory: %w", err)
	}

	fmt.Printf("📁 Output directory: %s\n", outputDir)

	// Create TUI
	m := NewTUI(config)

	// Run TUI
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	// Debug: Print what was selected
	fmt.Printf("🔍 TUI exited. Selections:\n")
	fmt.Printf("  - Family: '%s'\n", m.selectedFamily)
	fmt.Printf("  - Distro: '%s'\n", m.selectedDistro)
	fmt.Printf("  - Version: '%s'\n", m.selectedVersion)
	fmt.Printf("  - Should Download: %t\n", m.shouldDownload)

	// After TUI exits, handle the download
	if m.shouldDownload && m.selectedFamily != "" && m.selectedDistro != "" && m.selectedVersion != "" {
		// Resolve ISO URL
		isoURL, err := config.ResolveISOURL(m.selectedFamily, m.selectedDistro, m.selectedVersion)
		if err != nil {
			return fmt.Errorf("failed to resolve ISO URL: %w", err)
		}

		// Create downloader
		downloader := NewDownloader(outputDir)

		// Generate filename
		filename := fmt.Sprintf("%s-%s-%s.iso", m.selectedDistro, m.selectedVersion, "amd64")

		// Download ISO
		fmt.Printf("🚀 Starting download...\n")
		fmt.Printf("📥 URL: %s\n", isoURL)
		fmt.Printf("💾 Filename: %s\n", filename)

		if err := downloader.DownloadISO(isoURL, filename); err != nil {
			return fmt.Errorf("download failed: %w", err)
		}

		fmt.Printf("🎉 Download completed successfully!\n")
	} else {
		if !m.shouldDownload {
			fmt.Println("❌ No download requested. Exiting...")
		} else {
			fmt.Println("❌ Incomplete selection. Exiting...")
		}
	}

	return nil
}
