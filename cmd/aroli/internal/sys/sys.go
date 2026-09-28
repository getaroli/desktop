// Package sys holds the shared primitives every aroli command builds on:
// running processes, user confirmations, file trees, and checkout lookup.
// It depends only on the standard library.
package sys

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	Project       = "Aroli Desktop"
	RepositoryURL = "https://github.com/eduardoaugustolb/aroli-desktop.git"
	DefaultLang   = "pt-BR"
)

// Command runs a process attached to the user's terminal in dir.
func Command(dir, name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd
}

// Confirm asks a [s/N] question on reader.
func Confirm(reader *bufio.Reader, question string) bool {
	fmt.Printf("%s [s/N] ", question)
	answer, _ := reader.ReadString('\n')
	answer = strings.ToLower(strings.TrimSpace(answer))
	return answer == "s" || answer == "sim" || answer == "y" || answer == "yes"
}

// CopyTree copies files, directories, and symlinks from src to dst.
func CopyTree(src, dst string) error {
	return filepath.WalkDir(src, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := dst
		if rel != "." {
			target = filepath.Join(dst, rel)
		}
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode().Perm())
		}
		if info.Mode()&os.ModeSymlink != 0 {
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			return os.Symlink(link, target)
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		out, err := os.OpenFile(target, os.O_CREATE|os.O_EXCL|os.O_WRONLY, info.Mode().Perm())
		if err != nil {
			return err
		}
		_, err = io.Copy(out, in)
		closeErr := out.Close()
		if err != nil {
			return err
		}
		return closeErr
	})
}

// ValidRepo reports path when it looks like an Aroli Desktop checkout.
func ValidRepo(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	if info, err := os.Stat(filepath.Join(path, "install.sh")); err != nil || info.IsDir() {
		return "", fmt.Errorf("%s não parece ser um checkout do Aroli Desktop", path)
	}
	return path, nil
}

// CloneRepository and CloneRepositoryQuiet are variables so tests can stub
// the network out.
var CloneRepository = func(target string) error {
	return Command("", "git", "clone", "--depth", "1", RepositoryURL, target).Run()
}

var CloneRepositoryQuiet = func(target string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", "clone", "--depth", "1", RepositoryURL, target)
	cmd.Stdin = os.Stdin
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	return cmd.Run()
}

// EnsureRepository makes a CLI installation useful even when the user has
// never cloned the dotfiles repository. Existing checkouts always win.
func EnsureRepository(explicit string, dryRun bool) (string, error) {
	path, _, err := EnsureRepositoryWithCleanup(explicit, dryRun)
	return path, err
}

func EnsureRepositoryWithCleanup(explicit string, dryRun bool) (string, func(), error) {
	noCleanup := func() {}
	if explicit != "" {
		path, err := ValidRepo(explicit)
		return path, noCleanup, err
	}
	if env := os.Getenv("AROLI_DESKTOP_REPO"); env != "" {
		if path, err := ValidRepo(env); err == nil {
			return path, noCleanup, nil
		}
	}
	if cwd, err := os.Getwd(); err == nil {
		if path, err := ValidRepo(cwd); err == nil {
			return path, noCleanup, nil
		}
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", noCleanup, err
	}
	target := filepath.Join(home, ".local", "share", "aroli-desktop")
	if path, err := ValidRepo(target); err == nil {
		return path, noCleanup, nil
	}
	if dryRun {
		parent, err := os.MkdirTemp("", "aroli-dry-run-")
		if err != nil {
			return "", noCleanup, err
		}
		dryRunTarget := filepath.Join(parent, "checkout")
		if err := CloneRepositoryQuiet(dryRunTarget); err != nil {
			_ = os.RemoveAll(parent)
			return "", noCleanup, fmt.Errorf("não foi possível preparar o aroli para o dry-run: %w", err)
		}
		return dryRunTarget, func() { _ = os.RemoveAll(parent) }, nil
	}
	if _, err := exec.LookPath("git"); err != nil {
		return "", noCleanup, errors.New("git é necessário para baixar o aroli")
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return "", noCleanup, err
	}
	fmt.Println("Baixando o aroli pela primeira vez…")
	if err := CloneRepository(target); err != nil {
		return "", noCleanup, fmt.Errorf("não foi possível baixar o aroli: %w", err)
	}
	path, err := ValidRepo(target)
	return path, noCleanup, err
}

// RunBackend delegates to the shell backend: the read-only diagnose script
// or the version/update/rollback/prune dispatcher. The checkout is prepared
// for dry-run so these commands never clone on the user's behalf.
func RunBackend(kind string, args []string) error {
	repo, cleanup, err := EnsureRepositoryWithCleanup("", true)
	if err != nil {
		return err
	}
	defer cleanup()
	if _, err := ValidRepo(repo); err != nil {
		return errors.New("aroli ainda não foi instalado; execute: aroli install")
	}
	if kind == "diagnose" {
		return Command(repo, "bash", append([]string{"diagnose"}, args...)...).Run()
	}
	legacy := filepath.Join(repo, "home", ".local", "bin", "aroli-backend")
	if _, err := os.Stat(legacy); err != nil {
		return fmt.Errorf("backend de atualização ausente: %w", err)
	}
	return Command(repo, "bash", append([]string{legacy}, args...)...).Run()
}

// UserScript resolves a helper installed in ~/.local/bin.
func UserScript(name string) (string, error) {
	home, err := os.UserHomeDir()
	return filepath.Join(home, ".local", "bin", name), err
}

// ScriptState runs a helper with one action and returns its trimmed output.
func ScriptState(script, action string) (string, error) {
	path, err := UserScript(script)
	if err != nil {
		return "", err
	}
	out, err := exec.Command(path, action).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
