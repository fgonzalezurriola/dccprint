package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fgonzalezurriola/dccprint/internal/theme"
)

type PrintView struct {
	pdfs         []string
	cursor       int
	selectedItem string
	theme        *theme.Theme
	width        int
	height       int
}

func NewPrintView(pdfs []string, theme *theme.Theme) PrintView {
	return PrintView{
		pdfs:  pdfs,
		theme: theme,
	}
}

func (s PrintView) Init() tea.Cmd {
	return nil
}

func (s PrintView) View() string {
	if len(s.pdfs) == 0 {
		return "No PDF files found in the current directory."
	}

	view := ""
	for i, pdf := range s.pdfs {
		cursor := " "
		textStyle := lipgloss.NewStyle().Foreground(s.theme.Unselected)

		if s.cursor == i {
			cursor = lipgloss.NewStyle().Foreground(s.theme.Selected).Render(">")
			textStyle = lipgloss.NewStyle().Foreground(s.theme.Selected)
		}
		view += lipgloss.NewStyle().Width(s.width).Render(
			lipgloss.JoinHorizontal(lipgloss.Left,
				cursor,
				" ",
				textStyle.Render(pdf),
			),
		) + "\n"
	}
	return view
}

func (s PrintView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if s.cursor > 0 {
				s.cursor--
			}
		case "down", "j":
			if s.cursor < len(s.pdfs)-1 {
				s.cursor++
			}
		case "enter":
			if len(s.pdfs) > 0 {
				s.selectedItem = s.pdfs[s.cursor]
			}
		case "ctrl+c", "q":
			return s, tea.Quit
		}
	}
	return s, nil
}

func (s *PrintView) SelectedItem() string {
	return s.selectedItem
}

func (s *PrintView) Reset() {
	s.selectedItem = ""
	s.cursor = 0
}

func (s *PrintView) SetTheme(theme *theme.Theme) {
	s.theme = theme
}

func (s *PrintView) SetSize(width, height int) {
	s.width = width
	s.height = height
}
