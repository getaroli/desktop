// Package recover unwinds stuck modes and optionally restores the newest
// local snapshot. It confirms before changing anything.
package recover

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/getaroli/desktop/cmd/aroli/internal/gaming"
	"github.com/getaroli/desktop/cmd/aroli/internal/reading"
	"github.com/getaroli/desktop/cmd/aroli/internal/snapshot"
	"github.com/getaroli/desktop/cmd/aroli/internal/sys"
)

// RecoverDesktop implements `aroli recover [--dry-run|--yes] [--snapshot NOME]`.
func RecoverDesktop(args []string) error {
	dry, yes, snap := false, false, ""
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--dry-run":
			dry = true
		case "--yes":
			yes = true
		case "--snapshot":
			if i+1 == len(args) {
				return errors.New("--snapshot precisa de um nome")
			}
			i++
			snap = args[i]
		default:
			return errors.New("use: aroli recover [--dry-run|--yes] [--snapshot NOME]")
		}
	}
	if snap == "latest" {
		root, err := snapshot.SnapshotRoot()
		if err != nil {
			return err
		}
		entries, err := os.ReadDir(root)
		if err != nil {
			return errors.New("não há snapshot local para restaurar")
		}
		var latest string
		var latestTime time.Time
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			info, err := entry.Info()
			if err == nil && (latest == "" || info.ModTime().After(latestTime)) {
				latest, latestTime = entry.Name(), info.ModTime()
			}
		}
		if latest == "" {
			return errors.New("não há snapshot local para restaurar")
		}
		snap = latest
	}
	fmt.Println("Recuperação: desligaria Game Mode e Modo leitura, recarregaria Hyprland e reiniciaria Quickshell.")
	if snap != "" {
		fmt.Printf("Também restauraria o snapshot %q.\n", snap)
	}
	if dry {
		return nil
	}
	if !yes && !sys.Confirm(bufio.NewReader(os.Stdin), "Aplicar esta recuperação?") {
		return nil
	}
	_, _ = sys.ScriptState("game-mode", "off")
	_ = reading.Reading([]string{"off"})
	_ = gaming.ClearModeCoordination()
	if snap != "" {
		if err := snapshot.Snapshots([]string{"restore", snap, "--yes"}); err != nil {
			return err
		}
	}
	if _, err := exec.LookPath("hyprctl"); err == nil {
		_ = exec.Command("hyprctl", "reload").Run()
	}
	if _, err := exec.LookPath("systemctl"); err == nil {
		_ = exec.Command("systemctl", "--user", "restart", "quickshell.service").Run()
	}
	fmt.Println("Recuperação concluída. Execute aroli diagnose para um relatório completo.")
	return nil
}
