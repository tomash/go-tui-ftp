package ui

import (
	"simple-ftp-client/state"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewLayout(state *state.AppState, app *tview.Application) tview.Primitive {
	// Setup Status Bar: Green text, strictly no word wrapping to force horizontal flow
	statusBar := tview.NewTextView().
		SetWordWrap(false). // Critical for Windows consoles
		SetRegions(false).
		SetTextColor(tcell.ColorGreen).
		SetTextAlign(tview.AlignLeft)

	mainContent := tview.NewBox().SetBorder(true).SetTitle("File Browser")

	// Use Grid with explicit row sizes.
	// Row 0 = Fixed height (1 line for status bar).
	// Row 1 = Flexible (-1 takes all remaining space).
	layout := tview.NewGrid().SetRows(1, -1)

	// AddItem signature: (Primitive, rowStart, colStart, rowSpan, colSpan, xAlign, yAlign, focus)
	// Note: 'false' for the last argument prevents focus from stealing keyboard events.
	layout.AddItem(statusBar, 0, 0, 1, 1, tview.AlignLeft, 0, false)
	layout.AddItem(mainContent, 1, 0, 1, 1, tview.AlignCenter, 0, false)

	// Background goroutine to safely update UI from background FTP calls
	go func() {
		for ev := range state.UpdateCh {
			switch ev.Type {
			case "status":
				msg, ok := ev.Data.(string)
				if !ok {
					msg = "[unknown]"
				}
				app.QueueUpdateDraw(func() {
					statusBar.SetText(msg) // Plain text; color handled by SetTextColor()
				})
			default:
				// Future phases will handle "files", "progress", etc.
			}
		}
	}()

	state.UpdateStatus("Ready. Connect to a server when implemented.")

	return layout
}
