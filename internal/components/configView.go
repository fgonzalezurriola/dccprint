package components

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fgonzalezurriola/dccprint/internal/config"
	"github.com/fgonzalezurriola/dccprint/internal/theme"
)

type ConfigFinishedMsg struct{}

type configStep int

const (
	stepPrinter configStep = iota
	stepMode
	stepBorder
)

type ConfigView struct {
	currentStep configStep
	printerMenu Menu
	modeMenu    Menu
	borderMenu  Menu
	theme       *theme.Theme
	width       int
	height      int
}

func NewConfigView(t *theme.Theme) *ConfigView {
	printerMenuItems := []string{"Salita", "Toqui"}
	modeMenuItems := []string{"Simple", "Doble"}
	borderMenuItems := []string{"Corto", "Largo"}

	return &ConfigView{
		currentStep: stepPrinter,
		printerMenu: NewMenu(printerMenuItems, t),
		modeMenu:    NewMenu(modeMenuItems, t),
		borderMenu:  NewMenu(borderMenuItems, t),
		theme:       t,
	}
}

func (v *ConfigView) Init() tea.Cmd {
	return nil
}

func (v *ConfigView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
		switch v.currentStep {
		case stepPrinter:
			selectedPrinter := v.printerMenu.SelectedItem()
			if err := config.SavePrinter(selectedPrinter); err != nil {
				log.Printf("could not save printer config: %v", err)
			}
			v.currentStep = stepMode
			return v, nil
		case stepMode:
			selectedMode := v.modeMenu.SelectedItem()
			if err := config.SaveMode(selectedMode); err != nil {
				log.Printf("could not save mode config: %v", err)
			}
			if selectedMode == "Simple" {
				// Salta el borde y termina
				return v, func() tea.Msg { return ConfigFinishedMsg{} }
			} else {
				v.currentStep = stepBorder
			}
			return v, nil
		case stepBorder:
			selectedBorder := v.borderMenu.SelectedItem()
			if err := config.SaveBorder(selectedBorder); err != nil {
				log.Printf("could not save border config: %v", err)
			}
			// Termina el flujo
			return v, func() tea.Msg { return ConfigFinishedMsg{} }
		}
	}

	// Delega el update al menú actual
	switch v.currentStep {
	case stepPrinter:
		newMenu, menuCmd := v.printerMenu.Update(msg)
		v.printerMenu = newMenu.(Menu)
		cmd = menuCmd
	case stepMode:
		newMenu, menuCmd := v.modeMenu.Update(msg)
		v.modeMenu = newMenu.(Menu)
		cmd = menuCmd
	case stepBorder:
		newMenu, menuCmd := v.borderMenu.Update(msg)
		v.borderMenu = newMenu.(Menu)
		cmd = menuCmd
	}

	return v, cmd
}

func (v *ConfigView) View() string {
	switch v.currentStep {
	case stepPrinter:
		return v.printerMenu.View()
	case stepMode:
		return v.modeMenu.View()
	case stepBorder:
		return v.borderMenu.View()
	default:
		return ""
	}
}

func (v *ConfigView) SetSize(width, height int) {
	v.width = width
	v.height = height
	v.printerMenu.SetSize(width, height)
	v.modeMenu.SetSize(width, height)
	v.borderMenu.SetSize(width, height)
}

func (v *ConfigView) Reset() {
	v.currentStep = stepPrinter
	v.printerMenu.Reset()
	v.modeMenu.Reset()
	v.borderMenu.Reset()
}

func (v *ConfigView) SetTheme(theme *theme.Theme) {
	v.theme = theme
	v.printerMenu.SetTheme(theme)
	v.modeMenu.SetTheme(theme)
	v.borderMenu.SetTheme(theme)
}