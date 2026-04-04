package ui

import (
	"journaltui/app"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func SetupMainMenu(state *app.AppState) {
	menu := tview.NewList().
		AddItem("New Entry", "", 0, nil).
		AddItem("Read / Edit Entries", "", 0, nil).
		AddItem("Quit", "", 0, nil)

	menu.SetBorder(false)
	menu.SetHighlightFullLine(true)
	menu.SetWrapAround(true)
	menu.SetBackgroundColor(ColorBackground)
	menu.SetMainTextColor(ColorText)
	menu.SetSelectedTextColor(ColorBackground)
	menu.SetSelectedBackgroundColor(ColorAccent)
	menu.SetShortcutColor(ColorAccent)

	menu.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		switch index {
		case 0:
			// OpenEditor(state, "", false)
			existingID := findTodaysEntry(state)
			if existingID != "" {
				OpenEditor(state, existingID, false)
			} else {
				OpenEditor(state, "", false)
			}
		case 1:
			OpenTable(state)
		case 2:
			state.TviewApp.Stop()
		}
	})

	menu.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			return nil
		}
		return event
	})

	root := tview.NewFlex().SetDirection(tview.FlexRow)
	root.SetBackgroundColor(ColorBackground)

	root.AddItem(tview.NewBox().SetBackgroundColor(ColorBackground), 0, 1, false)

	asciiArt := `       _                              __________  ______
      (_)___  __  ___________  ____ _/ /_  __/ / / /  _/
     / / __ \/ / / / ___/ __ \/ __ ' / / / / / / / // /  
    / / /_/ / /_/ / /  / / / / /_/ / / / / / /_/ // /   
 __/ /\____/\__,_/_/  /_/ /_/\__,_/_/ /_/  \____/___/   
/___/                                                   `
	title := tview.NewTextView().
		SetText(asciiArt).
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(false)
	title.SetTextColor(ColorAccent)
	title.SetBackgroundColor(ColorBackground)
	root.AddItem(title, 6, 0, false)

	// title := tview.NewTextView().
	// 	SetText("JournalTUI").
	// 	SetTextAlign(tview.AlignCenter).
	// 	SetDynamicColors(false)
	// title.SetTextColor(ColorAccent)
	// title.SetBackgroundColor(ColorBackground)
	// root.AddItem(title, 1, 0, false)

	root.AddItem(tview.NewBox().SetBackgroundColor(ColorBackground), 1, 0, false)

	hFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	hFlex.SetBackgroundColor(ColorBackground)
	hFlex.AddItem(tview.NewBox().SetBackgroundColor(ColorBackground), 0, 1, false)
	hFlex.AddItem(menu, 30, 0, true)
	hFlex.AddItem(tview.NewBox().SetBackgroundColor(ColorBackground), 0, 1, false)

	root.AddItem(hFlex, 8, 0, true)

	root.AddItem(tview.NewBox().SetBackgroundColor(ColorBackground), 0, 1, false)

	state.Pages.AddPage("main_menu", root, true, true)
	state.Pages.SetBackgroundColor(ColorBackground)
}

func findTodaysEntry(state *app.AppState) string {
	today := time.Now().Format("02-01-2006")
	for id, entry := range state.Entries {
		if entry.Date == today {
			return id
		}
	}
	return ""
}
