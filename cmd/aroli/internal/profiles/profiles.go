// Package profiles names curated phase sets. Optional apps inside a
// profile stay visible and confirmable, never silent.
package profiles

import (
	"errors"
	"fmt"
	"strings"

	"github.com/getaroli/desktop/cmd/aroli/internal/install"
	"github.com/getaroli/desktop/cmd/aroli/internal/plugins"
)

// AroliProfile is one named, transparent installation profile.
type AroliProfile struct {
	Name, Description string
	Phases            []string
	Plugins           []string
}

// AroliProfiles is the curated catalog. Optional entries name packages the
// user still confirms separately.
var AroliProfiles = []AroliProfile{
	{"minimal", "Base visual do Aroli Desktop, sem aplicativos opcionais.", []string{"base", "config", "cursor", "final"}, nil},
	{"desktop", "Desktop completo: pacotes, gráficos, serviços e configuração.", []string{"base", "packages", "cursor", "config", "graphics", "services", "final"}, nil},
	{"creator", "Perfil desktop com ferramentas de criação explicitamente listadas.", []string{"base", "packages", "cursor", "config", "graphics", "services", "final"}, []string{"neovim", "yazi", "onefetch", "visual-studio-code-bin"}},
	{"gaming", "Perfil desktop; drivers e jogos continuam escolhas conscientes.", []string{"base", "packages", "cursor", "config", "graphics", "services", "final"}, nil},
}

// Profiles implements `aroli profile list|show|install NOME`.
func Profiles(args []string) error {
	if len(args) == 0 || args[0] == "list" {
		fmt.Println("\nPerfis de instalação")
		for _, p := range AroliProfiles {
			fmt.Printf("  %-10s %s\n", p.Name, p.Description)
		}
		return nil
	}
	if len(args) < 2 || (args[0] != "show" && args[0] != "install") {
		return errors.New("use: aroli profile list|show|install NOME")
	}
	p, err := profileByName(args[1])
	if err != nil {
		return err
	}
	fmt.Printf("%s\n  fases: %s\n", p.Description, strings.Join(p.Phases, ", "))
	if len(p.Plugins) > 0 {
		fmt.Printf("  opcionais: %s\n", strings.Join(p.Plugins, ", "))
	}
	if args[0] == "show" {
		return nil
	}
	forward := args[2:]
	if len(forward) == 0 {
		forward = []string{}
	}
	// The installer retains ownership of backups, preflight, and confirmations.
	if err := install.Install(append(forward, p.Phases...)); err != nil {
		return err
	}
	if len(p.Plugins) == 0 {
		return nil
	}
	// A profile names the optional packages, but still gives the user a separate
	// opportunity to decline them unless they deliberately supplied --yes.
	pluginArgs := []string{"install"}
	for _, arg := range forward {
		// Only these flags have a meaning for plugin installation. Installation
		// flags such as --lang and --repo must not leak into its parser.
		if arg == "--dry-run" || arg == "--yes" {
			pluginArgs = append(pluginArgs, arg)
		}
	}
	return plugins.Plugins(append(pluginArgs, p.Plugins...))
}

func profileByName(name string) (AroliProfile, error) {
	for _, p := range AroliProfiles {
		if p.Name == name {
			return p, nil
		}
	}
	return AroliProfile{}, fmt.Errorf("perfil %q não existe; veja: aroli profile list", name)
}
