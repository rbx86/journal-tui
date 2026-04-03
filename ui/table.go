package ui

import (
	"fmt"
	"time"

	"journaltui/app"
	"journaltui/storage"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func SetupTable(state *app.AppState) {
	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false)

	hint := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText("[yellow]↑↓[white] navigate  [yellow]enter[white] open  [yellow]t[white] rename  [yellow]d[white] delete  [yellow]esc[white] back")

	tableBox := tview.NewFlex().SetDirection(tview.FlexRow)
	tableBox.SetBorder(true).
		SetTitle(" Entries ").
		SetTitleAlign(tview.AlignLeft)

	tableBox.
		AddItem(table, 0, 1, true).
		AddItem(hint, 1, 0, false)

	modal := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false).
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(tableBox, 0, 17, true).
				AddItem(nil, 0, 1, false),
			0, 17, true).
		AddItem(nil, 0, 1, false)

	populateTable := func() {
		table.Clear()

		headers := []string{"ID", "Title", "Date", "Last Edited"}
		for col, h := range headers {
			cell := tview.NewTableCell(h).
				SetTextColor(tcell.ColorTeal).
				SetSelectable(false).
				SetExpansion(1)
			table.SetCell(0, col, cell)
		}

		row := 1
		for id, entry := range state.Entries {
			idCell := tview.NewTableCell(id).SetExpansion(1)
			titleCell := tview.NewTableCell(entry.Title).SetExpansion(1)
			dateCell := tview.NewTableCell(entry.Date).SetExpansion(1)
			lastEditedCell := tview.NewTableCell(formatRelativeTime(entry.LastEdited)).SetExpansion(1)

			table.SetCell(row, 0, idCell)
			table.SetCell(row, 1, titleCell)
			table.SetCell(row, 2, dateCell)
			table.SetCell(row, 3, lastEditedCell)
			row++
		}

		table.ScrollToBeginning()
	}

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			row, _ := table.GetSelection()
			if row == 0 {
				return nil
			}
			entryID := table.GetCell(row, 0).Text
			closeTable(state)
			OpenEditor(state, entryID, true)
			return nil

		case tcell.KeyEscape:
			closeTable(state)
			return nil
		}

		switch event.Rune() {
		case 't':
			row, _ := table.GetSelection()
			if row == 0 {
				return nil
			}
			entryID := table.GetCell(row, 0).Text
			showRenamePrompt(state, entryID, table, populateTable)
			return nil

		case 'd':
			row, _ := table.GetSelection()
			if row == 0 {
				return nil
			}
			entryID := table.GetCell(row, 0).Text
			ShowDeleteConfirm(state,
				func() {
					deleteEntry(state, entryID)
					populateTable()
					state.TviewApp.SetFocus(table)
				},
				func() {
					state.TviewApp.SetFocus(table)
				},
			)
			return nil
		}

		return event
	})

	state.Pages.AddPage("entry_table", modal, true, false)

	state.TableWidget = table
	state.PopulateTable = populateTable
}

func OpenTable(state *app.AppState) {
	state.PopulateTable()
	state.Pages.ShowPage("entry_table")
	state.TviewApp.SetFocus(state.TableWidget)
}

func closeTable(state *app.AppState) {
	state.Pages.HidePage("entry_table")
	state.Pages.SwitchToPage("main_menu")
}

func deleteEntry(state *app.AppState, entryID string) {
	storage.DeleteEntry(entryID)
	delete(state.Entries, entryID)
	storage.SaveMetadata(state.Entries)
}

func showRenamePrompt(state *app.AppState, entryID string, table *tview.Table, repopulate func()) {
	ShowTitlePrompt(state,
		func(newTitle string) {
			entry := state.Entries[entryID]
			entry.Title = newTitle
			entry.LastEdited = time.Now().UTC().Format(time.RFC3339)
			state.Entries[entryID] = entry
			storage.SaveMetadata(state.Entries)
			repopulate()
			state.TviewApp.SetFocus(table)
		},
		func() {
			state.TviewApp.SetFocus(table)
		},
	)
}

func formatRelativeTime(timestamp string) string {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return "unknown"
	}

	diff := time.Since(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		mins := int(diff.Minutes())
		return fmt.Sprintf("%d min. ago", mins)
	case diff < 24*time.Hour:
		hrs := int(diff.Hours())
		return fmt.Sprintf("%d hr. ago", hrs)
	case diff < 30*24*time.Hour:
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%d days ago", days)
	case diff < 365*24*time.Hour:
		months := int(diff.Hours() / 24 / 30)
		return fmt.Sprintf("%d months ago", months)
	default:
		years := int(diff.Hours() / 24 / 365)
		return fmt.Sprintf("%d years ago", years)
	}
}
