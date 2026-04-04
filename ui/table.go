package ui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"journaltui/app"
	"journaltui/storage"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func SetupTable(state *app.AppState) {

	activeTitle := ""
	activeDate := ""

	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false).
		SetSelectedStyle(tcell.StyleDefault.
			Background(ColorAccent).
			Foreground(tcell.ColorWhite))

	hint := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText("[yellow]↑↓[#ffffe0] navigate  [yellow]enter[#ffffe0] open  [yellow]t[#ffffe0] rename  [yellow]d[#ffffe0] delete  [yellow]f[#ffffe0] find  [yellow]c[#ffffe0] clear filter  [yellow]esc[#ffffe0] back")

	tableBox := tview.NewFlex().SetDirection(tview.FlexRow)
	tableBox.SetBorder(true).
		SetTitle(" Entries ").
		SetTitleAlign(tview.AlignCenter)

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

	// populateTable := func() {
	// 	table.Clear()

	// 	headers := []string{"Title", "Date", "Last Edited"}
	// 	for col, h := range headers {
	// 		cell := tview.NewTableCell(h).
	// 			// SetTextColor(tcell.ColorTeal).
	// 			SetTextColor(ColorAccent).
	// 			SetSelectable(false).
	// 			SetExpansion(1)
	// 		table.SetCell(0, col, cell)
	// 	}

	// 	row := 1
	// 	for id, entry := range state.Entries {
	// 		// titleCell := tview.NewTableCell(entry.Title).SetExpansion(1)
	// 		titleCell := tview.NewTableCell(entry.Title).SetExpansion(1).SetReference(id)
	// 		dateCell := tview.NewTableCell(entry.Date).SetExpansion(1)
	// 		lastEditedCell := tview.NewTableCell(formatRelativeTime(entry.LastEdited)).SetExpansion(1)

	// 		table.SetCell(row, 0, titleCell)
	// 		table.SetCell(row, 1, dateCell)
	// 		table.SetCell(row, 2, lastEditedCell)
	// 		row++
	// 	}

	// 	table.ScrollToBeginning()
	// }

	populateTable := func() {
		table.Clear()

		headers := []string{"Title", "Date", "Last Edited"}
		for col, h := range headers {
			cell := tview.NewTableCell(h).
				SetTextColor(ColorAccent).
				SetSelectable(false).
				SetExpansion(1)
			table.SetCell(0, col, cell)
		}

		row := 1
		for id, entry := range state.Entries {
			if !entryMatchesFilters(entry, activeTitle, activeDate) {
				continue
			}
			titleCell := tview.NewTableCell(entry.Title).SetExpansion(1).SetReference(id)
			dateCell := tview.NewTableCell(entry.Date).SetExpansion(1)
			lastEditedCell := tview.NewTableCell(formatRelativeTime(entry.LastEdited)).SetExpansion(1)

			table.SetCell(row, 0, titleCell)
			table.SetCell(row, 1, dateCell)
			table.SetCell(row, 2, lastEditedCell)
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
			// entryID := table.GetCell(row, 0).Text
			entryID := table.GetCell(row, 0).GetReference().(string)
			closeTable(state)
			OpenEditor(state, entryID, true)
			return nil

		case tcell.KeyEscape:
			closeTable(state)
			return nil
		}

		switch event.Rune() {

		case 'f':
			ShowSearchPrompt(state,
				func(titleQuery, dateQuery string) {
					activeTitle = titleQuery
					activeDate = dateQuery
					populateTable()
					state.TviewApp.SetFocus(table)
				},
				func() {
					state.TviewApp.SetFocus(table)
				},
			)
			return nil

		case 'c':
			activeTitle = ""
			activeDate = ""
			populateTable()
			return nil

		case 't':
			row, _ := table.GetSelection()
			if row == 0 {
				return nil
			}
			// entryID := table.GetCell(row, 0).Text
			entryID := table.GetCell(row, 0).GetReference().(string)
			showRenamePrompt(state, entryID, table, populateTable)
			return nil

		case 'd':
			row, _ := table.GetSelection()
			if row == 0 {
				return nil
			}
			// entryID := table.GetCell(row, 0).Text
			entryID := table.GetCell(row, 0).GetReference().(string)
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

// parseDateQuery takes a natural language date string and returns
// matching day, month, year as ints. -1 means not specified.
func parseDateQuery(query string) (day, month, year int) {
	day, month, year = -1, -1, -1
	query = strings.TrimSpace(strings.ToLower(query))
	if query == "" {
		return
	}

	monthNames := map[string]int{
		"january": 1, "jan": 1,
		"february": 2, "feb": 2,
		"march": 3, "mar": 3,
		"april": 4, "apr": 4,
		"may":  5,
		"june": 6, "jun": 6,
		"july": 7, "jul": 7,
		"august": 8, "aug": 8,
		"september": 9, "sep": 9, "sept": 9,
		"october": 10, "oct": 10,
		"november": 11, "nov": 11,
		"december": 12, "dec": 12,
	}

	// Try dd-mm-yyyy format first
	parts := strings.Split(query, "-")
	if len(parts) == 3 {
		d, err1 := strconv.Atoi(parts[0])
		m, err2 := strconv.Atoi(parts[1])
		y, err3 := strconv.Atoi(parts[2])
		if err1 == nil && err2 == nil && err3 == nil {
			return d, m, y
		}
	}

	// Split by space and try to identify tokens
	tokens := strings.Fields(query)
	for _, token := range tokens {
		// Try as year (4 digits)
		if len(token) == 4 {
			if y, err := strconv.Atoi(token); err == nil {
				year = y
				continue
			}
		}
		// Try as day (1-2 digits)
		if n, err := strconv.Atoi(token); err == nil {
			day = n
			continue
		}
		// Try as month name
		if m, ok := monthNames[token]; ok {
			month = m
			continue
		}
	}

	return
}

// entryMatchesFilters returns true if the entry matches the title and date query
func entryMatchesFilters(entry app.EntryMeta, titleQuery, dateQuery string) bool {
	// Title filter
	if titleQuery != "" {
		if !strings.Contains(strings.ToLower(entry.Title), strings.ToLower(titleQuery)) {
			return false
		}
	}

	// Date filter
	if dateQuery != "" {
		filterDay, filterMonth, filterYear := parseDateQuery(dateQuery)

		// Parse entry date — stored as dd-mm-yyyy
		parts := strings.Split(entry.Date, "-")
		if len(parts) != 3 {
			return false
		}
		entryDay, _ := strconv.Atoi(parts[0])
		entryMonth, _ := strconv.Atoi(parts[1])
		entryYear, _ := strconv.Atoi(parts[2])

		if filterDay != -1 && filterDay != entryDay {
			return false
		}
		if filterMonth != -1 && filterMonth != entryMonth {
			return false
		}
		if filterYear != -1 && filterYear != entryYear {
			return false
		}
	}

	return true
}
