package ui

import (
	"journaltui/app"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func ShowTitlePrompt(state *app.AppState, onConfirm func(title string), onCancel func()) {
	input := tview.NewInputField().
		SetLabel("Entry title: ").
		SetFieldWidth(40)

	input.SetBorder(true).
		SetTitle(" Title ").
		SetTitleAlign(tview.AlignLeft)

	input.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			title := input.GetText()
			if title == "" {
				title = "Untitled"
			}
			state.Pages.RemovePage("title_prompt")
			state.TviewApp.SetFocus(state.Pages)
			onConfirm(title)
		case tcell.KeyEscape:
			state.Pages.RemovePage("title_prompt")
			state.TviewApp.SetFocus(state.Pages)
			onCancel()
		}
	})

	modal := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false).
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(input, 5, 0, true).
				AddItem(nil, 0, 1, false),
			50, 0, true,
		).
		AddItem(nil, 0, 1, false)

	state.Pages.AddPage("title_prompt", modal, true, true)
	state.TviewApp.SetFocus(input)
}

func ShowDeleteConfirm(state *app.AppState, onConfirm func(), onCancel func()) {
	text := tview.NewTextView().
		SetText("Delete this entry? (y/n)").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	box := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(text, 1, 0, false).
		AddItem(nil, 0, 1, false)

	box.SetBorder(true).
		SetTitle(" Delete ").
		SetTitleAlign(tview.AlignLeft)

	box.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'y', 'Y':
			state.Pages.RemovePage("delete_confirm")
			state.TviewApp.SetFocus(state.Pages)
			onConfirm()
			return nil
		case 'n', 'N':
			state.Pages.RemovePage("delete_confirm")
			state.TviewApp.SetFocus(state.Pages)
			onCancel()
			return nil
		}
		if event.Key() == tcell.KeyEscape {
			state.Pages.RemovePage("delete_confirm")
			state.TviewApp.SetFocus(state.Pages)
			onCancel()
			return nil
		}
		return event
	})

	modal := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false).
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(box, 5, 0, true).
				AddItem(nil, 0, 1, false),
			40, 0, true,
		).
		AddItem(nil, 0, 1, false)

	state.Pages.AddPage("delete_confirm", modal, true, true)
	state.TviewApp.SetFocus(box)
}

func ShowUnsavedConfirm(state *app.AppState, onConfirm func(), onCancel func()) {
	text := tview.NewTextView().
		SetText("Unsaved changes. Quit anyway? (y/n)").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	box := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(text, 1, 0, false).
		AddItem(nil, 0, 1, false)

	box.SetBorder(true).
		SetTitle(" Unsaved Changes ").
		SetTitleAlign(tview.AlignLeft)

	box.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'y', 'Y':
			state.Pages.RemovePage("unsaved_confirm")
			state.TviewApp.SetFocus(state.Pages)
			onConfirm()
			return nil
		case 'n', 'N':
			state.Pages.RemovePage("unsaved_confirm")
			state.TviewApp.SetFocus(state.Pages)
			onCancel()
			return nil
		}
		if event.Key() == tcell.KeyEscape {
			state.Pages.RemovePage("unsaved_confirm")
			state.TviewApp.SetFocus(state.Pages)
			onCancel()
			return nil
		}
		return event
	})

	modal := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false).
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(box, 5, 0, true).
				AddItem(nil, 0, 1, false),
			50, 0, true,
		).
		AddItem(nil, 0, 1, false)

	state.Pages.AddPage("unsaved_confirm", modal, true, true)
	state.TviewApp.SetFocus(box)
}
