package ui

import (
	"fmt"
	"journaltui/app"
	"journaltui/storage"
	"time"

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

func ShowSearchPrompt(state *app.AppState, onSearch func(titleQuery, dateQuery string), onCancel func()) {
	titleInput := tview.NewInputField().
		SetLabel("Title:  ").
		SetFieldWidth(30).
		SetLabelColor(ColorAccent).
		SetFieldTextColor(ColorText)

	dateInput := tview.NewInputField().
		SetLabel("Date:   ").
		SetFieldWidth(30).
		SetLabelColor(ColorAccent).
		SetFieldTextColor(ColorText)

	titleInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			state.TviewApp.SetFocus(dateInput)
			return nil
		}
		if event.Key() == tcell.KeyEnter {
			state.Pages.RemovePage("search_prompt")
			onSearch(titleInput.GetText(), dateInput.GetText())
			return nil
		}
		if event.Key() == tcell.KeyEscape {
			state.Pages.RemovePage("search_prompt")
			onCancel()
			return nil
		}
		return event
	})

	dateInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			state.TviewApp.SetFocus(titleInput)
			return nil
		}
		if event.Key() == tcell.KeyEnter {
			state.Pages.RemovePage("search_prompt")
			onSearch(titleInput.GetText(), dateInput.GetText())
			return nil
		}
		if event.Key() == tcell.KeyEscape {
			state.Pages.RemovePage("search_prompt")
			onCancel()
			return nil
		}
		return event
	})

	hint := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText("[#1DB954]tab[white] switch fields  [#1DB954]enter[white] search  [#1DB954]esc[white] cancel")

	box := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(titleInput, 1, 0, true).
		AddItem(tview.NewBox(), 1, 0, false).
		AddItem(dateInput, 1, 0, false).
		AddItem(tview.NewBox(), 1, 0, false).
		AddItem(hint, 1, 0, false)

	box.SetBorder(true).
		SetTitle(" Search ").
		SetTitleAlign(tview.AlignCenter)

	modal := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false).
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(box, 7, 0, true).
				AddItem(nil, 0, 1, false),
			60, 0, true,
		).
		AddItem(nil, 0, 1, false)

	state.Pages.AddPage("search_prompt", modal, true, true)
	state.TviewApp.SetFocus(titleInput)
}

func ShowExportProgress(state *app.AppState, statusText *tview.TextView) {
	rainbowColors := []string{
		"[#FF0000]", // red
		"[#FF7700]", // orange
		"[#FFFF00]", // yellow
		"[#00FF00]", // green
		"[#0000FF]", // blue
		"[#4B0082]", // indigo
		"[#8B00FF]", // violet
	}

	progressBar := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)

	percentText := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignRight)

	progressRow := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(progressBar, 0, 1, false).
		AddItem(percentText, 6, 0, false)

	box := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 1, 0, false).
		AddItem(progressRow, 1, 0, false).
		AddItem(tview.NewBox(), 1, 0, false)

	box.SetBorder(true).
		SetTitle(" Exporting ").
		SetTitleAlign(tview.AlignCenter)

	modal := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false).
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(box, 5, 0, false).
				AddItem(nil, 0, 1, false),
			50, 0, false,
		).
		AddItem(nil, 0, 1, false)

	state.Pages.AddPage("export_progress", modal, true, true)

	updateProgress := func(percent int) {
		state.TviewApp.QueueUpdateDraw(func() {
			barWidth := 38
			filled := int(float64(percent) / 100.0 * float64(barWidth))
			bar := ""
			for i := 0; i < filled; i++ {
				colorIdx := i % len(rainbowColors)
				bar += rainbowColors[colorIdx] + "█"
			}
			bar += "[#333333]"
			for i := filled; i < barWidth; i++ {
				bar += "░"
			}
			progressBar.SetText(bar)
			percentText.SetText(fmt.Sprintf("[white]%3d%%", percent))
		})
	}

	go func() {
		err := storage.ExportEntries(updateProgress)

		state.TviewApp.QueueUpdateDraw(func() {
			state.Pages.RemovePage("export_progress")
			if err != nil {
				statusText.SetText("[red]Export failed.")
			} else {
				statusText.SetText("[#1DB954]Exported entries to ~/Downloads")
			}
		})

		go func() {
			time.Sleep(8 * time.Second)
			state.TviewApp.QueueUpdateDraw(func() {
				statusText.SetText("")
			})
		}()
	}()
}
