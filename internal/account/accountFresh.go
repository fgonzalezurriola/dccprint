package account

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/fgonzalezurriola/dccprint/internal/config"
	"github.com/fgonzalezurriola/dccprint/internal/theme"
)

type FreshManager struct {
	AccountInput textinput.Model
}

func NewFreshManager(t *theme.Theme) FreshManager {
	ti := textinput.New()
	ti.Placeholder = "Configura tu cuenta (sin @)"
	ti.Focus()
	ti.CharLimit = 16
	ti.PromptStyle = lipgloss.NewStyle().Foreground(t.Selected)
	ti.TextStyle = lipgloss.NewStyle().Foreground(t.Header)

	return FreshManager{AccountInput: ti}
}

func (a *FreshManager) SaveAccount() {
	account := a.AccountInput.Value()
	config.SaveAccount(account)
}

func (a *FreshManager) View() string {
	bold := lipgloss.NewStyle().Bold(true).Render("BIENVENIDO A DCCPRINT")
	info := lipgloss.NewStyle().Render("Primero, el nombre de tu cuenta DCC (lo que va antes del @)")
	input := a.AccountInput.View()
	return lipgloss.JoinVertical(lipgloss.Left, bold, info, "", input)
}
