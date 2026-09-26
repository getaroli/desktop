package prefs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExportPreferencesCopiesOnlyPortableSettings(t *testing.T) {
	home, destination := t.TempDir(), filepath.Join(t.TempDir(), "preferences")
	t.Setenv("HOME", home)
	if err := os.MkdirAll(filepath.Join(home, ".config", "hypr"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, ".config", "quickshell-rice.json"), []byte(`{"language":"pt-BR"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, ".config", "hypr", "language.conf"), []byte("pt-BR"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := ExportPreferences([]string{destination}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(destination, "aroli-preferences.json")); err != nil {
		t.Fatalf("manifest not exported: %v", err)
	}
	if data, err := os.ReadFile(filepath.Join(destination, ".config", "hypr", "language.conf")); err != nil || string(data) != "pt-BR" {
		t.Fatalf("language export = %q, %v", data, err)
	}
}
