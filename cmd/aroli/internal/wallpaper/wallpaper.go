// Package wallpaper lists, applies, imports, and removes wallpapers.
// Applying one also derives its dynamic palette.
package wallpaper

import (
	"bufio"
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/getaroli/desktop/cmd/aroli/internal/sys"
)

// Wallpapers implements `aroli wallpaper list|set|random|import|remove`.
func Wallpapers(args []string) error {
	if len(args) == 0 {
		return errors.New("use: aroli wallpaper list|set|random|import|remove")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dir := filepath.Join(home, "Pictures", "wallpapers")
	entries, _ := os.ReadDir(dir)
	images := []string{}
	for _, e := range entries {
		if !e.IsDir() {
			ext := strings.ToLower(filepath.Ext(e.Name()))
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp" {
				images = append(images, e.Name())
			}
		}
	}
	sort.Strings(images)
	switch args[0] {
	case "list":
		for _, name := range images {
			fmt.Println(name)
		}
		return nil
	case "set":
		if len(args) != 2 {
			return errors.New("use: aroli wallpaper set ARQUIVO")
		}
		return SetWallpaper(filepath.Join(dir, filepath.Base(args[1])))
	case "random":
		if len(images) == 0 {
			return errors.New("não há wallpapers instalados")
		}
		var b [1]byte
		if _, err := rand.Read(b[:]); err != nil {
			return err
		}
		return SetWallpaper(filepath.Join(dir, images[int(b[0])%len(images)]))
	case "import":
		if len(args) != 2 {
			return errors.New("use: aroli wallpaper import /caminho/imagem")
		}
		source, err := filepath.Abs(args[1])
		if err != nil {
			return err
		}
		if _, err := os.Stat(source); err != nil {
			return err
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		dest := filepath.Join(dir, filepath.Base(source))
		if _, err := os.Stat(dest); err == nil {
			return fmt.Errorf("%s já existe", filepath.Base(source))
		}
		return sys.CopyTree(source, dest)
	case "remove":
		if len(args) < 2 {
			return errors.New("use: aroli wallpaper remove ARQUIVO [--yes]")
		}
		path := filepath.Join(dir, filepath.Base(args[1]))
		if _, err := os.Stat(path); err != nil {
			return err
		}
		if len(args) < 3 || args[2] != "--yes" {
			if !sys.Confirm(bufio.NewReader(os.Stdin), "Mover este wallpaper para a Lixeira?") {
				return nil
			}
		}
		if _, err := exec.LookPath("gio"); err == nil {
			return sys.Command("", "gio", "trash", path).Run()
		}
		return os.Rename(path, path+".removed-"+time.Now().Format("20060102-150405"))
	default:
		return errors.New("use: aroli wallpaper list|set|random|import|remove")
	}
}

// SetWallpaper applies a wallpaper file through the session helper.
func SetWallpaper(path string) error {
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("wallpaper não encontrado: %w", err)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	return sys.Command("", filepath.Join(home, ".config", "hypr", "set-wallpaper.sh"), path).Run()
}
