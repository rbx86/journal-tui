package app

import (
	"github.com/rivo/tview"
)

type AppState struct {
	TviewApp                *tview.Application
	Pages                   *tview.Pages
	Entries                 map[string]EntryMeta
	CurrentEntryID          string
	IsDirty                 bool
	IsNewEntry              bool
	IsReadMode              bool
	EditorTextArea          *tview.TextArea
	EditorBox               *tview.Flex
	UpdateEditorHint        func()
	UpdateEditorBorderTitle func(*tview.Flex, string)
	TableWidget             *tview.Table
	PopulateTable           func()
}

func New() *AppState {
	return &AppState{
		TviewApp:   tview.NewApplication(),
		Pages:      tview.NewPages(),
		Entries:    make(map[string]EntryMeta),
		IsNewEntry: true,
		IsReadMode: false,
		IsDirty:    false,
	}
}
