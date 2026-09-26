// Package doctor reconciles the small, user-local pieces of the delivery contract.
package doctor

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/getaroli/desktop/cmd/aroli/internal/materialize"
	"github.com/getaroli/desktop/cmd/aroli/internal/sys"
)

// Doctor implements `aroli doctor [--fix]`. Reconciliation only creates
// absent seeds; it never overwrites a local desktop choice.
func Doctor(args []string) error {
	fix := len(args) == 1 && args[0] == "--fix"
	if len(args) > 1 || (len(args) == 1 && !fix) {
		return errors.New("use: aroli doctor [--fix]")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	issues := []string{}
	userEdits := filepath.Join(home, ".config", "aroli", "user_edits")
	mimeapps := filepath.Join(home, ".config", "mimeapps.list")
	if _, err := os.Stat(userEdits); err != nil {
		issues = append(issues, "UserEdits: diretório de overrides ausente")
	}
	if _, err := os.Stat(mimeapps); err != nil {
		issues = append(issues, "MimeDefaults: mimeapps.list ainda não foi semeado")
	}
	for _, name := range []string{".zshrc.local", ".bashrc.local", ".zprofile.local", ".profile.local"} {
		if _, err := os.Stat(filepath.Join(home, name)); err != nil {
			issues = append(issues, "ShellLoad: "+name+" ausente")
		}
	}
	if !materialize.StatePresent(home) {
		issues = append(issues, "Manifest: nenhuma materialização registrada")
	}
	if len(issues) == 0 {
		fmt.Println("Doctor: UserEdits, MimeDefaults, ShellLoad e Manifest estão saudáveis.")
	} else {
		fmt.Println("Doctor encontrou:")
		for _, issue := range issues {
			fmt.Println("  -", issue)
		}
	}
	if !fix {
		fmt.Println("\nPara reconciliar seeds locais: aroli doctor --fix")
		return nil
	}
	if err := os.MkdirAll(userEdits, 0o755); err != nil {
		return err
	}
	if _, err := os.Stat(mimeapps); os.IsNotExist(err) {
		repo, cleanup, e := sys.EnsureRepositoryWithCleanup("", false)
		if e != nil {
			return e
		}
		defer cleanup()
		data, e := os.ReadFile(filepath.Join(repo, "seeds", "mimeapps.list"))
		if e != nil {
			return e
		}
		if e = os.WriteFile(mimeapps, data, 0o644); e != nil {
			return e
		}
	}
	for _, name := range []string{".zshrc.local", ".bashrc.local", ".zprofile.local", ".profile.local"} {
		path := filepath.Join(home, name)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := os.WriteFile(path, []byte("# Personal Aroli shell overrides. Updates never replace this file.\n"), 0o644); err != nil {
				return err
			}
		}
	}
	if !materialize.StatePresent(home) {
		fmt.Println("Doctor: execute 'aroli materialize' para criar o Manifest da entrega.")
	}
	// Package convergence is intentionally inside --fix and still asks: it can
	// reach pacman/sudo, unlike the four local reconcilers above.
	repo, cleanup, e := sys.EnsureRepositoryWithCleanup("", false)
	if e == nil {
		defer cleanup()
		if _, statErr := os.Stat(filepath.Join(repo, "manifest.json")); statErr == nil && sys.Confirm(bufio.NewReader(os.Stdin), "Convergir os pacotes deste box com manifest.json?") {
			if e := sys.Command(repo, "bash", "install.sh", "packages").Run(); e != nil {
				return e
			}
		}
	}
	fmt.Println("Doctor aplicou os reconcilers locais seguros.")
	return nil
}
