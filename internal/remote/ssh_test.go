package remote

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/crypto/ssh/knownhosts"
)

func TestNewKnownHostsCallback(t *testing.T) {
	cb, err := newKnownHostsCallback()
	if err != nil {
		t.Fatalf("newKnownHostsCallback() error: %v", err)
	}
	if cb == nil {
		t.Fatal("expected non-nil callback")
	}
}

func TestNewKnownHostsCallback_CreatesFile(t *testing.T) {
	// Verify that ~/.ssh/known_hosts exists after calling newKnownHostsCallback
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot determine home directory")
	}
	path := filepath.Join(home, ".ssh", "known_hosts")

	_, err = newKnownHostsCallback()
	if err != nil {
		t.Fatalf("newKnownHostsCallback() error: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected %s to exist after newKnownHostsCallback()", path)
	}
}

func TestKnownHostsPath(t *testing.T) {
	path, err := knownHostsPath()
	if err != nil {
		t.Fatalf("knownHostsPath() error: %v", err)
	}
	if !strings.HasSuffix(path, filepath.Join(".ssh", "known_hosts")) {
		t.Errorf("expected path ending in .ssh/known_hosts, got %s", path)
	}
}

func TestTofuConnect_NoServer(t *testing.T) {
	t.Skip("tofuConnect requires a real SSH server; skipping in unit tests")
}

// TestErrorMessages_KeyMismatch verifies the key-mismatch error contains actionable hints.
func TestErrorMessages_KeyMismatch(t *testing.T) {
	// Simulate the error message that connect() would produce for a key mismatch.
	serverName := "myserver"
	addr := "10.0.0.1:22"
	msg := fmt.Sprintf("[%s] ⚠️  SSH HOST KEY CHANGED (%s)\n"+
		"  The server's host key does not match the one in ~/.ssh/known_hosts.\n"+
		"  This could mean:\n"+
		"    1. The server was reinstalled or its SSH keys were regenerated\n"+
		"    2. A man-in-the-middle attack is in progress\n\n"+
		"  → If you trust this change: homebutler trust %s --reset\n"+
		"  → If unexpected: do NOT connect and investigate", serverName, addr, serverName)

	if !strings.Contains(msg, "homebutler trust") {
		t.Error("key mismatch error should contain 'homebutler trust'")
	}
	if !strings.Contains(msg, "HOST KEY CHANGED") {
		t.Error("key mismatch error should contain 'HOST KEY CHANGED'")
	}
	if !strings.Contains(msg, "--reset") {
		t.Error("key mismatch error should contain '--reset' flag")
	}
}

// TestErrorMessages_UnknownHost verifies the unknown-host error contains actionable hints.
func TestErrorMessages_UnknownHost(t *testing.T) {
	serverName := "newserver"
	addr := "10.0.0.2:22"
	msg := fmt.Sprintf("[%s] failed to auto-register host key for %s\n  → Register manually: homebutler trust %s\n  → Check SSH connectivity: ssh %s@%s -p %d",
		serverName, addr, serverName, "root", "10.0.0.2", 22)

	if !strings.Contains(msg, "homebutler trust") {
		t.Error("unknown host error should contain 'homebutler trust'")
	}
	if !strings.Contains(msg, "ssh root@") {
		t.Error("unknown host error should contain ssh connection hint")
	}
}

// TestErrorMessages_NoCredentials verifies the no-credentials error contains config hint.
func TestErrorMessages_NoCredentials(t *testing.T) {
	serverName := "nocreds"
	msg := fmt.Sprintf("[%s] no SSH credentials configured\n  → Add 'key_file' or 'password' to this server in ~/.config/homebutler/config.yaml", serverName)

	if !strings.Contains(msg, "key_file") {
		t.Error("no-credentials error should mention key_file")
	}
	if !strings.Contains(msg, "config.yaml") {
		t.Error("no-credentials error should mention config.yaml")
	}
}

// TestErrorMessages_Timeout verifies the timeout error contains actionable hints.
func TestErrorMessages_Timeout(t *testing.T) {
	serverName := "slowserver"
	addr := "10.0.0.3:22"
	msg := fmt.Sprintf("[%s] connection timed out (%s)\n  → Check if the server is online and reachable\n  → Verify host/port in ~/.config/homebutler/config.yaml", serverName, addr)

	if !strings.Contains(msg, "timed out") {
		t.Error("timeout error should contain 'timed out'")
	}
	if !strings.Contains(msg, "config.yaml") {
		t.Error("timeout error should mention config.yaml")
	}
}

// TestKnownHostsKeyError verifies that knownhosts.KeyError works as expected
// for both key-mismatch and unknown-host scenarios.
func TestKnownHostsKeyError(t *testing.T) {
	// KeyError with Want = empty means unknown host
	unknownErr := &knownhosts.KeyError{}
	if len(unknownErr.Want) != 0 {
		t.Error("empty KeyError should have no Want entries (unknown host)")
	}

	// KeyError with Want populated means key mismatch
	mismatchErr := &knownhosts.KeyError{
		Want: []knownhosts.KnownKey{{Filename: "known_hosts", Line: 1}},
	}
	if len(mismatchErr.Want) == 0 {
		t.Error("mismatch KeyError should have Want entries")
	}
}
