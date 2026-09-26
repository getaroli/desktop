package sys

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidRepo(t *testing.T) {
	directory := t.TempDir()
	if _, err := ValidRepo(directory); err == nil {
		t.Fatal("directory without installer was accepted")
	}
	if err := os.WriteFile(filepath.Join(directory, "install.sh"), []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	got, err := ValidRepo(directory)
	if err != nil {
		t.Fatal(err)
	}
	if got != directory {
		t.Fatalf("got %q, want %q", got, directory)
	}
}

func TestEnsureRepositoryDryRunBootstrapsTemporaryCheckout(t *testing.T) {
	originalClone := CloneRepository
	originalQuietClone := CloneRepositoryQuiet
	t.Cleanup(func() { CloneRepository = originalClone; CloneRepositoryQuiet = originalQuietClone })
	CloneRepository = func(target string) error {
		if err := os.MkdirAll(target, 0o755); err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(target, "install.sh"), []byte("#!/bin/sh\n"), 0o755)
	}
	CloneRepositoryQuiet = CloneRepository

	workingDirectory := t.TempDir()
	originalDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(workingDirectory); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(originalDirectory) })
	// Keep a real local installation from changing what this isolated test
	// exercises: it must reach the temporary dry-run checkout branch.
	t.Setenv("HOME", t.TempDir())

	path, cleanup, err := EnsureRepositoryWithCleanup("", true)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(cleanup)
	if _, err := ValidRepo(path); err != nil {
		t.Fatalf("dry-run checkout is invalid: %v", err)
	}
	if !strings.HasPrefix(filepath.Base(filepath.Dir(path)), "aroli-dry-run-") {
		t.Fatalf("dry-run checkout was not temporary: %q", path)
	}
	cleanup()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("temporary checkout still exists: %v", err)
	}
}

func TestCopyTreePreservesFilesAndLinks(t *testing.T) {
	source, target := t.TempDir(), filepath.Join(t.TempDir(), "snapshot")
	if err := os.WriteFile(filepath.Join(source, "theme.conf"), []byte("umbra"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(source, "nested"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "nested", "colors"), []byte("violet"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("theme.conf", filepath.Join(source, "current")); err != nil {
		t.Fatal(err)
	}
	if err := CopyTree(source, target); err != nil {
		t.Fatal(err)
	}
	if data, err := os.ReadFile(filepath.Join(target, "nested", "colors")); err != nil || string(data) != "violet" {
		t.Fatalf("copied data = %q, %v", data, err)
	}
	if link, err := os.Readlink(filepath.Join(target, "current")); err != nil || link != "theme.conf" {
		t.Fatalf("copied link = %q, %v", link, err)
	}
}
