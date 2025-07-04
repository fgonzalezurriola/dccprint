package components

import (
	"fmt"

	"github.com/fgonzalezurriola/dccprint/internal/theme"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Menu struct {
	items        []string
	cursor       int
	selectedItem string
	theme        *theme.Theme
	width        int
	height       int
}

func NewMenu(items []string, theme *theme.Theme) Menu {
	return Menu{
		items: items,
		theme: theme,
	}
}

func (m Menu) Init() tea.Cmd {
	return nil
}

func (m Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "enter":
			m.selectedItem = m.items[m.cursor]
		}
	}
	return m, nil
}

func (m Menu) View() string {
	s := ""
	for i, item := range m.items {
		cursor := " "
		textStyle := lipgloss.NewStyle().Foreground(m.theme.Unselected)

		if m.cursor == i {
			cursor = lipgloss.NewStyle().Foreground(m.theme.Selected).Render(">")
			textStyle = lipgloss.NewStyle().Foreground(m.theme.Selected)
		}
		s += fmt.Sprintf("%s %s\n", cursor, textStyle.Render(item))
	}
	return s
}

func (m *Menu) SelectedItem() string {
	return m.selectedItem
}

func (m *Menu) Reset() {
	m.selectedItem = ""
	m.cursor = 0
}

func (m *Menu) SetTheme(theme *theme.Theme) {
	m.theme = theme
}

func (m *Menu) SetSize(width, height int) {
	m.width = width
	m.height = height
}
