package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TUI model for the interactive interface
type model struct {
	families        list.Model
	distros         list.Model
	versions        list.Model
	config          *Config
	state           string
	selectedFamily  string
	selectedDistro  string
	selectedVersion string
	asciiArt        string
	showDownload    bool
	shouldDownload  bool
}

// Init initializes the TUI model
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles TUI updates
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Update list sizes when terminal resizes
		m.families.SetWidth(msg.Width - 4)
		m.families.SetHeight(msg.Height - 6)
		if m.distros.Items() != nil {
			m.distros.SetWidth(msg.Width - 4)
			m.distros.SetHeight(msg.Height - 6)
		}
		if m.versions.Items() != nil {
			m.versions.SetWidth(msg.Width - 4)
			m.versions.SetHeight(msg.Height - 6)
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			return m.handleEnter()
		case "esc":
			return m.handleEsc()
		case "d":
			return m.handleDownload()
		}
	}

	var cmd tea.Cmd
	switch m.state {
	case "family":
		m.families, cmd = m.families.Update(msg)
	case "distro":
		m.distros, cmd = m.distros.Update(msg)
	case "version":
		m.versions, cmd = m.versions.Update(msg)
	}
	return m, cmd
}

// handleEnter handles the enter key press
func (m model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.state {
	case "family":
		selected := m.families.SelectedItem()
		if selected != nil {
			m.selectedFamily = selected.FilterValue()
			fmt.Printf("🔍 Selected family: %s\n", m.selectedFamily)
			m.setupDistros()
			m.state = "distro"
		}
		return m, nil
	case "distro":
		selected := m.distros.SelectedItem()
		if selected != nil {
			m.selectedDistro = selected.FilterValue()
			fmt.Printf("🔍 Selected distro: %s\n", m.selectedDistro)
			m.setupVersions()
			m.setupAsciiArt()
			m.state = "version"
		}
		return m, nil
	case "version":
		selected := m.versions.SelectedItem()
		if selected != nil {
			m.selectedVersion = selected.FilterValue()
			m.showDownload = true
			fmt.Printf("🔍 Selected version: %s\n", m.selectedVersion)
		}
		return m, nil
	}
	return m, nil
}

// handleEsc handles the escape key press
func (m model) handleEsc() (tea.Model, tea.Cmd) {
	switch m.state {
	case "distro":
		m.state = "family"
	case "version":
		m.state = "distro"
	}
	return m, nil
}

// handleDownload handles the download key press
func (m model) handleDownload() (tea.Model, tea.Cmd) {
	if m.selectedFamily != "" && m.selectedDistro != "" && m.selectedVersion != "" {
		m.shouldDownload = true
		fmt.Printf("🔍 Starting download for %s %s %s\n", m.selectedFamily, m.selectedDistro, m.selectedVersion)
		return m, tea.Quit
	}
	return m, nil
}

// setupDistros populates the distros list
func (m *model) setupDistros() {
	var items []list.Item
	for name := range m.config.Families[m.selectedFamily].Distros {
		items = append(items, item{title: name, desc: "Linux distribution"})
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Foreground(lipgloss.Color("205")).
		Padding(0, 0, 0, 1)
	delegate.Styles.SelectedDesc = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))
	delegate.Styles.NormalTitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("7"))
	delegate.Styles.NormalDesc = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	m.distros = list.New(items, delegate, 0, 0)
	m.distros.Title = ""
	m.distros.SetShowStatusBar(false)
	m.distros.SetFilteringEnabled(false)
	m.distros.SetWidth(80)
	m.distros.SetHeight(15)
}

// setupVersions populates the versions list
func (m *model) setupVersions() {
	var items []list.Item
	versions := m.config.Families[m.selectedFamily].Distros[m.selectedDistro].Versions

	for _, version := range versions {
		items = append(items, item{title: version, desc: "Download " + version})
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Foreground(lipgloss.Color("205")).
		Padding(0, 0, 0, 1)
	delegate.Styles.SelectedDesc = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))
	delegate.Styles.NormalTitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("7"))
	delegate.Styles.NormalDesc = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	m.versions = list.New(items, delegate, 0, 0)
	m.versions.Title = ""
	m.versions.SetShowStatusBar(false)
	m.versions.SetFilteringEnabled(false)
	m.versions.SetWidth(80)
	m.versions.SetHeight(15)
}

// setupAsciiArt sets up ASCII art for the selected distro
func (m *model) setupAsciiArt() {
	asciiArts := map[string]string{
		"debian": `    ██████╗ ███████╗██████╗ ██╗ █████╗ ███╗   ██╗
    ██╔══██╗██╔════╝██╔══██╗██║██╔══██╗████╗  ██║
    ██║  ██║█████╗  ██████╔╝██║███████║██╔██╗ ██║
    ██║  ██║██╔══╝  ██╔══██╗██║██╔══██║██║╚██╗██║
    ██████╔╝███████╗██║  ██║██║██║  ██║██║ ╚████║
    ╚═════╝ ╚══════╝╚═╝  ╚═╝╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝`,
		"ubuntu": `    ██╗   ██╗██████╗ ██╗   ██╗ █████╗ ███╗   ██╗██╗   ██╗
    ██║   ██║██╔══██╗██║   ██║██╔══██╗████╗  ██║██║   ██║
    ██║   ██║██████╔╝██║   ██║███████║██╔██╗ ██║██║   ██║
    ██║   ██║██╔══██╗██║   ██║██╔══██║██║╚██╗██║██║   ██║
    ╚██████╔╝██║  ██║╚██████╔╝██║  ██║██║ ╚████║╚██████╔╝
     ╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═══╝ ╚═════╝`,
		"fedora": `    ███████╗███████╗██████╗  ██████╗ ██████╗  █████╗ 
    ██╔════╝██╔════╝██╔══██╗██╔═══██╗██╔══██╗██╔══██╗
    █████╗  █████╗  ██║  ██║██║   ██║██████╔╝███████║
    ██╔══╝  ██╔══╝  ██║  ██║██║   ██║██╔══██╗██╔══██║
    ██║     ███████╗██████╔╝╚██████╔╝██║  ██║██║  ██║
    ╚═╝     ╚══════╝╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝`,
		"arch": `    █████╗ ██████╗  ██████╗██╗  ██╗
    ██╔══██╗██╔══██╗██╔════╝██║  ██║
    ███████║██████╔╝██║     ███████║
    ██╔══██║██╔══██╗██║     ██╔══██║
    ██║  ██║██║  ██║╚██████╗██║  ██║
    ╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝`,
	}
	m.asciiArt = asciiArts[m.selectedDistro]
}

// View renders the TUI
func (m model) View() string {
	// Title style
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Align(lipgloss.Center).
		Margin(1, 0)

	// ASCII art style
	asciiStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Align(lipgloss.Center).
		Margin(1, 0)

	// Selection status style
	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("2")).
		Bold(true)

	var content string
	var title string

	switch m.state {
	case "family":
		title = "🌍 Select Distribution Family"
		content = m.families.View()
	case "distro":
		title = "🐧 Select Distribution"
		content = m.distros.View()
	case "version":
		title = "📦 Select Version"
		content = m.versions.View()
	}

	// Add ASCII art if available
	if m.asciiArt != "" {
		content = asciiStyle.Render(m.asciiArt) + "\n" + content
	}

	// Add selection status
	statusText := ""
	if m.selectedFamily != "" {
		statusText += "✅ Family: " + m.selectedFamily
	}
	if m.selectedDistro != "" {
		statusText += " | ✅ Distro: " + m.selectedDistro
	}
	if m.selectedVersion != "" {
		statusText += " | ✅ Version: " + m.selectedVersion
	}

	if statusText != "" {
		content = statusStyle.Render(statusText) + "\n\n" + content
	}

	// Combine title and content
	fullContent := titleStyle.Render(title) + "\n" + content

	// Add help text at the bottom
	helpText := "↑/↓ Navigate • Enter Select • Esc Back • q Quit"
	if m.showDownload {
		helpText += " • d Download"
	}

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Align(lipgloss.Center)

	return fullContent + "\n" + helpStyle.Render(helpText)
}

// item represents a list item
type item struct {
	title, desc string
}

func (i item) FilterValue() string { return i.title }
func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }

// NewTUI creates a new TUI model
func NewTUI(config *Config) model {
	// Create family items
	var items []list.Item
	for name := range config.Families {
		items = append(items, item{title: name, desc: "Distribution family"})
	}

	// Debug: Print items being created
	fmt.Printf("🔍 Creating TUI with %d family items\n", len(items))

	// Create delegate with custom styling
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Foreground(lipgloss.Color("205")).
		Padding(0, 0, 0, 1)
	delegate.Styles.SelectedDesc = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))
	delegate.Styles.NormalTitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("7"))
	delegate.Styles.NormalDesc = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	m := model{
		families: list.New(items, delegate, 0, 0),
		config:   config,
		state:    "family",
	}
	m.families.Title = ""
	m.families.SetShowStatusBar(false)
	m.families.SetFilteringEnabled(false)
	m.families.SetWidth(80)
	m.families.SetHeight(15)

	return m
}
