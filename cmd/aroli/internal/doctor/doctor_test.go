package doctor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDoctorFixSeedsMissingLocalFilesAndPreservesUserChoice(t *testing.T) {
	home := t.TempDir()
	repo := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("AROLI_DESKTOP_REPO", repo)
	if err := os.WriteFile(filepath.Join(repo, "install.sh"), []byte("#!/usr/bin/env bash\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(repo, "seeds"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "seeds", "mimeapps.list"), []byte("[Default Applications]\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	custom := filepath.Join(home, ".zshrc.local")
	if err := os.WriteFile(custom, []byte("export PATH=custom\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := Doctor([]string{"--fix"}); err != nil {
		t.Fatal(err)
	}
	if err := Doctor([]string{"--fix"}); err != nil {
		t.Fatal(err)
	}

	if got, err := os.ReadFile(custom); err != nil || string(got) != "export PATH=custom\n" {
		t.Fatalf("user shell override changed to %q, %v", got, err)
	}
	for _, path := range []string{
		filepath.Join(home, ".config", "aroli", "user_edits"),
		filepath.Join(home, ".config", "mimeapps.list"),
		filepath.Join(home, ".bashrc.local"),
		filepath.Join(home, ".zprofile.local"),
		filepath.Join(home, ".profile.local"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Errorf("doctor did not seed %s: %v", path, err)
		}
	}
}
