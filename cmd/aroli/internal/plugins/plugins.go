// Package plugins manages the opt-in catalog: every item needs an
// individual user decision and never installs by default.
package plugins

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/getaroli/desktop/cmd/aroli/internal/sys"
)

// Plugin is one opt-in application or tool.
type Plugin struct {
	Name, Description, Source string
}

// Plugins implements `aroli plugins list|install`.
func Plugins(args []string) error {
	if len(args) == 0 || args[0] == "list" {
		items, err := PluginCatalog()
		if err != nil {
			return err
		}
		fmt.Println("\nAplicativos e ferramentas opcionais")
		fmt.Println("Nada abaixo é necessário para o Aroli Desktop funcionar. Instale só o que fizer sentido para você.")
		for _, item := range items {
			fmt.Printf("  %-30s %-7s %s\n", item.Name, item.Source, item.Description)
		}
		return nil
	}
	if args[0] != "install" {
		return errors.New("use: aroli plugins list ou aroli plugins install NOME")
	}
	fs := flag.NewFlagSet("plugins install", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	dry := fs.Bool("dry-run", false, "mostra o que seria instalado")
	yes := fs.Bool("yes", false, "não pede confirmação")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	if len(fs.Args()) == 0 {
		return errors.New("escolha pelo menos um nome; veja as opções com: aroli plugins list")
	}
	items, err := PluginCatalog()
	if err != nil {
		return err
	}
	chosen := make([]Plugin, 0, len(fs.Args()))
	for _, name := range fs.Args() {
		found := false
		for _, item := range items {
			if item.Name == name {
				chosen = append(chosen, item)
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("%q não é um opcional conhecido; use: aroli plugins list", name)
		}
	}
	if !*yes && !*dry {
		fmt.Println("\nVocê escolheu:")
		for _, item := range chosen {
			fmt.Printf("  • %s, %s\n", item.Name, item.Description)
		}
		if !sys.Confirm(bufio.NewReader(os.Stdin), "Instalar estes itens agora?") {
			fmt.Println("Tudo bem, nada foi instalado.")
			return nil
		}
	}
	return InstallPlugins(chosen, *dry)
}

// PluginCatalog reads the optional manifests from the managed checkout.
func PluginCatalog() ([]Plugin, error) {
	repo, cleanup, err := sys.EnsureRepositoryWithCleanup("", true)
	if err != nil {
		return nil, err
	}
	defer cleanup()
	if _, err := sys.ValidRepo(repo); err != nil {
		return nil, errors.New("instale o aroli antes de gerenciar opcionais: aroli install")
	}
	items := []Plugin{}
	for _, manifest := range []struct{ file, source string }{{"optional-pacman.txt", "oficial"}, {"optional-aur.txt", "AUR"}} {
		file, err := os.Open(filepath.Join(repo, "packages", manifest.file))
		if err != nil {
			return nil, err
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "\t", 2)
			item := Plugin{Name: strings.TrimSpace(parts[0]), Source: manifest.source}
			if len(parts) == 2 {
				item.Description = strings.TrimSpace(parts[1])
			}
			items = append(items, item)
		}
		if err := scanner.Err(); err != nil {
			file.Close()
			return nil, err
		}
		file.Close()
	}
	return items, nil
}

// InstallPlugins installs the chosen catalog items, officials first.
func InstallPlugins(items []Plugin, dry bool) error {
	official, aur := []string{}, []string{}
	for _, item := range items {
		if item.Source == "AUR" {
			aur = append(aur, item.Name)
		} else {
			official = append(official, item.Name)
		}
	}
	if dry {
		if len(official) > 0 {
			fmt.Printf("seria instalado dos repositórios oficiais: %s\n", strings.Join(official, ", "))
		}
		if len(aur) > 0 {
			fmt.Printf("seria instalado do AUR: %s\n", strings.Join(aur, ", "))
		}
		return nil
	}
	if len(official) > 0 {
		if _, err := exec.LookPath("omarchy-pkg-add"); err == nil {
			if err := installWith("omarchy-pkg-add", official); err != nil {
				return err
			}
		} else {
			args := append([]string{"env", "OMARCHY_ALLOW_DIRECT_PACMAN=1", "pacman", "-S", "--needed", "--noconfirm", "--"}, official...)
			if err := sys.Command("", "sudo", args...).Run(); err != nil {
				return fmt.Errorf("não foi possível instalar os itens oficiais: %w", err)
			}
		}
	}
	if len(aur) > 0 {
		if _, err := exec.LookPath("omarchy-pkg-aur-add"); err == nil {
			if err := installWith("omarchy-pkg-aur-add", aur); err != nil {
				return err
			}
		} else {
			if _, err := exec.LookPath("yay"); err != nil {
				return errors.New("estes itens vêm do AUR, mas o yay não está instalado. Rode: aroli install aur")
			}
			args := append([]string{"-S", "--needed", "--noconfirm", "--"}, aur...)
			if err := sys.Command("", "yay", args...).Run(); err != nil {
				return fmt.Errorf("não foi possível instalar os itens do AUR: %w", err)
			}
		}
	}
	fmt.Println("\nPronto! Os itens escolhidos foram instalados.")
	return nil
}

func installWith(manager string, packages []string) error {
	for _, pkg := range packages {
		if err := sys.Command("", manager, pkg).Run(); err != nil {
			return fmt.Errorf("%s não conseguiu instalar %s: %w", manager, pkg, err)
		}
	}
	return nil
}
