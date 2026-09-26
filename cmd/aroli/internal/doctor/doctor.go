// Package doctor checks the local desktop integration and offers safe
// session repairs. It reports; only --fix with confirmation acts.
package doctor

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/getaroli/desktop/cmd/aroli/internal/sys"
)

// Doctor implements `aroli doctor [--fix]`.
func Doctor(args []string) error {
	fix := len(args) == 1 && args[0] == "--fix"
	if len(args) > 1 || (len(args) == 1 && !fix) {
		return errors.New("use: aroli doctor [--fix]")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	script := filepath.Join(home, ".config", "hypr", "set-wallpaper.sh")
	issues := []string{}
	if _, err := os.Stat(script); err != nil {
		issues = append(issues, "script de wallpaper ausente")
	}
	if _, err := exec.LookPath("hyprctl"); err != nil {
		issues = append(issues, "hyprctl não está disponível nesta sessão")
	}
	if _, err := exec.LookPath("quickshell"); err != nil {
		issues = append(issues, "quickshell não está instalado")
	}
	if len(issues) == 0 {
		fmt.Println("Doctor: integrações locais essenciais parecem saudáveis.")
	} else {
		fmt.Println("Doctor encontrou:")
		for _, issue := range issues {
			fmt.Println("  -", issue)
		}
	}
	if !fix {
		fmt.Println("\nPara reaplicar integrações seguras da sessão: aroli doctor --fix")
		return nil
	}
	if !sys.Confirm(bufio.NewReader(os.Stdin), "Recarregar Hyprland e reiniciar o serviço Quickshell do usuário?") {
		return nil
	}
	if _, err := exec.LookPath("hyprctl"); err == nil {
		_ = sys.Command("", "hyprctl", "reload").Run()
	}
	if _, err := exec.LookPath("systemctl"); err == nil {
		_ = sys.Command("", "systemctl", "--user", "restart", "quickshell.service").Run()
	}
	fmt.Println("Doctor aplicou as reparações de sessão disponíveis.")
	return nil
}
