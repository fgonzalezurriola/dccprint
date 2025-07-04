package app

import (
	"github.com/fgonzalezurriola/dccprint/internal/components"
	"github.com/fgonzalezurriola/dccprint/internal/config"
	"github.com/fgonzalezurriola/dccprint/internal/theme"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	mainView viewState = iota
	themeView
)

type viewState int

type Model struct {
	currentView viewState
	mainMenu    components.Menu
	themeMenu   components.Menu
	theme       *theme.Theme
	width       int
	height      int
}

func New() Model {
	cfg := config.Load()
	t := theme.New(cfg.Theme)

	mainMenuItems := []string{"Imprimir", "Cuenta", "Theme", "Salir"}
	themeMenuItems := []string{"Default", "Cadcc", "Anakena"}

	//printerMenuItems := []string{"Salita", "Toqui"}

	return Model{
		currentView: mainView,
		mainMenu:    components.NewMenu(mainMenuItems, t),
		themeMenu:   components.NewMenu(themeMenuItems, t),
		theme:       t,
	}

}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.mainMenu.SetSize(msg.Width, msg.Height)
		m.themeMenu.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.currentView == mainView {
				return m, tea.Quit
			}
		case "esc":
			if m.currentView != mainView {
				m.currentView = mainView
			}
		}
	}

	switch m.currentView {
	case mainView:
		newMenu, newCmd := m.mainMenu.Update(msg)
		m.mainMenu = newMenu.(components.Menu)
		cmd = newCmd
		switch m.mainMenu.SelectedItem() {
		case "Cambiar Theme":
			m.currentView = themeView
			m.mainMenu.Reset()
		case "Salir":
			return m, tea.Quit

		}

	case themeView:
		newMenu, newCmd := m.themeMenu.Update(msg)
		m.themeMenu = newMenu.(components.Menu)
		cmd = newCmd

		switch m.themeMenu.SelectedItem() {
		case "Default", "Cadcc", "Anakena":
			config.Save(m.themeMenu.SelectedItem())
			m.theme = theme.New(m.themeMenu.SelectedItem())
			m.mainMenu.SetTheme(m.theme)
			m.themeMenu.SetTheme(m.theme)
			m.currentView = mainView
			m.themeMenu.Reset()
		case "Volver":
			m.currentView = mainView
			m.themeMenu.Reset()
		}
	}

	return m, cmd
}

func (m Model) View() string {
	header := m.renderHeader()
	var view string

	switch m.currentView {
	case mainView:
		view = m.mainMenu.View()
	case themeView:
		view = m.themeMenu.View()
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

func (m Model) renderHeader() string {
	headerStyle := lipgloss.NewStyle().
		Background(m.theme.Header).
		Foreground(m.theme.Foreground).
		Padding(1, 2)

	title := headerStyle.Render("Titulo de prueba")

	line := lipgloss.NewStyle().
		Background(m.theme.Background).
		Height(1).
		Render("")

	return lipgloss.JoinVertical(lipgloss.Left, line, title, line)
}
