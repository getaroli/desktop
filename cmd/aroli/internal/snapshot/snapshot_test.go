package snapshot

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateRestoreKeepsBackupAndExcludesUnmanagedData(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	config := filepath.Join(home, ".config", "hypr")
	if err := os.MkdirAll(config, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(config, "hyprland.lua"), []byte("before\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	secret := filepath.Join(home, ".ssh", "private-key")
	if err := os.MkdirAll(filepath.Dir(secret), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(secret, []byte("secret"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := Snapshots([]string{"create", "before-change"}); err != nil {
		t.Fatal(err)
	}
	stored := filepath.Join(home, ".local", "share", "aroli-desktop", "snapshots", "before-change")
	if _, err := os.Stat(filepath.Join(stored, ".config", "hypr", "hyprland.lua")); err != nil {
		t.Fatalf("snapshot did not include supported config: %v", err)
	}
	if _, err := os.Stat(filepath.Join(stored, ".ssh")); !os.IsNotExist(err) {
		t.Fatalf("snapshot included private data: %v", err)
	}

	if err := os.WriteFile(filepath.Join(config, "hyprland.lua"), []byte("after\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := Snapshots([]string{"restore", "before-change", "--yes"}); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(config, "hyprland.lua"))
	if err != nil || string(got) != "before\n" {
		t.Fatalf("restored config = %q, %v", got, err)
	}
	backups, err := filepath.Glob(config + ".before-snapshot-*")
	if err != nil || len(backups) != 1 {
		t.Fatalf("backup paths = %v, %v; want one preserved backup", backups, err)
	}
	backup, err := os.ReadFile(filepath.Join(backups[0], "hyprland.lua"))
	if err != nil || string(backup) != "after\n" {
		t.Fatalf("preserved backup = %q, %v", backup, err)
	}
}

func TestSnapshotRejectsPathLikeNames(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	for _, name := range []string{"../escape", ".", "..", ""} {
		if err := Snapshots([]string{"create", name}); err == nil {
			t.Errorf("Snapshots(create, %q) accepted a path-like name", name)
		}
	}
}
