package components

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fgonzalezurriola/dccprint/internal/theme"
)

type PrintView struct {
	pdfs          []string
	cursor        int
	selectedItem  string
	theme         *theme.Theme
	width         int
	height        int
	StatusMessage string
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
	var currDir, err = os.Getwd()
	if err != nil {
		log.Fatalf("Error in PrintView View(): %v", err)
	}

	if len(s.pdfs) == 0 {
		return fmt.Sprintf("PDFs no encontrados en %s", currDir)
	}

	var lines []string
	if s.StatusMessage != "" {
		msgStyle := lipgloss.NewStyle().Foreground(s.theme.Selected).Bold(true)
		lines = append(lines, msgStyle.Render(s.StatusMessage))
	}
	for i, pdf := range s.pdfs {
		cursor := " "
		textStyle := lipgloss.NewStyle().Foreground(s.theme.Unselected)

		if s.cursor == i {
			cursor = lipgloss.NewStyle().Foreground(s.theme.Selected).Render(">")
			textStyle = lipgloss.NewStyle().Foreground(s.theme.Selected)
		}
		line := lipgloss.JoinHorizontal(lipgloss.Left,
			cursor,
			" ",
			textStyle.Render(pdf),
		)
		lines = append(lines, line)
	}
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
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
