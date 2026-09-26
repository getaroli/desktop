// aroli is the supported command-line interface for Aroli Desktop.
//
// The installer itself deliberately remains a compatibility backend for now:
// it contains years of idempotency and backup rules. This program owns the
// public interface, bootstrap flow and beginner-friendly terminal menu.
//
// Command implementations live in internal/, one concern per package; this
// file is only the dispatcher, the help text, and the version.
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/getaroli/desktop/cmd/aroli/internal/doctor"
	"github.com/getaroli/desktop/cmd/aroli/internal/gaming"
	"github.com/getaroli/desktop/cmd/aroli/internal/install"
	"github.com/getaroli/desktop/cmd/aroli/internal/logs"
	"github.com/getaroli/desktop/cmd/aroli/internal/plugins"
	"github.com/getaroli/desktop/cmd/aroli/internal/power"
	"github.com/getaroli/desktop/cmd/aroli/internal/prefs"
	"github.com/getaroli/desktop/cmd/aroli/internal/profiles"
	"github.com/getaroli/desktop/cmd/aroli/internal/reading"
	"github.com/getaroli/desktop/cmd/aroli/internal/recover"
	"github.com/getaroli/desktop/cmd/aroli/internal/selfupdate"
	"github.com/getaroli/desktop/cmd/aroli/internal/session"
	"github.com/getaroli/desktop/cmd/aroli/internal/snapshot"
	"github.com/getaroli/desktop/cmd/aroli/internal/status"
	"github.com/getaroli/desktop/cmd/aroli/internal/sys"
	"github.com/getaroli/desktop/cmd/aroli/internal/tui"
	"github.com/getaroli/desktop/cmd/aroli/internal/wallpaper"
)

var version = "dev" // set with -ldflags "-X main.version=vX.Y.Z"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "aroli:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return tui.Tui(cliVersion())
	}
	switch args[0] {
	case "install":
		return install.Install(args[1:])
	case "plugins", "plugin":
		return plugins.Plugins(args[1:])
	case "logs":
		return logs.ShowLogs(args[1:])
	case "doctor":
		return doctor.Doctor(args[1:])
	case "profile", "profiles":
		return profiles.Profiles(args[1:])
	case "snapshot", "snapshots":
		return snapshot.Snapshots(args[1:])
	case "wallpaper", "wallpapers":
		return wallpaper.Wallpapers(args[1:])
	case "gaming", "game":
		return gaming.Gaming(args[1:])
	case "battery", "efficiency":
		return power.Battery(args[1:])
	case "session", "mode":
		return session.SessionProfile(args[1:])
	case "reading":
		return reading.Reading(args[1:])
	case "recover":
		return recover.RecoverDesktop(args[1:])
	case "export":
		return prefs.ExportPreferences(args[1:])
	case "import":
		return prefs.ImportPreferences(args[1:])
	case "cli":
		return selfupdate.UpdateCLI(cliVersion(), args[1:])
	case "diagnose":
		return sys.RunBackend("diagnose", args[1:])
	case "version":
		fmt.Println(cliVersion())
		return nil
	case "help", "--help", "-h":
		usage(os.Stdout)
		return nil
	case "status":
		return status.Status(cliVersion())
	case "check", "update", "rollback", "prune":
		return sys.RunBackend("aroli", args)
	default:
		return fmt.Errorf("comando desconhecido %q (use 'aroli help')", args[0])
	}
}

func usage(w io.Writer) {
	fmt.Fprintf(w, `%s CLI

Uso:
  aroli                     abre o assistente interativo
  aroli install [opções] [fase...]
  aroli diagnose
  aroli logs [--list]         mostra o registro da instalação mais recente
  aroli doctor [--fix]        verifica e repara integrações locais conhecidas
  aroli profile list|show|install NOME
  aroli snapshot create|list|restore NOME
  aroli wallpaper list|set|random|import|remove
  aroli gaming status|on|off|toggle|launch COMANDO [args...]
  aroli battery status|available|set power-saver|balanced|performance|next
  aroli session list|status|apply laptop|desktop|gaming|creator
  aroli reading on|off|status
  aroli recover [--dry-run|--yes] [--snapshot NOME]
  aroli export DIRETÓRIO | aroli import DIRETÓRIO [--yes]
  aroli cli update [--dry-run] atualiza somente o binário da CLI
  aroli plugins list|install [opções] [nome...]
  aroli status | check | update | rollback | prune

Instalação:
  --dry-run                mostra o plano, sem alterar nada
  --yes                    não pergunta confirmações
  --copy                   copia os arquivos em vez de criar links
  --link                   cria links para o checkout (padrão)
  --lang pt-BR|en|es       idioma da interface
  --repo CAMINHO           usa este checkout em vez do gerenciado pela CLI
  --resume                 continua uma instalação interrompida

Opcionais:
  aroli plugins list         mostra apps e ferramentas que podem ser adicionadas
  aroli plugins install NOME instala somente os itens escolhidos

Fases: base, aur, packages, repos, cursor, config, system, graphics,
       services, sddm, spicetify, final ou restore.

Exemplos:
  aroli install --dry-run
  aroli install --lang pt-BR
  aroli install packages services
  aroli update --dry-run
`, sys.Project)
}

func cliVersion() string {
	if version != "dev" {
		return strings.TrimPrefix(version, "v")
	}
	if repo, err := sys.ValidRepo("."); err == nil {
		if data, err := os.ReadFile(filepath.Join(repo, "VERSION")); err == nil {
			return strings.TrimSpace(string(data))
		}
	}
	return version
}
