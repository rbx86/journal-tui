package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	ColorBackground = tcell.NewRGBColor(18, 18, 18)
	ColorSurface    = tcell.NewRGBColor(28, 28, 28)
	ColorBorder     = tcell.NewRGBColor(50, 50, 50)
	ColorText       = tcell.NewRGBColor(220, 220, 220)
	ColorMuted      = tcell.NewRGBColor(120, 120, 120)
	ColorAccent     = tcell.NewRGBColor(29, 185, 84)
	ColorAccentDark = tcell.NewRGBColor(20, 130, 60)
)

func ApplyTheme() {
	tview.Styles.PrimitiveBackgroundColor = ColorBackground
	tview.Styles.ContrastBackgroundColor = ColorSurface
	tview.Styles.MoreContrastBackgroundColor = ColorSurface
	tview.Styles.BorderColor = ColorBorder
	tview.Styles.TitleColor = ColorAccent
	tview.Styles.GraphicsColor = ColorBorder
	tview.Styles.PrimaryTextColor = ColorText
	tview.Styles.SecondaryTextColor = ColorMuted
	tview.Styles.TertiaryTextColor = ColorMuted
	tview.Styles.InverseTextColor = ColorBackground
	tview.Styles.ContrastSecondaryTextColor = ColorMuted
}