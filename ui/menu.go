package ui

import (
	"journaltui/app"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func SetupMainMenu(state *app.AppState) {
	outerFlex := tview.NewFlex().
		SetDirection(tview.FlexRow)

	title := tview.NewTextView().
		SetText("JournalTUI").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	title.SetTextColor(tcell.ColorTeal)

	menu := tview.NewList().
		AddItem("New Entry", "", 0, nil).
		AddItem("Read / Edit Entries", "", 0, nil).
		AddItem("Quit", "", 0, nil)

	menu.SetBorder(false)
	menu.SetHighlightFullLine(true)
	menu.SetWrapAround(true)

	menu.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		switch index {
		case 0:
			OpenEditor(state, "", false)
		case 1:
			OpenTable(state)
		case 2:
			state.TviewApp.Stop()
		}
	})

	menu.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			return nil // swallow the event
		}
		return event
	})

	centerFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false).
		AddItem(menu, 30, 0, true).
		AddItem(nil, 0, 1, false)

	outerFlex.
		AddItem(nil, 0, 1, false).
		AddItem(title, 1, 0, false).
		AddItem(tview.NewBox(), 1, 0, false).
		AddItem(centerFlex, 5, 0, true).
		AddItem(nil, 0, 1, false)

	state.Pages.AddPage("main_menu", outerFlex, true, true)
}
