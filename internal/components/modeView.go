package components

import (
	"github.com/fgonzalezurriola/dccprint/internal/theme"
)

type ModeView struct {
	Menu
}

func NewModeView(items []string, theme *theme.Theme) ModeView {
	return ModeView{Menu: NewMenu(items, theme)}
}
