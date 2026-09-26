// Package gaming coordinates the reversible compositor layer with Feral
// GameMode. The latter is deliberately opt-in: it only wraps a process
// passed to launch.
package gaming

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/getaroli/desktop/cmd/aroli/internal/reading"
	"github.com/getaroli/desktop/cmd/aroli/internal/sys"
)

type modeCoordination struct {
	BatteryBeforeGame string `json:"battery_before_game,omitempty"`
	UpdatedAt         string `json:"updated_at"`
}

func modeCoordinationPath() (string, error) {
	home, err := os.UserHomeDir()
	return filepath.Join(home, ".local", "state", "aroli-desktop", "mode-coordination.json"), err
}

func readModeCoordination() (modeCoordination, error) {
	path, err := modeCoordinationPath()
	if err != nil {
		return modeCoordination{}, err
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return modeCoordination{}, nil
	}
	if err != nil {
		return modeCoordination{}, err
	}
	var state modeCoordination
	return state, json.Unmarshal(data, &state)
}

func writeModeCoordination(state modeCoordination) error {
	path, err := modeCoordinationPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	state.UpdatedAt = time.Now().Format(time.RFC3339)
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// ClearModeCoordination drops the saved pre-game power profile, if any.
func ClearModeCoordination() error {
	path, err := modeCoordinationPath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// PrepareGameMode resolves the only meaningful conflict: power-saver can cap
// a game before Feral GameMode has a chance to optimise its process. Reading
// mode and Game Mode both own compositor effects, so reading is turned off
// first. The captured power profile is restored when Game Mode ends.
func PrepareGameMode() error {
	if err := reading.Reading([]string{"off"}); err != nil {
		return err
	}
	battery, err := sys.ScriptState("battery-efficiency", "status")
	if err != nil || battery != "power-saver" {
		return nil
	}
	path, err := sys.UserScript("battery-efficiency")
	if err != nil {
		return err
	}
	if out, err := exec.Command(path, "set", "balanced").Output(); err != nil {
		return fmt.Errorf("não foi possível preparar o perfil de energia para jogos: %w", err)
	} else if strings.TrimSpace(string(out)) != "balanced" {
		return errors.New("o perfil de energia não mudou para balanced")
	}
	return writeModeCoordination(modeCoordination{BatteryBeforeGame: battery})
}

// RestoreGameModeCoordination brings back the power profile captured before
// Game Mode started.
func RestoreGameModeCoordination() error {
	state, err := readModeCoordination()
	if err != nil || state.BatteryBeforeGame == "" {
		return err
	}
	path, err := sys.UserScript("battery-efficiency")
	if err != nil {
		return err
	}
	if _, err := exec.Command(path, "set", state.BatteryBeforeGame).Output(); err != nil {
		return err
	}
	return ClearModeCoordination()
}

// Gaming implements `aroli gaming status|on|off|toggle|launch`.
func Gaming(args []string) error {
	if len(args) == 0 {
		args = []string{"status"}
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	helper := filepath.Join(home, ".local", "bin", "game-mode")
	runHelper := func(action string) (string, error) {
		cmd := exec.Command(helper, action)
		out, err := cmd.Output()
		if err != nil {
			return "", fmt.Errorf("Game Mode não está disponível: %w", err)
		}
		return strings.TrimSpace(string(out)), nil
	}
	switch args[0] {
	case "status":
		if len(args) > 2 || (len(args) == 2 && args[1] != "--json") {
			return errors.New("use: aroli gaming status [--json]")
		}
		state, err := runHelper("status")
		if err != nil {
			return err
		}
		if len(args) == 2 {
			available := false
			if _, err := exec.LookPath("gamemoderun"); err == nil {
				available = true
			}
			return json.NewEncoder(os.Stdout).Encode(map[string]any{"active": state == "on", "gamemode_available": available})
		}
		fmt.Println(state)
		return nil
	case "on", "off", "toggle":
		if len(args) != 1 {
			return errors.New("use: aroli gaming on|off|toggle")
		}
		before, err := runHelper("status")
		if err != nil {
			return err
		}
		if (args[0] == "on" || args[0] == "toggle") && before != "on" {
			if err := PrepareGameMode(); err != nil {
				return err
			}
		}
		state, err := runHelper(args[0])
		if err != nil {
			if before != "on" {
				_ = RestoreGameModeCoordination()
			}
			return err
		}
		if before == "on" && state == "off" {
			if err := RestoreGameModeCoordination(); err != nil {
				return err
			}
		}
		fmt.Println(state)
		return nil
	case "launch":
		if len(args) < 2 {
			return errors.New("use: aroli gaming launch COMANDO [args...]")
		}
		wasActive, err := runHelper("status")
		if err != nil {
			return err
		}
		if wasActive != "on" {
			if err := PrepareGameMode(); err != nil {
				return err
			}
			if _, err := runHelper("on"); err != nil {
				_ = RestoreGameModeCoordination()
				return err
			}
		}
		program, programArgs := args[1], args[2:]
		if _, err := exec.LookPath("gamemoderun"); err == nil {
			programArgs = append([]string{program}, programArgs...)
			program = "gamemoderun"
		}
		cmd := sys.Command("", program, programArgs...)
		err = cmd.Run()
		if wasActive != "on" {
			if _, offErr := runHelper("off"); offErr != nil && err == nil {
				err = offErr
			}
			if restoreErr := RestoreGameModeCoordination(); restoreErr != nil && err == nil {
				err = restoreErr
			}
		}
		return err
	default:
		return errors.New("use: aroli gaming status|on|off|toggle|launch COMANDO [args...]")
	}
}
