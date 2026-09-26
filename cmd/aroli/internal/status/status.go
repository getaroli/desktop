// Package status reports the installed CLI, session health, pending
// updates, and local snapshots. Read-only by design.
package status

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/getaroli/desktop/cmd/aroli/internal/snapshot"
)

// Status implements `aroli status`.
func Status(version string) error {
	fmt.Printf("Aroli Desktop CLI %s\n", version)
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	checks := []struct {
		label, command string
		args           []string
	}{{"Hyprland", "hyprctl", []string{"version"}}, {"Quickshell", "systemctl", []string{"--user", "is-active", "quickshell.service"}}, {"wallpaper", "", nil}}
	healthy := 0
	for _, check := range checks {
		if check.label == "wallpaper" {
			if _, err := os.Stat(filepath.Join(home, ".cache", "wal", "lockbg")); err == nil {
				fmt.Printf("%-12s saudável\n", check.label)
				healthy++
			} else {
				fmt.Printf("%-12s atenção\n", check.label)
			}
			continue
		}
		if _, err := exec.LookPath(check.command); err != nil {
			fmt.Printf("%-12s indisponível\n", check.label)
			continue
		}
		cmd := exec.Command(check.command, check.args...)
		cmd.Stdout, cmd.Stderr = io.Discard, io.Discard
		if cmd.Run() == nil {
			fmt.Printf("%-12s saudável\n", check.label)
			healthy++
		} else {
			fmt.Printf("%-12s atenção\n", check.label)
		}
	}
	fmt.Printf("saúde       %d/%d\n", healthy, len(checks))
	var update struct {
		Remote    string `json:"remote"`
		Available bool   `json:"update_available"`
	}
	if data, err := os.ReadFile(filepath.Join(home, ".cache", "aroli-desktop", "update.json")); err == nil && json.Unmarshal(data, &update) == nil && update.Remote != "" {
		state := "em dia"
		if update.Available {
			state = "atualização disponível: " + update.Remote
		}
		fmt.Printf("atualização  %s\n", state)
	} else {
		fmt.Println("atualização  ainda não verificada (aroli check --force)")
	}
	if root, err := snapshot.SnapshotRoot(); err == nil {
		if entries, err := os.ReadDir(root); err == nil {
			count := 0
			for _, entry := range entries {
				if entry.IsDir() {
					count++
				}
			}
			fmt.Printf("snapshots    %d local(is)\n", count)
		}
	}
	return nil
}
