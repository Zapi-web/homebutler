// cmd/demo/main.go — Renders fake CLI/TUI output for GIF recording.
// Uses the same Lip Gloss styles as the real TUI but prints to stdout (no alt screen).
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const tuiWidth = 100

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
	cmd := "watch"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "status":
		printStatus()
	case "docker":
		printDocker()
	case "ports":
		printPorts()
	case "alerts":
		printAlerts()
	case "watch":
		printTUI()
	}
}

func printStatus() {
	emoji := lipgloss.NewStyle().Bold(true)
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	fmt.Println(emoji.Render("🖥  homelab-server") + dim.Render(" (linux/arm64)"))
	fmt.Println(dim.Render("   Uptime:  ") + "42d 7h")
	fmt.Println(dim.Render("   CPU:     ") + greenStyle.Render("23.4%") + dim.Render(" (10 cores)"))
	fmt.Println(dim.Render("   Memory:  ") + "7.9 / 16.0 GB " + greenStyle.Render("(49.2%)"))
	fmt.Println(dim.Render("   Disk /:  ") + "11 / 460 GB " + greenStyle.Render("(3%)"))
}

func printDocker() {
	hdr := headerStyle.Render(fmt.Sprintf("%-20s %-12s %-15s %s", "NAME", "STATE", "IMAGE", "STATUS"))
	fmt.Println(hdr)

	type ctr struct{ name, state, image, status string }
	containers := []ctr{
		{"homebridge", "running", "homebridge/hb", "Up 12 days"},
		{"portainer", "running", "portainer/ce", "Up 12 days"},
		{"pihole", "exited", "pihole/pihole", "Exited (0) 3d ago"},
		{"prometheus", "running", "prom/prometheus", "Up 5 days"},
		{"grafana", "running", "grafana/grafana", "Up 5 days"},
	}
	for _, c := range containers {
		st := okStyle.Render(fmt.Sprintf("%-12s", c.state))
		if c.state == "exited" {
			st = redStyle.Render(fmt.Sprintf("%-12s", c.state))
		}
		fmt.Printf("%-20s %s %-15s %s\n", c.name, st, c.image, c.status)
	}
}

func printPorts() {
	hdr := headerStyle.Render(fmt.Sprintf("%-8s %-8s %-24s %s", "PORT", "PROTO", "PROCESS", "PID"))
	fmt.Println(hdr)

	type port struct{ port, proto, proc, pid string }
	ports := []port{
		{"22", "TCP", "sshd", "1234"},
		{"53", "UDP", "pihole-FTL", "5678"},
		{"80", "TCP", "nginx", "9012"},
		{"443", "TCP", "nginx", "9012"},
		{"3000", "TCP", "grafana", "3456"},
		{"8080", "TCP", "homebridge", "7890"},
		{"9090", "TCP", "prometheus", "2345"},
	}
	for _, p := range ports {
		fmt.Printf("%-8s %-8s %-24s %s\n", p.port, p.proto, p.proc, p.pid)
	}
}

func printAlerts() {
	fmt.Printf("  %s  %s\n", okStyle.Render("●"), fmt.Sprintf("CPU      %s  (threshold: 90%%)", greenStyle.Render("23.4%")))
	fmt.Printf("  %s  %s\n", okStyle.Render("●"), fmt.Sprintf("Memory   %s  (threshold: 85%%)", greenStyle.Render("49.2%")))
	fmt.Printf("  %s  %s\n", okStyle.Render("●"), fmt.Sprintf("Disk /   %s   (threshold: 90%%)", greenStyle.Render("3.0%")))
	fmt.Println()
	fmt.Println(okStyle.Render("  ✓ All systems healthy"))
}

func printTUI() {
	var b strings.Builder

	// Tabs
	tab1 := activeTabStyle.Render(" [1] homelab-server ")
	tab2 := inactiveTabStyle.Render(" [2] rpi5 ")
	serverCount := dimStyle.Render("  (2 available · Tab to switch)")
	b.WriteString(tab1 + tab2 + serverCount + "\n")

	// Left panel — System
	leftWidth := tuiWidth/3 + 2
	barW := leftWidth - 20
	if barW < 8 {
		barW = 8
	}
	var sysLines []string
	sysLines = append(sysLines, titleStyle.Render("⚡ homelab-server"))
	sysLines = append(sysLines, "")
	sysLines = append(sysLines, fmt.Sprintf("  CPU  %s %5.1f%%", progressBar(23.4, barW), 23.4))
	sysLines = append(sysLines, fmt.Sprintf("  Mem  %s %5.1f%%", progressBar(49.2, barW), 49.2))
	sysLines = append(sysLines, fmt.Sprintf("  /    %s %5.1f%%", progressBar(3.0, barW), 3.0))
	sysLines = append(sysLines, "")
	sysLines = append(sysLines, "  Uptime:  42d 7h")
	sysLines = append(sysLines, "  OS:      linux/arm64")
	sysLines = append(sysLines, "  Cores:   10")
	sysLines = append(sysLines, "  Memory:  7.9 / 16.0 GB")
	leftPanel := panelStyle.Width(leftWidth).Render(strings.Join(sysLines, "\n"))

	// Right panel — Docker
	rightWidth := tuiWidth - leftWidth - 4
	var dockLines []string
	dockLines = append(dockLines, titleStyle.Render("Docker Containers"))
	dockLines = append(dockLines, "")
	hdr := fmt.Sprintf("  %-18s %-10s %-10s %s", "NAME", "STATE", "IMAGE", "STATUS")
	dockLines = append(dockLines, headerStyle.Render(hdr))

	type ctr struct{ name, state, image, status string }
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
	footer := panelStyle.Width(tuiWidth - 4).Render(strings.Join(footerParts, "\n"))
	b.WriteString(footer)

	fmt.Print(b.String())
}
