package style

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

// Color palette
var (
	ColorBackground = lipgloss.Color("#0f1419")
	ColorSurface    = lipgloss.Color("#1c2128")
	ColorBorder     = lipgloss.Color("#30363d")
	ColorPrimary    = lipgloss.Color("#58a6ff")
	ColorAccent     = lipgloss.Color("#f78166")
	ColorSuccess    = lipgloss.Color("#3fb950")
	ColorDanger     = lipgloss.Color("#f85149")
	ColorWarning    = lipgloss.Color("#d29922")
	ColorMuted      = lipgloss.Color("#8b949e")
	ColorText       = lipgloss.Color("#c9d1d9")
	ColorBright     = lipgloss.Color("#ffffff")
)

// Text styles
var (
	StyleTitle = lipgloss.NewStyle().
			Foreground(ColorBright).
			Bold(true).
			Padding(0, 1)

	StyleHeader = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true)

	StyleLabel = lipgloss.NewStyle().
			Foreground(ColorMuted)

	StyleValue = lipgloss.NewStyle().
			Foreground(ColorText)

	StyleMuted = lipgloss.NewStyle().
			Foreground(ColorMuted)

	StyleSuccess = lipgloss.NewStyle().
			Foreground(ColorSuccess)

	StyleDanger = lipgloss.NewStyle().
			Foreground(ColorDanger)

	StyleWarning = lipgloss.NewStyle().
			Foreground(ColorWarning)

	// Box styles
	StyleBox = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(1, 2).
			Foreground(ColorText)

	StyleBoxSelected = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(ColorPrimary).
				Padding(1, 2).
				Foreground(ColorText)

	StyleBoxTitle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true)

	// Table styles
	StyleTableHeader = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Bold(true)

	StyleTableRow = lipgloss.NewStyle().
			Foreground(ColorText)

	StyleTableRowAlt = lipgloss.NewStyle().
			Foreground(ColorMuted)

	// Form styles
	StyleInput = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(0, 1).
			Foreground(ColorText)

	StyleInputFocused = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(ColorPrimary).
				Padding(0, 1).
				Foreground(ColorBright)

	StyleButton = lipgloss.NewStyle().
			Background(ColorPrimary).
			Foreground(ColorBright).
			Padding(0, 2).
			Margin(0, 1)

	StyleButtonDanger = lipgloss.NewStyle().
				Background(ColorDanger).
				Foreground(ColorBright).
				Padding(0, 2)

	StyleButtonMuted = lipgloss.NewStyle().
				Background(ColorSurface).
				Foreground(ColorMuted).
				Padding(0, 2)

	// Status badges
	StyleBadgeAdmin = lipgloss.NewStyle().
			Background(ColorAccent).
			Foreground(ColorBright).
			Padding(0, 1).
			Margin(0, 1)

	StyleBadgePlayer = lipgloss.NewStyle().
			Background(ColorPrimary).
			Foreground(ColorBright).
			Padding(0, 1).
			Margin(0, 1)

	StyleBadgeClassless = lipgloss.NewStyle().
			Background(ColorWarning).
			Foreground(ColorBright).
			Padding(0, 1).
			Margin(0, 1)

	// Divider
	StyleDivider = lipgloss.NewStyle().
			Foreground(ColorBorder)
)

// RenderDivider renders a horizontal divider
func RenderDivider(width int) string {
	div := ""
	for i := 0; i < width; i++ {
		div += "─"
	}
	return lipgloss.Style{}.Foreground(ColorBorder).Render(div)
}

// Success renders a success message
func Success(msg string) string {
	return StyleSuccess.Render("✓ ") + msg
}

// Error renders an error message
func Error(msg string) string {
	return StyleDanger.Render("✗ ") + msg
}

// Info renders an info message
func Info(msg string) string {
	return lipgloss.Style{}.Foreground(ColorPrimary).Render("ℹ ") + msg
}

// fatal prints an error and exits
func Fatal(msg string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", Error(msg), err)
	} else {
		fmt.Fprintln(os.Stderr, Error(msg))
	}
	os.Exit(1)
}
