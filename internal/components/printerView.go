package components

import (
	"github.com/fgonzalezurriola/dccprint/internal/theme"
)

type PrinterView struct {
	Menu
}

func NewPrinterView(items []string, theme *theme.Theme) PrinterView {
	return PrinterView{Menu: NewMenu(items, theme)}
}
