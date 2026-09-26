// Package tui is the guided terminal assistant for first-time users.
// Explicit subcommands stay the auditable path for agents and automation.
package tui

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"

	"github.com/getaroli/desktop/cmd/aroli/internal/install"
	"github.com/getaroli/desktop/cmd/aroli/internal/logs"
	"github.com/getaroli/desktop/cmd/aroli/internal/plugins"
	"github.com/getaroli/desktop/cmd/aroli/internal/selfupdate"
	"github.com/getaroli/desktop/cmd/aroli/internal/sys"
)

type tuiScreen uint8

const (
	tuiHome tuiScreen = iota
	tuiLanguage
	tuiPlugins
	tuiConfirm
	tuiLoading
	tuiRunning
	tuiDone
)

type tuiAction struct {
	label, hint string
	command     func(*tuiModel) error
}

type tuiResult struct{ err error }

type tuiPluginsResult struct {
	items []plugins.Plugin
	err   error
}

type tuiActionCommand struct {
	run            func() error
	stdin          io.Reader
	stdout, stderr io.Writer
}

func (c *tuiActionCommand) Run() error            { return c.run() }
func (c *tuiActionCommand) SetStdin(r io.Reader)  { c.stdin = r }
func (c *tuiActionCommand) SetStdout(w io.Writer) { c.stdout = w }
func (c *tuiActionCommand) SetStderr(w io.Writer) { c.stderr = w }

type tuiModel struct {
	screen       tuiScreen
	cursor       int
	pluginCursor int
	action       int
	language     string
	dryRun       bool
	version      string
	plugins      []plugins.Plugin
	selected     map[int]bool
	message      string
	width        int
	height       int
}

var (
	tuiTitle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F2CC8F"))
	tuiAccent   = lipgloss.NewStyle().Foreground(lipgloss.Color("#81B29A"))
	tuiMuted    = lipgloss.NewStyle().Foreground(lipgloss.Color("#8D99AE"))
	tuiSelected = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F2CC8F"))
	tuiPanel    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#3D405B")).Padding(1, 2)
)

func newTUIModel(version string) tuiModel {
	return tuiModel{language: sys.DefaultLang, selected: map[int]bool{}, version: version}
}

func (m tuiModel) actions() []tuiAction {
	return []tuiAction{
		{label: "Instalar ou reparar", hint: "aplicar o aroli no sistema", command: func(m *tuiModel) error { return install.Install([]string{"--yes", "--lang", m.language}) }},
		{label: "Ver plano de instalação", hint: "simular sem alterar arquivos", command: func(m *tuiModel) error { return install.Install([]string{"--dry-run", "--yes", "--lang", m.language}) }},
		{label: "Adicionar plugins", hint: "escolher aplicativos oficiais e AUR", command: nil},
		{label: "Diagnosticar", hint: "verificar problemas conhecidos", command: func(*tuiModel) error { return sys.RunBackend("diagnose", nil) }},
		{label: "Verificar atualizações", hint: "consultar a versão estável", command: func(*tuiModel) error { return sys.RunBackend("aroli", []string{"check", "--force"}) }},
		{label: "Atualizar aroli", hint: "aplicar a versão estável mais recente", command: func(*tuiModel) error { return sys.RunBackend("aroli", []string{"update"}) }},
		{label: "Atualizar CLI", hint: "baixar a CLI com checksum SHA-256", command: func(m *tuiModel) error { return selfupdate.UpdateCLI(m.version, []string{"update"}) }},
		{label: "Ver último log", hint: "abrir o registro da instalação", command: func(*tuiModel) error { return logs.ShowLogs(nil) }},
	}
}

// Tui opens the guided assistant. It needs an interactive terminal.
func Tui(version string) error {
	if !termIsInteractive() {
		return errors.New("a TUI precisa de um terminal interativo; use 'aroli help' para os comandos")
	}
	_, err := tea.NewProgram(newTUIModel(version)).Run()
	return err
}

func termIsInteractive() bool {
	info, err := os.Stdout.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}

func (m tuiModel) Init() tea.Cmd { return nil }

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tuiResult:
		m.screen = tuiDone
		if msg.err != nil {
			m.message = "Erro: " + msg.err.Error()
		} else {
			m.message = "Ação concluída com sucesso."
		}
	case tuiPluginsResult:
		if msg.err != nil {
			m.message = "Não foi possível carregar os plugins: " + msg.err.Error()
			m.screen = tuiDone
		} else {
			m.plugins, m.selected, m.pluginCursor = msg.items, map[int]bool{}, 0
			m.screen = tuiPlugins
		}
	case tea.KeyPressMsg:
		key := msg.String()
		if key == "ctrl+c" || key == "q" || key == "esc" {
			if m.screen == tuiHome || m.screen == tuiDone {
				return m, tea.Quit
			}
			m.screen = tuiHome
			return m, nil
		}
		switch m.screen {
		case tuiHome:
			return m.updateHome(key)
		case tuiLanguage:
			return m.updateLanguage(key)
		case tuiPlugins:
			return m.updatePlugins(key)
		case tuiConfirm:
			if key == "y" || key == "enter" {
				m.screen = tuiRunning
				return m, m.runAction()
			}
			if key == "n" || key == "backspace" {
				m.screen = tuiHome
			}
		case tuiDone:
			if key == "enter" || key == "r" {
				m.screen = tuiHome
			}
		}
	}
	return m, nil
}

func (m tuiModel) updateHome(key string) (tea.Model, tea.Cmd) {
	actions := m.actions()
	switch key {
	case "up", "k":
		m.cursor = (m.cursor + len(actions) - 1) % len(actions)
	case "down", "j":
		m.cursor = (m.cursor + 1) % len(actions)
	case "enter", "1", "2", "3", "4", "5", "6", "7", "8":
		if key != "enter" {
			m.cursor = int(key[0] - '1')
		}
		m.action = m.cursor
		if m.cursor == 0 || m.cursor == 1 {
			m.dryRun = m.cursor == 1
			m.screen = tuiLanguage
		} else if m.cursor == 2 {
			m.screen = tuiLoading
			return m, m.loadPlugins()
		} else {
			m.screen = tuiConfirm
		}
	}
	return m, nil
}

func (m tuiModel) loadPlugins() tea.Cmd {
	return func() tea.Msg {
		var items []plugins.Plugin
		err := silenceTerminal(func() error {
			var err error
			items, err = plugins.PluginCatalog()
			return err
		})
		return tuiPluginsResult{items: items, err: err}
	}
}

func (m tuiModel) updateLanguage(key string) (tea.Model, tea.Cmd) {
	langs := []string{"pt-BR", "en", "es"}
	var index int
	for i, lang := range langs {
		if lang == m.language {
			index = i
		}
	}
	switch key {
	case "left", "h", "up":
		index = (index + len(langs) - 1) % len(langs)
	case "right", "l", "down":
		index = (index + 1) % len(langs)
	case "enter":
		m.language, m.screen = langs[index], tuiConfirm
	}
	if key != "enter" {
		m.language = langs[index]
	}
	return m, nil
}

func (m tuiModel) updatePlugins(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "up", "k":
		m.pluginCursor = (m.pluginCursor + len(m.plugins) - 1) % len(m.plugins)
	case "down", "j":
		m.pluginCursor = (m.pluginCursor + 1) % len(m.plugins)
	case "space":
		m.selected[m.pluginCursor] = !m.selected[m.pluginCursor]
	case "enter":
		if len(m.selectedPlugins()) == 0 {
			m.message = "Selecione pelo menos um plugin com espaço."
			return m, nil
		}
		m.screen = tuiConfirm
	}
	return m, nil
}

func (m tuiModel) selectedPlugins() []plugins.Plugin {
	items := []plugins.Plugin{}
	for i, item := range m.plugins {
		if m.selected[i] {
			items = append(items, item)
		}
	}
	return items
}

func (m tuiModel) runAction() tea.Cmd {
	return tea.Exec(&tuiActionCommand{run: func() error {
		if m.action == 2 {
			return plugins.InstallPlugins(m.selectedPlugins(), false)
		}
		return m.actions()[m.action].command(&m)
	}}, func(err error) tea.Msg { return tuiResult{err: err} })
}

// Bubble Tea owns the terminal while the TUI is active. The installer and its
// shell backend are intentionally verbose, so keep their output from writing
// over Bubble Tea's alternate screen and report only the final result in the
// interface. Direct CLI commands remain verbose as before.
func silenceTerminal(run func() error) error {
	stdout, stderr := os.Stdout, os.Stderr
	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		return run()
	}
	os.Stdout, os.Stderr = devNull, devNull
	defer func() {
		os.Stdout, os.Stderr = stdout, stderr
		_ = devNull.Close()
	}()
	return run()
}

func (m tuiModel) View() tea.View {
	var body string
	switch m.screen {
	case tuiHome:
		body = m.viewHome()
	case tuiLanguage:
		body = m.viewLanguage()
	case tuiPlugins:
		body = m.viewPlugins()
	case tuiConfirm:
		body = m.viewConfirm()
	case tuiLoading:
		body = "\n  Carregando plugins…\n\n  A conexão pode levar alguns segundos.\n  Pressione q ou Esc para voltar.\n"
	case tuiRunning:
		body = "\n  Executando…\n"
	case tuiDone:
		body = tuiTitle.Render("\n  "+m.message+"\n\n") + tuiMuted.Render("  Enter/r: voltar   q: sair")
	}
	var view tea.View
	view.SetContent(tuiPanel.Render(tuiTitle.Render("Aroli Desktop") + "\n" + tuiMuted.Render("aroli · assistente") + "\n" + tuiAccent.Render(m.breadcrumb()) + "\n\n" + body + "\n\n" + tuiMuted.Render("↑/↓ navegar · Enter selecionar · q sair")))
	view.AltScreen = true
	return view
}

func (m tuiModel) breadcrumb() string {
	switch m.screen {
	case tuiLanguage:
		if m.dryRun {
			return "Início › Instalação › Plano › Idioma"
		}
		return "Início › Instalação › Idioma"
	case tuiPlugins:
		return "Início › Plugins"
	case tuiConfirm:
		if m.action == 2 {
			return "Início › Plugins › Confirmar"
		}
		return "Início › Instalação › Confirmar"
	case tuiRunning:
		return "Início › Executando"
	case tuiDone:
		return "Início › Resultado"
	default:
		return "Início"
	}
}

func (m tuiModel) viewHome() string {
	lines := []string{"\n  Início", ""}
	for i, action := range m.actions() {
		prefix := "  "
		if i == m.cursor {
			prefix = tuiAccent.Render("› ")
		}
		lines = append(lines, prefix+tuiSelected.Render(fmt.Sprintf("%d. %-25s", i+1, action.label))+"  "+tuiMuted.Render(action.hint))
	}
	return strings.Join(lines, "\n")
}

func (m tuiModel) viewLanguage() string {
	return "\n  Idioma da instalação\n\n  " + tuiSelected.Render("‹  "+m.language+"  ›") + "\n\n  Use ←/→ e confirme com Enter."
}

func (m tuiModel) viewPlugins() string {
	lines := []string{"\n  Plugins opcionais · Espaço marca, Enter confirma", ""}
	for i, item := range m.plugins {
		mark := "○"
		if m.selected[i] {
			mark = tuiAccent.Render("●")
		}
		prefix := "  "
		if i == m.pluginCursor {
			prefix = tuiAccent.Render("› ")
		}
		lines = append(lines, prefix+mark+" "+item.Name+"  "+tuiMuted.Render(item.Source+" · "+item.Description))
	}
	return strings.Join(lines, "\n")
}

func (m tuiModel) viewConfirm() string {
	action := m.actions()[m.action].label
	if m.action == 2 {
		action = fmt.Sprintf("instalar %d plugin(s)", len(m.selectedPlugins()))
	}
	if m.dryRun {
		action += " (simulação)"
	}
	return "\n  Confirmar: " + tuiSelected.Render(action) + "\n\n  " + tuiAccent.Render("Enter/y") + " executar    " + tuiMuted.Render("n voltar")
}
