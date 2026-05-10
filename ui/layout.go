package ui

import (
	"fmt"
	"simple-ftp-client/state"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FTPDriver interface {
	Connect(host string, port int, user, pass string) error
	ListDir(path string) ([]string, error)
	Disconnect() error
}

func NewLayout(state *state.AppState, app *tview.Application, ftpClient FTPDriver) tview.Primitive {
	statusBar := tview.NewTextView().SetWordWrap(false).SetTextColor(tcell.ColorGreen).SetTextAlign(tview.AlignLeft)
	fileList := tview.NewList()
	fileList.SetBorder(true)
	fileList.SetTitle("File Browser")

	go func() {
		for ev := range state.UpdateCh {
			switch ev.Type {
			case "status":
				msg, _ := ev.Data.(string)
				app.QueueUpdateDraw(func() { statusBar.SetText(msg) })
			case "files":
				files, ok := ev.Data.([]string)
				if !ok || len(files) == 0 {
					files = []string{"No files found"}
				}
				app.QueueUpdateDraw(func() {
					fileList.Clear()
					for _, f := range files {
						fileList.AddItem(f, "", 0, nil)
					}
				})
			}
		}
	}()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == 'c' && !state.Connected {
			state.UpdateStatus("Connecting...")
			go func() {
				err := ftpClient.Connect("ftp.funet.fi", 21, "anonymous", "")
				if err != nil {
					state.UpdateStatus(fmt.Sprintf("Error: %v", err))
					return
				}
				state.Connected = true
				files, err := ftpClient.ListDir("/")
				if err != nil {
					state.UpdateStatus(fmt.Sprintf("List Error: %v", err))
					return
				}
				state.UpdateFiles(files)
				state.UpdateStatus("Connected! Files listed.")
			}()
			return nil
		}
		if event.Key() == tcell.KeyRune && event.Rune() == 'd' && state.Connected {
			ftpClient.Disconnect()
			state.Connected = false
			fileList.Clear()
			state.UpdateStatus("Disconnected.")
			return nil
		}
		return event
	})

	state.UpdateStatus("Ready. Press 'c' to connect, 'd' to disconnect.")

	layout := tview.NewGrid().SetRows(1, -1)
	layout.AddItem(statusBar, 0, 0, 1, 1, tview.AlignLeft, 0, false)
	layout.AddItem(fileList, 1, 0, 1, 1, tview.AlignCenter, 0, true)

	return layout
}
