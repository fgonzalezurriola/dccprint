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
	Accent     lipgloss.Color
}

var themes = map[string]*Theme{
	"Default": {
		Selected:   lipgloss.Color("#3700ffff"), // Pink
		Unselected: lipgloss.Color("#71e7f7ff"), // Gray
		Foreground: lipgloss.Color("252"), // Light Gray
		Background: lipgloss.Color("235"), // Dark Gray
		Border:     lipgloss.Color("238"), // Darker Gray
		Header:     lipgloss.Color("#71e7f7ff"), // Pink
		Footer:     lipgloss.Color("#71e7f7ff"), // Pink
		Accent:     lipgloss.Color("ffffff"),
	},
	"Dcc...": {
		Selected:   lipgloss.Color("32"),  // Blue
		Unselected: lipgloss.Color("240"), // Mid grey
		Foreground: lipgloss.Color("228"), // Light yellow
		Background: lipgloss.Color("235"), // Gris oscuro
		Border:     lipgloss.Color("23"),  // Midnight green
		Header:     lipgloss.Color("228"),
		Footer:     lipgloss.Color("228"),
		Accent:     lipgloss.Color("49"), // Emerald
	},
	"Anakena": {
		Selected:   lipgloss.Color("#2f9c49"), // Anakena green
		Unselected: lipgloss.Color("#252"),    // Green
		Foreground: lipgloss.Color("#d00a05"), // Light gray // text
		Background: lipgloss.Color("234"),     // Very dark gray
		Border:     lipgloss.Color("34"),      // Leaf green
		Header:     lipgloss.Color("215"),     // Coral orange
		Footer:     lipgloss.Color("215"),     // Leaf green
		Accent:     lipgloss.Color("67"),      // Muted blue
	},
}

// function New returns a Theme by its name
// If the name is not found, returns "Default"
func New(name string) *Theme {
	if t, ok := themes[name]; ok {
		return t
	}
	return themes["Default"]
}
