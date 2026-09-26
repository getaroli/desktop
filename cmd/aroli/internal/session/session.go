// Package session applies coordinated daily-use profiles: one command,
// one predictable combination of game, power, and reading modes.
package session

import (
	"errors"
	"fmt"

	"github.com/getaroli/desktop/cmd/aroli/internal/gaming"
	"github.com/getaroli/desktop/cmd/aroli/internal/power"
	"github.com/getaroli/desktop/cmd/aroli/internal/reading"
	"github.com/getaroli/desktop/cmd/aroli/internal/sys"
)

type sessionMode struct {
	name, description, battery string
	game                       bool
}

var sessionModes = []sessionMode{
	{"laptop", "Economia de energia, sem efeitos de jogo.", "power-saver", false},
	{"desktop", "Equilíbrio para uso diário.", "balanced", false},
	{"gaming", "Desliga efeitos do compositor; mantém energia balanceada.", "balanced", true},
	{"creator", "Equilíbrio para criação e desenvolvimento.", "balanced", false},
}

// SessionProfile implements `aroli session list|status|apply ...`.
func SessionProfile(args []string) error {
	if len(args) == 0 || args[0] == "list" {
		for _, mode := range sessionModes {
			fmt.Printf("  %-8s %s\n", mode.name, mode.description)
		}
		return nil
	}
	if args[0] == "status" && len(args) == 1 {
		game, _ := sys.ScriptState("game-mode", "status")
		batteryState, _ := sys.ScriptState("battery-efficiency", "status")
		read, _ := reading.ReadingState()
		fmt.Printf("game: %s\nreading: %s\npower: %s\n", game, read, batteryState)
		return nil
	}
	if len(args) != 2 || args[0] != "apply" {
		return errors.New("use: aroli session list|status|apply laptop|desktop|gaming|creator")
	}
	var wanted *sessionMode
	for i := range sessionModes {
		if sessionModes[i].name == args[1] {
			wanted = &sessionModes[i]
		}
	}
	if wanted == nil {
		return fmt.Errorf("perfil de sessão %q não existe", args[1])
	}
	if wanted.game {
		if err := power.Battery([]string{"set", wanted.battery}); err != nil {
			return err
		}
		if err := reading.Reading([]string{"off"}); err != nil {
			return err
		}
		return gaming.Gaming([]string{"on"})
	}
	if err := gaming.Gaming([]string{"off"}); err != nil {
		return err
	}
	if err := power.Battery([]string{"set", wanted.battery}); err != nil {
		return err
	}
	return reading.Reading([]string{"off"})
}
