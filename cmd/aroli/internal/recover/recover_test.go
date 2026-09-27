package recover

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/getaroli/desktop/cmd/aroli/internal/snapshot"
)

func TestRecoverDryRunWithLatestSnapshotDoesNotChangeConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	config := filepath.Join(home, ".config", "hypr")
	if err := os.MkdirAll(config, 0o755); err != nil {
		t.Fatal(err)
	}
	configFile := filepath.Join(config, "hyprland.lua")
	if err := os.WriteFile(configFile, []byte("stable state\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := snapshot.Snapshots([]string{"create", "recovery-check"}); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configFile, []byte("current state\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := RecoverDesktop([]string{"--dry-run", "--snapshot", "latest"}); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(configFile)
	if err != nil || string(got) != "current state\n" {
		t.Fatalf("dry-run changed config to %q, %v", got, err)
	}
}
