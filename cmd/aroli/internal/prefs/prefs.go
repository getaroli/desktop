// Package prefs moves portable visual preferences between machines.
// Exports carry no credentials or personal data; optionals are listed,
// never installed automatically.
package prefs

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/getaroli/desktop/cmd/aroli/internal/profiles"
	"github.com/getaroli/desktop/cmd/aroli/internal/sys"
	"github.com/getaroli/desktop/cmd/aroli/internal/wallpaper"
)

type portablePreferences struct {
	Version         int      `json:"version"`
	ExportedAt      string   `json:"exported_at"`
	Language        string   `json:"language,omitempty"`
	Wallpaper       string   `json:"wallpaper,omitempty"`
	Session         string   `json:"session,omitempty"`
	OptionalPlugins []string `json:"optional_plugins,omitempty"`
}

// ExportPreferences implements `aroli export DIRETÓRIO`.
func ExportPreferences(args []string) error {
	if len(args) != 1 {
		return errors.New("use: aroli export DIRETÓRIO")
	}
	dest, err := filepath.Abs(args[0])
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return err
	}
	if _, err := os.Stat(filepath.Join(dest, "aroli-preferences.json")); err == nil {
		return errors.New("o diretório já contém um export do Aroli Desktop; escolha outro destino")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	for _, rel := range []string{".config/quickshell-aroli.json", ".config/hypr/language.conf"} {
		source := filepath.Join(home, rel)
		if _, err := os.Lstat(source); err == nil {
			if err := sys.CopyTree(source, filepath.Join(dest, rel)); err != nil {
				return err
			}
		}
	}
	pref := portablePreferences{Version: 1, ExportedAt: time.Now().Format(time.RFC3339)}
	if data, err := os.ReadFile(filepath.Join(home, ".config", "hypr", "language.conf")); err == nil {
		pref.Language = strings.TrimSpace(string(data))
	}
	if target, err := os.Readlink(filepath.Join(home, ".cache", "wal", "lockbg")); err == nil {
		pref.Wallpaper = filepath.Base(target)
	}
	for _, p := range profiles.AroliProfiles {
		for _, name := range p.Plugins {
			if _, err := exec.Command("pacman", "-Q", name).Output(); err == nil {
				pref.OptionalPlugins = append(pref.OptionalPlugins, name)
			}
		}
	}
	data, err := json.MarshalIndent(pref, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dest, "aroli-preferences.json"), data, 0o600); err != nil {
		return err
	}
	fmt.Printf("Preferências exportadas para %s\n", dest)
	return nil
}

// ImportPreferences implements `aroli import DIRETÓRIO [--yes]`.
func ImportPreferences(args []string) error {
	if len(args) < 1 || len(args) > 2 || (len(args) == 2 && args[1] != "--yes") {
		return errors.New("use: aroli import DIRETÓRIO [--yes]")
	}
	source, err := filepath.Abs(args[0])
	if err != nil {
		return err
	}
	data, err := os.ReadFile(filepath.Join(source, "aroli-preferences.json"))
	if err != nil {
		return err
	}
	var pref portablePreferences
	if err := json.Unmarshal(data, &pref); err != nil || pref.Version != 1 {
		return errors.New("arquivo de preferências inválido ou incompatível")
	}
	if len(args) == 1 && !sys.Confirm(bufio.NewReader(os.Stdin), "Importar preferências visuais deste diretório?") {
		return nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	for _, rel := range []string{".config/quickshell-aroli.json", ".config/hypr/language.conf"} {
		from := filepath.Join(source, rel)
		if _, err := os.Lstat(from); err == nil {
			target := filepath.Join(home, rel)
			if _, err := os.Lstat(target); err == nil {
				backup := target + ".before-import-" + time.Now().Format("20060102-150405")
				if err := os.Rename(target, backup); err != nil {
					return err
				}
			}
			if err := sys.CopyTree(from, target); err != nil {
				return err
			}
		}
	}
	if pref.Wallpaper != "" {
		wall := filepath.Join(home, "Pictures", "wallpapers", filepath.Base(pref.Wallpaper))
		if _, err := os.Stat(wall); err == nil {
			_ = wallpaper.SetWallpaper(wall)
		}
	}
	if len(pref.OptionalPlugins) > 0 {
		fmt.Printf("Opcionais exportados (não instalados automaticamente): %s\n", strings.Join(pref.OptionalPlugins, ", "))
	}
	fmt.Println("Preferências importadas. As versões anteriores receberam o sufixo .before-import-…")
	return nil
}
