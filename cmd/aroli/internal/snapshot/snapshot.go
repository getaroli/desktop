// Package snapshot saves and restores the supported visual configuration
// paths locally. Snapshots never include credentials or personal documents.
package snapshot

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/getaroli/desktop/cmd/aroli/internal/sys"
)

var snapshotPaths = []string{
	".config/hypr", ".config/quickshell", ".config/kitty", ".config/gtk-3.0", ".config/gtk-4.0", ".config/fastfetch", ".config/yazi",
}

// SnapshotRoot is the local snapshot store.
func SnapshotRoot() (string, error) {
	h, err := os.UserHomeDir()
	return filepath.Join(h, ".local", "share", "aroli-desktop", "snapshots"), err
}

func validSnapshotName(name string) bool {
	return name != "" && filepath.Base(name) == name && name != "." && name != ".."
}

// Snapshots implements `aroli snapshot create|list|restore NOME`.
func Snapshots(args []string) error {
	if len(args) == 0 {
		return errors.New("use: aroli snapshot create|list|restore NOME")
	}
	root, err := SnapshotRoot()
	if err != nil {
		return err
	}
	switch args[0] {
	case "list":
		entries, err := os.ReadDir(root)
		if os.IsNotExist(err) {
			fmt.Println("Nenhum snapshot local.")
			return nil
		}
		if err != nil {
			return err
		}
		for _, e := range entries {
			if e.IsDir() {
				fmt.Println(e.Name())
			}
		}
		return nil
	case "create":
		if len(args) != 2 || !validSnapshotName(args[1]) {
			return errors.New("use: aroli snapshot create NOME (somente um nome simples)")
		}
		dest := filepath.Join(root, args[1])
		if _, err := os.Stat(dest); err == nil {
			return fmt.Errorf("o snapshot %q já existe", args[1])
		}
		if err := os.MkdirAll(dest, 0o755); err != nil {
			return err
		}
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		for _, rel := range snapshotPaths {
			src := filepath.Join(home, rel)
			if _, err := os.Lstat(src); err == nil {
				if err := sys.CopyTree(src, filepath.Join(dest, rel)); err != nil {
					return err
				}
			}
		}
		return os.WriteFile(filepath.Join(dest, "METADATA"), []byte("created="+time.Now().Format(time.RFC3339)+"\n"), 0o644)
	case "restore":
		if len(args) < 2 || !validSnapshotName(args[1]) {
			return errors.New("use: aroli snapshot restore NOME [--yes]")
		}
		assume := len(args) == 3 && args[2] == "--yes"
		src := filepath.Join(root, args[1])
		if _, err := os.Stat(src); err != nil {
			return fmt.Errorf("snapshot %q não existe", args[1])
		}
		if !assume && !sys.Confirm(bufio.NewReader(os.Stdin), "Restaurar este snapshot sobre as configurações atuais?") {
			fmt.Println("Nada foi alterado.")
			return nil
		}
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		for _, rel := range snapshotPaths {
			from := filepath.Join(src, rel)
			if _, err := os.Lstat(from); err == nil {
				target := filepath.Join(home, rel)
				backup := target + ".before-snapshot-" + time.Now().Format("20060102-150405")
				if _, err := os.Lstat(target); err == nil {
					if err := os.Rename(target, backup); err != nil {
						return err
					}
				}
				if err := sys.CopyTree(from, target); err != nil {
					return err
				}
			}
		}
		fmt.Println("Snapshot restaurado. As versões anteriores foram mantidas com .before-snapshot-…")
		return nil
	default:
		return errors.New("use: aroli snapshot create|list|restore NOME")
	}
}
