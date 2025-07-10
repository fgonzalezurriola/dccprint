package theme

import (
	"github.com/fgonzalezurriola/dccprint/internal/config"
)

type Manager struct {
	Current string
}

func NewManager(current string) *Manager {
	return &Manager{Current: current}
}

func (tm *Manager) ChangeTheme(selected string) {
	config.SaveTheme(selected)
	tm.Current = selected
}
