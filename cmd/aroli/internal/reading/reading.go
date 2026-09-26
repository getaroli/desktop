// Package reading drives the eye-strain reading mode helper.
package reading

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/getaroli/desktop/cmd/aroli/internal/sys"
)

// Reading implements `aroli reading on|off|status`.
func Reading(args []string) error {
	if len(args) != 1 || (args[0] != "on" && args[0] != "off" && args[0] != "status") {
		return errors.New("use: aroli reading on|off|status")
	}
	if args[0] == "on" {
		if state, _ := sys.ScriptState("game-mode", "status"); state == "on" {
			if _, err := sys.ScriptState("game-mode", "off"); err != nil {
				return err
			}
		}
	}
	pathHome, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	out, err := exec.Command(filepath.Join(pathHome, ".config", "hypr", "scripts", "reading-mode.sh"), args[0]).Output()
	if err != nil {
		return fmt.Errorf("modo leitura indisponível: %w", err)
	}
	fmt.Println(strings.TrimSpace(string(out)))
	return nil
}

// ReadingState reports the helper's current mode.
func ReadingState() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	out, err := exec.Command(filepath.Join(home, ".config", "hypr", "scripts", "reading-mode.sh"), "status").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
