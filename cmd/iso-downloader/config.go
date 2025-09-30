package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Config represents the configuration structure
type Config struct {
	Families map[string]Family `json:"families"`
}

// Family represents a distribution family
type Family struct {
	Distros map[string]Distro `json:"distros"`
}

// Distro represents a distribution
type Distro struct {
	Versions []string `json:"versions"`
	BaseURL  string   `json:"base_url"`
}

// LoadConfig loads configuration from local file or GitHub
func LoadConfig() (*Config, error) {
	// Try local file first
	localPath := "data/distros.json"
	if _, err := os.Stat(localPath); err == nil {
		return loadLocalConfig(localPath)
	}

	// Fallback to GitHub
	return loadRemoteConfig()
}

// loadLocalConfig loads configuration from local file
func loadLocalConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	return &config, nil
}

// loadRemoteConfig loads configuration from GitHub
func loadRemoteConfig() (*Config, error) {
	url := "https://raw.githubusercontent.com/yourname/iso-downloader/main/data/distros.json"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch config: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch config: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := json.Unmarshal(body, &config); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	return &config, nil
}

// ResolveISOURL resolves the ISO URL for a given distro and version
func (c *Config) ResolveISOURL(family, distro, version string) (string, error) {
	familyData, exists := c.Families[family]
	if !exists {
		return "", fmt.Errorf("family %s not found", family)
	}

	distroData, exists := familyData.Distros[distro]
	if !exists {
		return "", fmt.Errorf("distro %s not found in family %s", distro, family)
	}

	// For now, return a simple URL construction
	// In a real implementation, you'd resolve the actual ISO URL
	baseURL := distroData.BaseURL
	if baseURL == "" {
		baseURL = "https://releases.ubuntu.com/" // Default fallback
	}

	// Simple URL construction - this would need to be more sophisticated
	// for different distros with different URL patterns
	isoURL := fmt.Sprintf("%s/%s/%s", baseURL, version, "ubuntu-"+version+"-desktop-amd64.iso")

	return isoURL, nil
}
