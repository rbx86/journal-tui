package main

import (
	"fmt"
	"log"
	"os"

	"journaltui/app"
	"journaltui/storage"
	"journaltui/ui"

	"golang.org/x/term"
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
	ui.ApplyTheme()
	ui.SetupMainMenu(state)
	ui.SetupEditor(state)
	ui.SetupMainMenu(state)
	ui.SetupEditor(state)
	ui.SetupTable(state)

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err == nil && (width < 100 || height < 25) {
		fmt.Fprintf(os.Stderr, "Terminal too small (%dx%d). Please resize to at least 100x25.\n", width, height)
		os.Exit(1)
	}

	if err := state.TviewApp.SetRoot(state.Pages, true).Run(); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
