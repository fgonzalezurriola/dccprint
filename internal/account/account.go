package account

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/fgonzalezurriola/dccprint/internal/config"
	"github.com/fgonzalezurriola/dccprint/internal/theme"
)

type Manager struct {
	AccountInput textinput.Model
}

func NewManager(t *theme.Theme, cfg config.Config) Manager {
	ti := textinput.New()
	ti.Placeholder = "Ingresa el nombre de cuenta (sin @)"
	ti.Focus()
	ti.CharLimit = 16
	ti.SetValue(cfg.Account)
	ti.PromptStyle = lipgloss.NewStyle().Foreground(t.Selected)
	ti.TextStyle = lipgloss.NewStyle().Foreground(t.Header)

	return Manager{AccountInput: ti}
}

func (a *Manager) SaveAccount() {
	account := a.AccountInput.Value()
	config.SaveAccount(account)
}
