package main

import (
	"log"

	"journaltui/app"
	"journaltui/storage"
	"journaltui/ui"
)

func main() {
	if err := storage.Init(); err != nil {
		log.Fatalf("failed to initialize data directory: %v", err)
	}

	entries, err := storage.LoadMetadata()
	if err != nil {
		log.Fatalf("failed to load metadata: %v", err)
	}

	state := app.New()
	state.Entries = entries
	ui.SetupMainMenu(state)
	ui.SetupEditor(state)
	ui.SetupMainMenu(state)
	ui.SetupEditor(state)
	ui.SetupTable(state)

	if err := state.TviewApp.SetRoot(state.Pages, true).Run(); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
