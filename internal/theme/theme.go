package theme

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Selected   lipgloss.Color
	Unselected lipgloss.Color
	Foreground lipgloss.Color
	Background lipgloss.Color
	Border     lipgloss.Color
	Header     lipgloss.Color
	Footer     lipgloss.Color
}

// "Terminal" theme uses terminal's default foreground/background and ANSI colors
// Selected, Unselected, Header, etc. use basic ANSI colors for visibility
var themes = map[string]*Theme{
	"Default": {
		Selected:   lipgloss.Color("#ffc934ff"),
		Unselected: lipgloss.Color("#8d6f1dff"),
		Foreground: lipgloss.Color("#default"),
		Background: lipgloss.Color("#default"),
		Header:     lipgloss.Color("#430fffff"),
	},
	"Dcc...": {
		Selected:   lipgloss.Color("#00BCF0"),
		Unselected: lipgloss.Color("#5d8c99ff"),
		Foreground: lipgloss.Color("228"),
		Background: lipgloss.Color("#36454F"),
		Header:     lipgloss.Color("228"),
	},
	"Anakena": {
		Selected:   lipgloss.Color("#32a74eff"),
		Unselected: lipgloss.Color("#22883aff"),
		Foreground: lipgloss.Color("#d00a05"),
		Background: lipgloss.Color("#36454F"),
		Header:     lipgloss.Color("#ffa149ff"),
	},
}

// function New returns a Theme by its name
// If the name is not found, returns "Default"
// To use terminal colors, select theme "Terminal"
func New(name string) *Theme {
	if t, ok := themes[name]; ok {
		return t
	}
	return themes["Default"]
}
