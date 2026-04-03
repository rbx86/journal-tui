package ui

import (
	"fmt"
	"time"

	"journaltui/app"
	"journaltui/storage"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func SetupEditor(state *app.AppState) {
	textArea := tview.NewTextArea()

	hint := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	updateHint := func() {
		if state.IsReadMode {
			hint.SetText("[yellow]ctrl+e[white] edit  [yellow]esc[white] back")
		} else {
			hint.SetText("[yellow]ctrl+s[white] save  [yellow]esc[white] back")
		}
	}

	updateBorderTitle := func(editorBox *tview.Flex, title string) {
		if state.IsReadMode {
			editorBox.SetTitle(fmt.Sprintf(" %s [read] ", title))
		} else {
			editorBox.SetTitle(fmt.Sprintf(" %s [editing] ", title))
		}
	}

	editorBox := tview.NewFlex().SetDirection(tview.FlexRow)
	editorBox.SetBorder(true).
		SetTitleAlign(tview.AlignLeft)

	editorBox.
		AddItem(textArea, 0, 1, true).
		AddItem(hint, 1, 0, false)

	modal := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false).
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(editorBox, 0, 17, true).
				AddItem(nil, 0, 1, false),
			0, 17, true).
		AddItem(nil, 0, 1, false)

	textArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {

		case tcell.KeyCtrlS:
			if state.IsReadMode {
				return nil
			}
			handleSave(state, textArea, editorBox, updateBorderTitle)
			return nil

		case tcell.KeyCtrlE:
			if !state.IsReadMode {
				return nil
			}
			state.IsReadMode = false
			textArea.SetDisabled(false)
			updateHint()
			updateBorderTitle(editorBox, getEditorTitle(state))
			return nil

		case tcell.KeyEscape:
			if state.IsDirty {
				ShowUnsavedConfirm(state,
					func() {
						closeEditor(state)
					},
					func() {
						state.TviewApp.SetFocus(textArea)
					},
				)
			} else {
				closeEditor(state)
			}
			return nil
		}

		if !state.IsReadMode {
			state.IsDirty = true
			updateHint()
		}

		return event
	})

	state.Pages.AddPage("entry_editor", modal, true, false)

	state.EditorTextArea = textArea
	state.EditorBox = editorBox
	state.UpdateEditorHint = updateHint
	state.UpdateEditorBorderTitle = updateBorderTitle
}

func OpenEditor(state *app.AppState, entryID string, readMode bool) {
	state.CurrentEntryID = entryID
	state.IsReadMode = readMode
	state.IsDirty = false

	textArea := state.EditorTextArea
	editorBox := state.EditorBox

	if entryID == "" {
		// New entry
		state.IsNewEntry = true
		textArea.SetText("", true)
		textArea.SetDisabled(false)
		editorBox.SetTitle(" New Entry [editing] ")
	} else {
		// Existing entry
		state.IsNewEntry = false
		content, err := storage.LoadEntry(entryID)
		if err != nil {
			content = ""
		}
		textArea.SetText(content, true)
		textArea.SetDisabled(readMode)
		title := getEditorTitle(state)
		if readMode {
			editorBox.SetTitle(fmt.Sprintf(" %s [read] ", title))
		} else {
			editorBox.SetTitle(fmt.Sprintf(" %s [editing] ", title))
		}
	}

	state.UpdateEditorHint()
	state.Pages.ShowPage("entry_editor")
	state.TviewApp.SetFocus(state.EditorTextArea)
}

func handleSave(state *app.AppState, textArea *tview.TextArea, editorBox *tview.Flex, updateBorderTitle func(*tview.Flex, string)) {
	content := textArea.GetText()

	if state.IsNewEntry {
		entryID := time.Now().Format("20060102-150405")
		state.CurrentEntryID = entryID

		ShowTitlePrompt(state,
			func(title string) {
				saveEntry(state, entryID, title, content)
				state.IsNewEntry = false
				state.IsDirty = false
				updateBorderTitle(editorBox, title)
			},
			func() {
				saveEntry(state, entryID, "Untitled", content)
				state.IsNewEntry = false
				state.IsDirty = false
				updateBorderTitle(editorBox, "Untitled")
			},
		)
	} else {
		entry := state.Entries[state.CurrentEntryID]
		entry.LastEdited = time.Now().UTC().Format(time.RFC3339)
		state.Entries[state.CurrentEntryID] = entry

		if err := storage.SaveEntry(state.CurrentEntryID, content); err != nil {
			return
		}
		if err := storage.SaveMetadata(state.Entries); err != nil {
			return
		}
		state.IsDirty = false
	}
}

func saveEntry(state *app.AppState, entryID, title, content string) {
	now := time.Now()

	entriesDir, err := storage.EntriesDir()
	if err != nil {
		return
	}

	meta := app.EntryMeta{
		Title:      title,
		Date:       now.Format("02-01-2006"),
		CreatedAt:  now.UTC().Format(time.RFC3339),
		LastEdited: now.UTC().Format(time.RFC3339),
		Path:       entriesDir + "/" + entryID + ".md",
	}

	state.Entries[entryID] = meta

	storage.SaveEntry(entryID, content)
	storage.SaveMetadata(state.Entries)
}

func closeEditor(state *app.AppState) {
	state.CurrentEntryID = ""
	state.IsDirty = false
	state.IsNewEntry = true
	state.IsReadMode = false
	state.Pages.HidePage("entry_editor")
	state.Pages.SwitchToPage("main_menu")
}

func getEditorTitle(state *app.AppState) string {
	if entry, ok := state.Entries[state.CurrentEntryID]; ok {
		return entry.Title
	}
	return "Untitled"
}
