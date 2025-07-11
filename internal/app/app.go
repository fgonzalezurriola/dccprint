package app

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/fgonzalezurriola/dccprint/internal/account"
	"github.com/fgonzalezurriola/dccprint/internal/components"
	"github.com/fgonzalezurriola/dccprint/internal/components/scripts"
	"github.com/fgonzalezurriola/dccprint/internal/config"
	"github.com/fgonzalezurriola/dccprint/internal/theme"
)

type viewState int

type setupStep int

const (
	setupWelcome setupStep = iota
	setupAccount
	setupPrintConfig
	setupConfirm
)

type Model struct {
	config         config.Config
	viewController *ViewController
	mainMenu       components.Menu
	PrintView      components.PrintView
	ConfigView     components.Menu
	themeMenu      components.Menu
	theme          *theme.Theme
	themeManager   *theme.Manager
	width          int
	height         int
	accountManager account.Manager

	// Setup state
	inSetup        bool
	setupStep      setupStep
	setupAccountIn textinput.Model
	setupConfigV   *components.ConfigView
	setupDone      bool
}

// --- Component Initializers ---
func newMainMenu(t *theme.Theme) components.Menu {
	mainMenuItems := []string{"Imprimir PDF", "Configuración de Impresión", "Configurar Cuenta", "Cambiar Theme", "Salir"}
	return components.NewMenu(mainMenuItems, t)
}

func newThemeMenu(t *theme.Theme) components.Menu {
	themeMenuItems := []string{"Default", "Cadcc", "Anakena"}
	return components.NewMenu(themeMenuItems, t)
}

func newPrintView(t *theme.Theme) components.PrintView {
	return components.NewPrintView(scripts.GetPDFFiles(), t)
}

func newAccountManager(t *theme.Theme, cfg config.Config) account.Manager {
	return account.NewManager(t, cfg)
}

// --- Model ---
func NewModel() *Model {
	cfg := config.Load()
	t := theme.New(cfg.Theme)

	ti := textinput.New()
	ti.Placeholder = "Ingresa el nombre de cuenta (sin @)"
	ti.Focus()
	ti.CharLimit = 16
	ti.SetValue(cfg.Account)
	ti.PromptStyle = lipgloss.NewStyle().Foreground(t.Selected)
	ti.TextStyle = lipgloss.NewStyle().Foreground(t.Header)

	themeManager := theme.NewManager(cfg.Theme)
	vc := NewViewController()
	if !cfg.SetupCompleted {
		vc.Set(SetupView)
	}
	return &Model{
		config:         cfg,
		viewController: vc,
		mainMenu:       newMainMenu(t),
		PrintView:      newPrintView(t),
		themeMenu:      newThemeMenu(t),
		theme:          t,
		themeManager:   themeManager,
		accountManager: newAccountManager(t, cfg),
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

// --- Main Update ---
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle global messages first
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.mainMenu.SetSize(msg.Width, msg.Height)
		m.themeMenu.SetSize(msg.Width, msg.Height)
		m.PrintView.SetSize(msg.Width, msg.Height)
		if m.setupConfigV != nil {
			m.setupConfigV.SetSize(msg.Width, msg.Height)
		}
		return m, nil

	// Handle global keybindings
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			if m.viewController.Get() != MainView {
				m.mainMenu.Reset()
				m.viewController.Set(MainView)
				return m, nil
			}
		}
	}

	switch m.viewController.Get() {
	case SetupView:
		return m.updateSetupView(msg)
	case MainView:
		return m.updateMainView(msg)
	case PrintView:
		return m.updatePrintView(msg)
	case ThemeView:
		return m.updateThemeView(msg)
	case AccountView:
		return m.updateAccountView(msg)
	}
	return m, nil
}

// --- Setup View Update ---
func (m *Model) updateSetupView(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.inSetup {
		m.inSetup = true
		m.setupStep = setupWelcome
		m.setupAccountIn = m.accountManager.AccountInput
		m.setupConfigV = components.NewConfigView(m.theme)
		m.setupConfigV.SetSize(m.width, m.height)
		return m, nil
	}

	switch m.setupStep {
	case setupWelcome:
		if key, ok := msg.(tea.KeyMsg); ok && (key.String() == "enter" || key.String() == " ") {
			m.setupStep = setupAccount
			return m, nil
		}
		return m, nil
	case setupAccount:
		var inputCmd tea.Cmd
		m.setupAccountIn, inputCmd = m.setupAccountIn.Update(msg)
		if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
			account := m.setupAccountIn.Value()
			// Empty account verification (dcc accounts are exactly 8 chars)
			if len(account) < 1 {
				return m, nil
			}
			m.config.Account = account
			_ = config.SaveAccount(account)
			m.setupStep = setupPrintConfig
			return m, nil
		}
		return m, inputCmd
	case setupPrintConfig:
		newV, cmd := m.setupConfigV.Update(msg)
		m.setupConfigV = newV.(*components.ConfigView)
		if cmd != nil {
			msgResult := cmd()
			if _, ok := msgResult.(components.ConfigFinishedMsg); ok {
				// Guardar config actualizada
				m.config = config.Load() // Recarga config con los cambios
				m.setupStep = setupConfirm
				return m, nil
			}
			return m, cmd
		}
		return m, nil
	case setupConfirm:
		if key, ok := msg.(tea.KeyMsg); ok && (key.String() == "enter" || key.String() == " ") {
			m.config.SetupCompleted = true
			_ = config.SaveConfig(m.config) // Guardar toda la config, incluyendo el flag
			m.inSetup = false
			m.viewController.Set(MainView)
			return m, nil
		}
		return m, nil
	}
	return m, nil
}

// --- Update helpers ---
func (m *Model) updateMainView(msg tea.Msg) (tea.Model, tea.Cmd) {
	newMenu, menuCmd := m.mainMenu.Update(msg)
	m.mainMenu = newMenu.(components.Menu)
	_ = menuCmd

	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
		switch m.mainMenu.SelectedItem() {
		case "Imprimir PDF":
			m.viewController.Set(PrintView)
		case "Configurar Cuenta":
			m.viewController.Set(AccountView)
			m.accountManager.AccountInput.Focus()
		case "Cambiar Theme":
			m.themeMenu.Reset()
			m.viewController.Set(ThemeView)
		case "Salir":
			return m, tea.Quit
		}
	}
	return m, menuCmd
}

func (m *Model) updatePrintView(msg tea.Msg) (tea.Model, tea.Cmd) {
	newSelector, selectorCmd := m.PrintView.Update(msg)
	m.PrintView = newSelector.(components.PrintView)
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
		// TODO: lógica de impresión
		m.viewController.Set(MainView)
	}
	return m, selectorCmd
}

func (m *Model) updateThemeView(msg tea.Msg) (tea.Model, tea.Cmd) {
	newMenu, menuCmd := m.themeMenu.Update(msg)
	m.themeMenu = newMenu.(components.Menu)
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
		selectedTheme := m.themeMenu.SelectedItem()
		m.themeManager.ChangeTheme(selectedTheme)
		m.theme = theme.New(m.themeManager.Current)
		m.mainMenu.SetTheme(m.theme)
		m.themeMenu.SetTheme(m.theme)
		m.accountManager.AccountInput.PromptStyle = lipgloss.NewStyle().Foreground(m.theme.Selected)
		m.accountManager.AccountInput.TextStyle = lipgloss.NewStyle().Foreground(m.theme.Header)
		m.mainMenu.Reset()
		m.viewController.Set(MainView)
	}
	return m, menuCmd
}

func (m *Model) updateAccountView(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
		m.accountManager.SaveAccount()
		m.mainMenu.Reset()
		m.viewController.Set(MainView)
		return m, nil
	} else {
		var inputCmd tea.Cmd
		m.accountManager.AccountInput, inputCmd = m.accountManager.AccountInput.Update(msg)
		return m, inputCmd
	}
}

// --- Main View ---
func (m *Model) View() string {
	header := components.RenderHeader(m.width, m.theme)
	var view string

	switch m.viewController.Get() {
	case SetupView:
		view = m.viewSetup()
	case MainView:
		view = m.viewMain()
	case PrintView:
		view = m.viewPrint()
	case ConfigView:
		view = m.viewMain() // Todo
	case AccountView:
		view = m.viewAccount()
	case ThemeView:
		view = m.viewTheme()
	}
	content := lipgloss.JoinVertical(lipgloss.Left, header, view)
	centeredContent := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(content)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		centeredContent,
	)
}

// --- Setup View Render ---
func (m *Model) viewSetup() string {
	headerStyle := lipgloss.NewStyle().Foreground(m.theme.Header).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(m.theme.Selected)
	promptStyle := lipgloss.NewStyle().Foreground(m.theme.Unselected)
	stepStyle := lipgloss.NewStyle().Foreground(m.theme.Selected).Bold(true)

	switch m.setupStep {
	case setupWelcome:
		header := headerStyle.Render("Bienvenido a DCCPrint")
		desc := descStyle.Render("Continuaremos con la configuración inicial (se puede cambiar a futuro)")
		prompt := promptStyle.Render("[Presiona Enter para continuar]")
		return lipgloss.JoinVertical(lipgloss.Left, "", header, "", desc, "", prompt)
	case setupAccount:
		step := stepStyle.Render("Paso 1: Configura tu usuario")
		prompt := promptStyle.Render("[Enter para continuar]")
		return lipgloss.JoinVertical(lipgloss.Left, "", step, m.setupAccountIn.View(), "", prompt)
	case setupPrintConfig:
		step := stepStyle.Render("Paso 2: Configura la impresión")
		// Detectar si el modo seleccionado es Simple
		mode := m.config.Mode
		if mode == "Simple" {
			prompt := promptStyle.Render("[Enter para finalizar]")
			return lipgloss.JoinVertical(lipgloss.Left, "", step, m.setupConfigV.View(), "", prompt)
		}
		prompt := promptStyle.Render("[Enter para continuar]")
		return lipgloss.JoinVertical(lipgloss.Left, "", step, m.setupConfigV.View(), "", prompt)
	case setupConfirm:
		done := headerStyle.Render("¡Listo!")
		desc := descStyle.Render("Tu configuración ha sido guardada.")
		prompt := promptStyle.Render("[Presiona Enter para ir al menú principal]")
		return lipgloss.JoinVertical(lipgloss.Left, "", done, "", desc, "", prompt)
	default:
		return ""
	}
}

// --- View helpers ---
func (m *Model) viewMain() string {
	return m.mainMenu.View()
}

func (m *Model) viewPrint() string {
	return m.PrintView.View()
}

func (m *Model) viewTheme() string {
	return m.themeMenu.View()
}

func (m *Model) viewAccount() string {
	return m.accountManager.AccountInput.View()
}
