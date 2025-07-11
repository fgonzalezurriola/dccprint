package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/fgonzalezurriola/dccprint/internal/theme"
)

const Logo = `
 ██████╗   ██████╗  ██████╗   ██████╗  ██████╗  ██╗ ███╗   ██╗ ████████╗
 ██╔══██╗ ██╔════╝ ██╔════╝   ██╔══██╗ ██╔══██╗ ██║ ████╗  ██║ ╚══██╔══╝
 ██║  ██║ ██║      ██║        ██████╔╝ ██████╔╝ ██║ ██╔██╗ ██║    ██║
 ██║  ██║ ██║      ██║        ██╔═══╝  ██╔══██╗ ██║ ██║╚██╗██║    ██║
 ██████╔╝ ╚██████╗ ╚██████╗   ██║      ██║  ██║ ██║ ██║ ╚████║    ██║
 ╚═════╝   ╚═════╝  ╚═════╝   ╚═╝      ╚═╝  ╚═╝ ╚═╝ ╚═╝  ╚═══╝    ╚═╝
`
const LogoWidth = 73

func RenderHeader(width int, theme *theme.Theme) string {

	baseStyle := lipgloss.NewStyle().
		Background(theme.Selected).
		Foreground(theme.Header)
	headerStyle := baseStyle
	var header string
	var headerMessage string
	if width > LogoWidth {
		headerMessage = Logo
		headerStyle = baseStyle.Padding(1, 1)
	} else {
		headerMessage = "DCC PRINT"
		headerStyle = baseStyle.Padding(1, 8)
	}
	header = headerStyle.Render(headerMessage)

	line := lipgloss.NewStyle().Background(theme.Background).Height(1).Render("")

	return lipgloss.JoinVertical(lipgloss.Left, header, line)
}
