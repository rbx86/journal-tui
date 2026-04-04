package ui

import (
	"fmt"
	"journaltui/app"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func ShowTitlePrompt(state *app.AppState, onConfirm func(title string), onCancel func()) {
	counter := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignRight).
		SetText("[#1DB954]0/30")

	input := tview.NewInputField().
		SetLabel("Entry title: ").
		SetFieldWidth(30).
		SetLabelColor(ColorAccent).
		SetFieldTextColor(ColorText)

	input.SetAcceptanceFunc(func(text string, lastChar rune) bool {
		return len([]rune(text)) <= 30
	})

	input.SetChangedFunc(func(text string) {
		count := len([]rune(text))
		if count >= 30 {
			counter.SetText("[red]30/30")
		} else {
			counter.SetText(fmt.Sprintf("[#1DB954]%d/30", count))
		}
	})

	inputRow := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(input, 0, 1, true).
		AddItem(counter, 6, 0, false)

	box := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(inputRow, 1, 0, true)

	box.SetBorder(true).
		SetTitle(" Title ").
		SetTitleAlign(tview.AlignCenter)

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
				AddItem(box, 3, 0, true).
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
