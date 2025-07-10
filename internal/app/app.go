package app

import (
	"fmt"

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
	return &Model{
		config:         cfg,
		viewController: NewViewController(),
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
		return m, nil

	case tea.KeyMsg:
		// Handle global keybindings
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

	accInfo := m.config.Account
	accText := "Cuenta configurada: "
	if len(accInfo) < 3 {
		accInfo = "Configurar Cuenta"
		accText = "Ingresa tu usuario DCC en "
	}

	accText = lipgloss.NewStyle().Foreground(m.theme.Selected).Render(accText)
	accInfo = lipgloss.NewStyle().Foreground(m.theme.Header).Render(accInfo)
	accRender := fmt.Sprintf("%s%s", accText, accInfo)
	content := lipgloss.JoinVertical(lipgloss.Left, header, view, accRender)
	centeredContent := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(content)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		centeredContent,
	)
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
