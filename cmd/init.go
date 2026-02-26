package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Higangssh/homebutler/internal/config"
	"gopkg.in/yaml.v3"
)

func runInit() error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("🏠 HomeButler Setup")
	fmt.Println()

	// Determine config path
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot determine home directory: %w", err)
	}
	cfgPath := filepath.Join(home, ".config", "homebutler", "config.yaml")

	// Check if config already exists
	if _, err := os.Stat(cfgPath); err == nil {
		fmt.Print("Config already exists. Overwrite? (y/n): ")
		if !scanLine(scanner) {
			return nil
		}
		if !isYes(scanner.Text()) {
			fmt.Println("Aborted.")
			return nil
		}
	}

	cfg := &config.Config{
		Alerts: config.AlertConfig{CPU: 90, Memory: 85, Disk: 90},
		Output: "json",
	}

	for {
		server, err := promptServer(scanner)
		if err != nil {
			return err
		}
		cfg.Servers = append(cfg.Servers, *server)

		fmt.Printf("✅ Added %s (%s)\n", server.Name, server.Host)
		fmt.Println()

		fmt.Print("Add another server? (y/n): ")
		if !scanLine(scanner) {
			break
		}
		if !isYes(scanner.Text()) {
			break
		}
		fmt.Println()
	}

	// Marshal and save
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	if err := os.WriteFile(cfgPath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	fmt.Println()
	fmt.Printf("Config saved to %s\n", cfgPath)
	fmt.Println("Run homebutler status to test!")
	return nil
}

func promptServer(scanner *bufio.Scanner) (*config.ServerConfig, error) {
	server := &config.ServerConfig{}

	// Name (required)
	name, err := promptRequired(scanner, "Server name: ")
	if err != nil {
		return nil, err
	}
	server.Name = name

	// Host (required)
	host, err := promptRequired(scanner, "Host/IP: ")
	if err != nil {
		return nil, err
	}
	server.Host = host

	// Local?
	fmt.Print("Is this the local machine? (y/n): ")
	if !scanLine(scanner) {
		return nil, fmt.Errorf("unexpected end of input")
	}
	if isYes(scanner.Text()) {
		server.Local = true
		return server, nil
	}

	// Remote — SSH details
	user, err := promptRequired(scanner, "SSH user: ")
	if err != nil {
		return nil, err
	}
	server.User = user

	fmt.Print("SSH port (22): ")
	if !scanLine(scanner) {
		return nil, fmt.Errorf("unexpected end of input")
	}
	portStr := strings.TrimSpace(scanner.Text())
	if portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid port: %s", portStr)
		}
		server.Port = port
	}

	fmt.Print("Auth method (key/password): ")
	if !scanLine(scanner) {
		return nil, fmt.Errorf("unexpected end of input")
	}
	authMethod := strings.TrimSpace(strings.ToLower(scanner.Text()))
	if authMethod == "password" {
		server.AuthMode = "password"
		password, err := promptRequired(scanner, "Password: ")
		if err != nil {
			return nil, err
		}
		server.Password = password
	} else {
		fmt.Print("Key file (~/.ssh/id_rsa): ")
		if !scanLine(scanner) {
			return nil, fmt.Errorf("unexpected end of input")
		}
		keyFile := strings.TrimSpace(scanner.Text())
		if keyFile == "" {
			keyFile = "~/.ssh/id_rsa"
		}
		server.KeyFile = keyFile
	}

	return server, nil
}

func promptRequired(scanner *bufio.Scanner, prompt string) (string, error) {
	for {
		fmt.Print(prompt)
		if !scanLine(scanner) {
			return "", fmt.Errorf("unexpected end of input")
		}
		val := strings.TrimSpace(scanner.Text())
		if val != "" {
			return val, nil
		}
		fmt.Println("  This field is required.")
	}
}

func scanLine(scanner *bufio.Scanner) bool {
	return scanner.Scan()
}

func isYes(s string) bool {
	return strings.TrimSpace(strings.ToLower(s)) == "y"
}
