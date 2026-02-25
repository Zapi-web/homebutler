// cmd/demo/main.go — Renders a static fake TUI dashboard for GIF recording.
// Uses the same Lip Gloss styles as the real TUI but prints to stdout (no alt screen).
package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const width = 100

var (
	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Padding(0, 2)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Padding(0, 2)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2)

	greenStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
	yellowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("226"))
	redStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("230"))
	dimStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62"))
	okStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("82"))

	barFull  = "\u2588"
	barEmpty = "\u2591"
)

func progressBar(percent float64, w int) string {
	filled := int(percent / 100 * float64(w))
	if filled > w {
		filled = w
	}
	var style lipgloss.Style
	switch {
	case percent >= 90:
		style = redStyle
	case percent >= 70:
		style = yellowStyle
	default:
		style = greenStyle
	}
	return style.Render(strings.Repeat(barFull, filled)) +
		dimStyle.Render(strings.Repeat(barEmpty, w-filled))
}

func main() {
	var b strings.Builder

	// Tabs
	tab1 := activeTabStyle.Render(" [1] mac-mini ")
	tab2 := inactiveTabStyle.Render(" [2] rpi5 ")
	serverCount := dimStyle.Render("  (2 available · Tab to switch)")
	b.WriteString(tab1 + tab2 + serverCount + "\n")

	// Left panel — System
	leftWidth := width/3 + 2
	barW := leftWidth - 20
	if barW < 8 {
		barW = 8
	}
	var sysLines []string
	sysLines = append(sysLines, titleStyle.Render("⚡ mac-mini"))
	sysLines = append(sysLines, "")
	sysLines = append(sysLines, fmt.Sprintf("  CPU  %s %5.1f%%", progressBar(23.4, barW), 23.4))
	sysLines = append(sysLines, fmt.Sprintf("  Mem  %s %5.1f%%", progressBar(49.2, barW), 49.2))
	sysLines = append(sysLines, fmt.Sprintf("  /    %s %5.1f%%", progressBar(3.0, barW), 3.0))
	sysLines = append(sysLines, "")
	sysLines = append(sysLines, "  Uptime:  42d 7h")
	sysLines = append(sysLines, "  OS:      darwin/arm64")
	sysLines = append(sysLines, "  Cores:   10")
	sysLines = append(sysLines, "  Memory:  7.9 / 16.0 GB")
	leftPanel := panelStyle.Width(leftWidth).Render(strings.Join(sysLines, "\n"))

	// Right panel — Docker
	rightWidth := width - leftWidth - 4
	var dockLines []string
	dockLines = append(dockLines, titleStyle.Render("Docker Containers"))
	dockLines = append(dockLines, "")
	hdr := fmt.Sprintf("  %-18s %-10s %-10s %s", "NAME", "STATE", "IMAGE", "STATUS")
	dockLines = append(dockLines, headerStyle.Render(hdr))

	type ctr struct {
		name, state, image, status string
	}
	containers := []ctr{
		{"homebridge", "running", "homebridge", "Up 12 days"},
		{"portainer", "running", "portainer", "Up 12 days"},
		{"pihole", "exited", "pihole", "Exited (0) 3d"},
		{"prometheus", "running", "prom/prom~", "Up 5 days"},
		{"grafana", "running", "grafana/g~", "Up 5 days"},
	}
	for _, c := range containers {
		stateStr := okStyle.Render(fmt.Sprintf("%-10s", c.state))
		if c.state == "exited" {
			stateStr = redStyle.Render(fmt.Sprintf("%-10s", c.state))
		}
		dockLines = append(dockLines, fmt.Sprintf("  %-18s %s %-10s %s", c.name, stateStr, c.image, c.status))
	}
	rightPanel := panelStyle.Width(rightWidth).Render(strings.Join(dockLines, "\n"))

	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel) + "\n")

	// Footer
	var footerParts []string
	alertLine := "  Alerts: " +
		okStyle.Render("CPU: 23%") + "  " +
		okStyle.Render("Mem: 49%") + "  " +
		okStyle.Render("Disk /: 3%")
	footerParts = append(footerParts, alertLine)

	keys := headerStyle.Render("Tab/Shift+Tab") + " switch server  │  " +
		headerStyle.Render("q") + " quit  │  " +
		dimStyle.Render("⟳ 2s")
	footerParts = append(footerParts, "  "+keys)
	footer := panelStyle.Width(width - 4).Render(strings.Join(footerParts, "\n"))
	b.WriteString(footer)

	fmt.Print(b.String())
}
