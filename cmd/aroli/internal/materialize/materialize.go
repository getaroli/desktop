// Package materialize applies the delivery contract for the versioned home
// tree while preserving files that are owned by the person using the desktop.
package materialize

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/getaroli/desktop/cmd/aroli/internal/sys"
)

type state struct {
	Files     map[string]string `json:"files"`
	UpdatedAt string            `json:"updated_at"`
}

// Materialize runs the proven config phase, safely retires files removed from
// home/, then overlays explicit user edits. Seeds remain the installer's
// responsibility and are therefore created only when absent.
func Materialize(args []string) error {
	fs := flag.NewFlagSet("materialize", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	dry := fs.Bool("dry-run", false, "mostra o plano sem alterar arquivos")
	repoFlag := fs.String("repo", "", "checkout do aroli")
	if err := fs.Parse(args); err != nil || fs.NArg() != 0 {
		if err != nil {
			return err
		}
		return errors.New("use: aroli materialize [--dry-run] [--repo CAMINHO]")
	}
	repo, cleanup, err := sys.EnsureRepositoryWithCleanup(*repoFlag, *dry)
	if err != nil {
		return err
	}
	defer cleanup()
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	old, _ := readState(home)
	if *dry {
		fmt.Println("Materialize: reaplicaria a base de home/, retiraria somente arquivos ainda gerenciados e aplicaria user_edits.")
	}
	backend := []string{"install.sh", "--yes"}
	if *dry {
		backend = append(backend, "--dry-run")
	}
	backend = append(backend, "config")
	if err := sys.Command(repo, "bash", backend...).Run(); err != nil {
		return err
	}
	next, err := inventory(filepath.Join(repo, "home"))
	if err != nil {
		return err
	}
	if err := pruneRemoved(home, repo, old.Files, next, *dry); err != nil {
		return err
	}
	if err := overlay(filepath.Join(home, ".config", "aroli", "user_edits"), home, *dry); err != nil {
		return err
	}
	if *dry {
		return nil
	}
	return writeState(home, state{Files: next, UpdatedAt: time.Now().Format(time.RFC3339)})
}

func statePath(home string) string {
	return filepath.Join(home, ".local", "state", "aroli-desktop", "materialize.json")
}

// StatePresent reports whether this machine has completed a materialization.
func StatePresent(home string) bool {
	_, err := os.Stat(statePath(home))
	return err == nil
}
func readState(home string) (state, error) {
	var s state
	b, err := os.ReadFile(statePath(home))
	if os.IsNotExist(err) {
		return state{Files: map[string]string{}}, nil
	}
	if err != nil {
		return s, err
	}
	err = json.Unmarshal(b, &s)
	if s.Files == nil {
		s.Files = map[string]string{}
	}
	return s, err
}
func writeState(home string, s state) error {
	if err := os.MkdirAll(filepath.Dir(statePath(home)), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(statePath(home), b, 0o600)
}

func inventory(root string) (map[string]string, error) {
	out := map[string]string{}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		sum, err := digest(path)
		if err != nil {
			return err
		}
		out[filepath.ToSlash(rel)] = sum
		return nil
	})
	return out, err
}
func digest(path string) (string, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return "", err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		v, err := os.Readlink(path)
		if err != nil {
			return "", err
		}
		sum := sha256.Sum256([]byte("link:" + v))
		return hex.EncodeToString(sum[:]), nil
	}
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	_, err = io.Copy(h, f)
	return hex.EncodeToString(h.Sum(nil)), err
}

func protected(rel string) bool {
	return strings.HasPrefix(rel, ".config/hypr/user.") || rel == ".config/mimeapps.list" || strings.Contains(rel, "/state/")
}
func pruneRemoved(home, repo string, old, next map[string]string, dry bool) error {
	keys := make([]string, 0)
	for rel := range old {
		if _, ok := next[rel]; !ok {
			keys = append(keys, rel)
		}
	}
	sort.Strings(keys)
	backup := filepath.Join(home, ".dotfiles-backup", "materialize-prune-"+time.Now().Format("20060102-150405"))
	for _, rel := range keys {
		if protected(rel) {
			continue
		}
		dest := filepath.Join(home, filepath.FromSlash(rel))
		got, err := digest(dest)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		owned := got == old[rel]
		if link, err := os.Readlink(dest); err == nil {
			owned = strings.HasPrefix(filepath.Clean(link), filepath.Clean(repo)+string(os.PathSeparator))
		}
		if !owned {
			fmt.Printf("Materialize: mantido (modificado por você): %s\n", dest)
			continue
		}
		if dry {
			fmt.Printf("Materialize: moveria removido para backup: %s\n", dest)
			continue
		}
		if err := os.MkdirAll(filepath.Join(backup, filepath.Dir(rel)), 0o755); err != nil {
			return err
		}
		if err := os.Rename(dest, filepath.Join(backup, rel)); err != nil {
			return err
		}
	}
	return nil
}

func overlay(root, home string, dry bool) error {
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		dest := filepath.Join(home, rel)
		if dry {
			fmt.Printf("Materialize: aplicaria user_edit: %s\n", dest)
			return nil
		}
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return err
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()
		tmp, err := os.CreateTemp(filepath.Dir(dest), ".aroli-edit-")
		if err != nil {
			return err
		}
		name := tmp.Name()
		defer os.Remove(name)
		if _, err = io.Copy(tmp, in); err == nil {
			err = tmp.Chmod(0o644)
		}
		if closeErr := tmp.Close(); err == nil {
			err = closeErr
		}
		if err != nil {
			return err
		}
		return os.Rename(name, dest)
	})
}
