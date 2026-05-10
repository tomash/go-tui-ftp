package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"simple-ftp-client/state"
	"simple-ftp-client/ui"

	"github.com/rivo/tview"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	app := tview.NewApplication()

	stateData := &state.AppState{
		UpdateCh: make(chan state.StateEvent, 10),
	}

	layout := ui.NewLayout(stateData, app)

	go func() {
		<-done
		app.Stop()
	}()

	if err := app.SetRoot(layout, true).Run(); err != nil {
		log.Fatalf("UI error: %v", err)
	}

	fmt.Println("\nShutting down cleanly. Goodbye!")
}
