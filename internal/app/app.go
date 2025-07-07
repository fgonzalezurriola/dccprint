package app

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/fgonzalezurriola/dccprint/internal/components"
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
	currentView  viewState
	mainMenu     components.Menu
	printView    components.Menu
	configView   components.Menu
	themeMenu    components.Menu
	theme        *theme.Theme
	width        int
	height       int
	accountInput textinput.Model
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
		currentView: mainView,
		mainMenu:    components.NewMenu(mainMenuItems, t),

		themeMenu:    components.NewMenu(themeMenuItems, t),
		theme:        t,
		accountInput: ti,
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
				// Todo
			case "Configurar Cuenta":
				m.currentView = accountView
				m.accountInput.Focus()
			case "Cambiar Theme":
				m.themeMenu.Reset()
				m.currentView = themeView
			case "Salir":
				return m, tea.Quit
			}
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
			m.accountInput.PromptStyle = lipgloss.NewStyle().Foreground(m.theme.Selected)
			m.accountInput.TextStyle = lipgloss.NewStyle().Foreground(m.theme.Header)
			m.mainMenu.Reset()
			m.currentView = mainView
		}

	case accountView:
		if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
			account := m.accountInput.Value()
			config.SaveAccount(account)
			m.mainMenu.Reset()
			m.currentView = mainView
		} else {
			var inputCmd tea.Cmd
			m.accountInput, inputCmd = m.accountInput.Update(msg)
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
		view = m.mainMenu.View()
	case configView:
		view = m.mainMenu.View()
	case accountView:
		view = m.accountInput.View()
	case themeView:
		view = m.themeMenu.View()
	}

	acc := config.Load().Account
	if len(acc) < 3 {
		acc = "Ingresa tu usuario DCC en \"Configurar Cuenta\""
	} else {
		acc = "Cuenta configurada: " + acc
	}

	accountText := lipgloss.NewStyle().Render(acc)
	content := lipgloss.JoinVertical(lipgloss.Left, header, view, accountText)
	centeredContent := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(content)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		centeredContent,
	)
}
