// Package power exposes the Power Profiles Daemon's live profiles through
// the session helper. Other desktop environments can change them too, so
// status always reads the daemon.
package power

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/getaroli/desktop/cmd/aroli/internal/sys"
)

// Battery implements `aroli battery ...`: status, available, next, set, or a
// bare profile name.
func Battery(args []string) error {
	if len(args) == 0 {
		args = []string{"status"}
	}
	valid := (len(args) == 1 && (args[0] == "status" || args[0] == "available" || args[0] == "next" || args[0] == "balanced" || args[0] == "performance" || args[0] == "power-saver")) || (len(args) == 2 && args[0] == "set" && (args[1] == "power-saver" || args[1] == "balanced" || args[1] == "performance"))
	if !valid {
		return errors.New("use: aroli battery status|available|set power-saver|balanced|performance|next")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	helper := filepath.Join(home, ".local", "bin", "battery-efficiency")
	helperArgs := args
	if len(args) == 1 && args[0] != "status" && args[0] != "available" && args[0] != "next" {
		helperArgs = []string{"set", args[0]}
	}
	if len(helperArgs) == 2 && helperArgs[0] == "set" && helperArgs[1] == "power-saver" {
		if game, _ := sys.ScriptState("game-mode", "status"); game == "on" {
			return errors.New("Game Mode está ativo; desligue-o antes de ativar economia de energia")
		}
	}
	cmd := exec.Command(helper, helperArgs...)
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("controle de perfil de energia indisponível: %w", err)
	}
	fmt.Println(strings.TrimSpace(string(out)))
	return nil
}
