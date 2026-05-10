package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"simple-ftp-client/ftp"
	"simple-ftp-client/state"
	"simple-ftp-client/ui"

	"github.com/rivo/tview"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	app := tview.NewApplication()

	// Initialize State and FTP Client separately to keep them decoupled
	stateData := &state.AppState{
		UpdateCh: make(chan state.StateEvent, 64),
	}
	ftpClient := ftp.NewClient()

	// Pass both into the UI builder
	layout := ui.NewLayout(stateData, app, ftpClient)

	go func() {
		<-done
		app.Stop()
	}()

	if err := app.SetRoot(layout, true).Run(); err != nil {
		log.Fatalf("UI error: %v", err)
	}

	println("\nShutting down cleanly.")
}
