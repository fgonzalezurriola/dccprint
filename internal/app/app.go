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

const (
	mainView viewState = iota
	printView
	configView
	accountView
	themeView
)

type viewState int

type Model struct {
	currentView    viewState
	mainMenu       components.Menu
	printView      components.PrintView
	configView     components.Menu
	themeMenu      components.Menu
	theme          *theme.Theme
	width          int
	height         int
	accountManager account.Manager
}

func NewModel() *Model {
	cfg := config.Load()
	t := theme.New(cfg.Theme)

	mainMenuItems := []string{"Imprimir PDF", "Configuración de Impresión", "Configurar Cuenta", "Cambiar Theme", "Salir"}
	themeMenuItems := []string{"Default", "Cadcc", "Anakena"}

	ti := textinput.New()
	ti.Placeholder = "Ingresa el nombre de cuenta (sin @)"
	ti.Focus()
	ti.CharLimit = 16
	ti.SetValue(cfg.Account)
	ti.PromptStyle = lipgloss.NewStyle().Foreground(t.Selected)
	ti.TextStyle = lipgloss.NewStyle().Foreground(t.Header)

	return &Model{
		currentView:    mainView,
		mainMenu:       components.NewMenu(mainMenuItems, t),
		printView:      components.NewPrintView(scripts.GetPDFFiles(), t),
		themeMenu:      components.NewMenu(themeMenuItems, t),
		theme:          t,
		accountManager: account.NewManager(t, cfg),
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle global messages first
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.mainMenu.SetSize(msg.Width, msg.Height)
		m.themeMenu.SetSize(msg.Width, msg.Height)
		m.printView.SetSize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		// Handle global keybindings
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			if m.currentView != mainView {
				m.mainMenu.Reset()
				m.currentView = mainView
				return m, nil
			}
		}
	}

	switch m.currentView {
	case mainView:
		// Update menu first
		newMenu, menuCmd := m.mainMenu.Update(msg)
		m.mainMenu = newMenu.(components.Menu)
		cmd = menuCmd

		if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
			switch m.mainMenu.SelectedItem() {
			case "Imprimir PDF":
				m.currentView = printView
			case "Configurar Cuenta":
				m.currentView = accountView
				m.accountManager.AccountInput.Focus()
			case "Cambiar Theme":
				m.themeMenu.Reset()
				m.currentView = themeView
			case "Salir":
				return m, tea.Quit
			}
		}

	case printView:
		newSelector, selectorCmd := m.printView.Update(msg)
		m.printView = newSelector.(components.PrintView)
		cmd = selectorCmd

		if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
			// TODO
			m.currentView = mainView
		}

	case themeView:
		newMenu, menuCmd := m.themeMenu.Update(msg)
		m.themeMenu = newMenu.(components.Menu)
		cmd = menuCmd

		if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
			selectedTheme := m.themeMenu.SelectedItem()
			config.SaveTheme(selectedTheme)
			m.theme = theme.New(selectedTheme)
			m.mainMenu.SetTheme(m.theme)
			m.themeMenu.SetTheme(m.theme)
			// Todo: Muy engorroso, abstraer
			m.accountManager.AccountInput.PromptStyle = lipgloss.NewStyle().Foreground(m.theme.Selected)
			m.accountManager.AccountInput.TextStyle = lipgloss.NewStyle().Foreground(m.theme.Header)
			m.mainMenu.Reset()
			m.currentView = mainView
		}

	case accountView:
		if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
			m.accountManager.SaveAccount()
			m.mainMenu.Reset()
			m.currentView = mainView
		} else {
			var inputCmd tea.Cmd
			m.accountManager.AccountInput, inputCmd = m.accountManager.AccountInput.Update(msg)
			cmd = inputCmd
		}
	}

	return m, cmd
}

func (m *Model) View() string {
	header := components.RenderHeader(m.width, m.theme)
	var view string

	switch m.currentView {
	case mainView:
		view = m.mainMenu.View()
	case printView:
		view = m.printView.View()
	case configView:
		view = m.mainMenu.View()
	case accountView:
		view = m.accountManager.AccountInput.View()
	case themeView:
		view = m.themeMenu.View()
	}

	accInfo := config.Load().Account
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
