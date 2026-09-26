// Package logs reads the local install records. Read-only by design.
package logs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ShowLogs implements `aroli logs [--list]`: the last 100 lines of the most
// recent install log, or the list of records.
func ShowLogs(args []string) error {
	if len(args) > 1 || (len(args) == 1 && args[0] != "--list") {
		return errors.New("use: aroli logs [--list]")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	directory := filepath.Join(home, ".local", "state", "aroli-desktop")
	files, err := filepath.Glob(filepath.Join(directory, "install-*.log"))
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return errors.New("ainda não há logs de instalação")
	}
	sort.Slice(files, func(i, j int) bool { return files[i] > files[j] })
	if len(args) == 1 {
		for _, file := range files {
			fmt.Println(file)
		}
		return nil
	}
	data, err := os.ReadFile(files[0])
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) > 100 {
		lines = lines[len(lines)-100:]
	}
	fmt.Printf("Último log: %s\n\n", files[0])
	fmt.Println(strings.Join(lines, "\n"))
	return nil
}
