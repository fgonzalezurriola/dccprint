package app

import (
	"fmt"
	"log"

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
	PrinterView    components.PrinterView
	ModeView       components.ModeView
	themeMenu      components.Menu
	theme          *theme.Theme
	themeManager   *theme.Manager
	accountManager account.Manager
	freshManager   account.FreshManager
	width          int
	height         int
	printCompleted bool
}

// --- Component Initializers ---
func newMainMenu(t *theme.Theme) components.Menu {
	mainMenuItems := []string{"Imprimir PDF", "Configuración de Impresión", "Configurar Cuenta", "Cambiar Theme", "Salir"}
	return components.NewMenu(mainMenuItems, t)
}

func newThemeMenu(t *theme.Theme) components.Menu {
	themeMenuItems := []string{"Default", "Dcc...", "Anakena"}
	return components.NewMenu(themeMenuItems, t)
}

func newPrintView(t *theme.Theme) components.PrintView {
	return components.NewPrintView(scripts.GetPDFFiles(), t)
}

func newAccountManager(t *theme.Theme, cfg config.Config) account.Manager {
	return account.NewManager(t, cfg)
}

func newPrinterView(t *theme.Theme) components.PrinterView {
	printerMenuItems := []string{"Salita", "Toqui"}
	return components.NewPrinterView(printerMenuItems, t)
}

func newModeView(t *theme.Theme) components.ModeView {
	modeMenuItems := []string{"Doble cara, Borde largo (Recomendado)", "Doble cara, Borde corto", "Simple (Reverso en blanco)"}
	return components.NewModeView(modeMenuItems, t)
}

func newTextInput(ti textinput.Model, t *theme.Theme, cfg config.Config) {
	ti.Placeholder = "Ingresa el nombre de cuenta (sin @)"
	ti.Focus()
	ti.CharLimit = 16
	ti.SetValue(cfg.Account)
	ti.PromptStyle = lipgloss.NewStyle().Foreground(t.Selected)
	ti.TextStyle = lipgloss.NewStyle().Foreground(t.Header)
}

// --- Model ---
func NewModel() *Model {
	cfg := config.Load()
	t := theme.New(cfg.Theme)
	newTextInput(textinput.New(), t, cfg)
	themeManager := theme.NewManager(cfg.Theme)
	vc := NewViewController()

	model := &Model{
		config:         cfg,
		viewController: vc,
		mainMenu:       newMainMenu(t),
		PrintView:      newPrintView(t),
		PrinterView:    newPrinterView(t),
		ModeView:       newModeView(t),
		themeMenu:      newThemeMenu(t),
		theme:          t,
		themeManager:   themeManager,
		accountManager: newAccountManager(t, cfg),
		freshManager:   account.NewFreshManager(t),
	}

	if cfg.Account == "" {
		model.viewController.Set(FreshView)
	}

	return model
}

func (m *Model) Init() tea.Cmd {
	return nil
}

// --- Main  ---
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle global messages first
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.mainMenu.SetSize(msg.Width, msg.Height)
		m.themeMenu.SetSize(msg.Width, msg.Height)
		m.PrintView.SetSize(msg.Width, msg.Height)
		m.PrinterView.SetSize(msg.Width, msg.Height)
		m.ModeView.SetSize(msg.Width, msg.Height)

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
	case MainView:
		return m.updateMainView(msg)
	case PrintView:
		return m.updatePrintView(msg)
	case PrinterView:
		return m.updatePrinterView(msg)
	case ModeView:
		return m.updateModeView(msg)
	case ThemeView:
		return m.updateThemeView(msg)
	case AccountView:
		return m.updateAccountView(msg)
	case FreshView:
		return m.updateFreshView(msg)
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
		case "Configuración de Impresión":
			m.viewController.Set(PrinterView)
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

	// Meanwhile printCompleted is active, just the view is shown
	if m.printCompleted {
		if key, ok := msg.(tea.KeyMsg); ok && (key.String() == "enter" || key.String() == "q" || key.String() == "ctrl+c") {
			return m, tea.Quit
		}
		return m, nil
	}

	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
		filename := m.PrintView.SelectedItem()
		scriptName, err := scripts.CreateScript(filename)
		if err != nil {
			log.Fatalf("Error creando script: %v\n", err)
		}
		command := fmt.Sprintf("./%s", scriptName)
		if err := scripts.CopyToClipboard(command); err != nil {
			m.PrintView.StatusMessage = fmt.Sprintf("Error copiando al clipboard: %v", err)
		} else {
			m.PrintView.StatusMessage = "Script generado exitosamente!\n" +
				"Nombre del script generado: " + scriptName + "\n" +
				"Comando copiado al clipboard: " + command +
				"\nInstrucciones:\n" +
				"> Ctrl+Shift+V + Enter para ejecutar el script\n" +
				"> Ingresa tu contraseña SSH cuando se solicite\n" +
				"\nPresiona Enter, q o Ctrl+C para salir."
		}
		m.printCompleted = true
		return m, nil
	}

	return m, selectorCmd
}
func (m *Model) updatePrinterView(msg tea.Msg) (tea.Model, tea.Cmd) {
	newMenu, menuCmd := m.PrinterView.Menu.Update(msg)
	m.PrinterView.Menu = newMenu.(components.Menu)
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
		selectedPrinter := m.PrinterView.Menu.SelectedItem()
		config.SavePrinter(selectedPrinter)
		m.viewController.Set(ModeView)
	}
	return m, menuCmd
}

func (m *Model) updateModeView(msg tea.Msg) (tea.Model, tea.Cmd) {
	newMenu, menuCmd := m.ModeView.Menu.Update(msg)
	m.ModeView.Menu = newMenu.(components.Menu)
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
		selectedMode := m.ModeView.Menu.SelectedItem()
		config.SaveMode(selectedMode)
		m.viewController.Set(MainView)
	}
	return m, menuCmd
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
	case PrinterView:
		view = m.viewPrinter()
	case ModeView:
		view = m.viewMode()
	case AccountView:
		view = m.viewAccount()
	case ThemeView:
		view = m.viewTheme()
	case FreshView:
		view = m.viewFreshView()
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

func (m *Model) viewPrinter() string {
	return m.PrinterView.View()
}

func (m *Model) viewMode() string {
	return m.ModeView.View()
}

func (m *Model) updateFreshView(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
		m.freshManager.SaveAccount()
		m.viewController.Set(PrinterView)
		return m, nil
	} else {
		var inputCmd tea.Cmd
		m.freshManager.AccountInput, inputCmd = m.freshManager.AccountInput.Update(msg)
		return m, inputCmd
	}
}

func (m *Model) viewFreshView() string {
	return m.freshManager.View()
}
