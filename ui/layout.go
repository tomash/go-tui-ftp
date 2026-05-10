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

func NewLayout(appState *state.AppState, app *tview.Application, ftpClient FTPDriver) tview.Primitive {
	statusBar := tview.NewTextView().SetWordWrap(false).SetTextColor(tcell.ColorGreen).SetTextAlign(tview.AlignLeft)
	fileList := tview.NewList()
	fileList.SetBorder(true)
	fileList.SetTitle("File Browser")

	go func() {
		for ev := range appState.UpdateCh {
			switch ev.Type {
			case "status":
				msg, _ := ev.Data.(string)
				app.QueueUpdateDraw(func() { statusBar.SetText(msg) })
			case "files":
				data, ok := ev.Data.(state.FileListData)
				if !ok {
					app.QueueUpdateDraw(func() {
						fileList.Clear()
						fileList.AddItem("(invalid file list update)", "", 0, nil)
					})
					continue
				}
				app.QueueUpdateDraw(func() {
					fileList.Clear()
					switch {
					case data.Err != nil:
						fileList.AddItem("(listing failed)", "", 0, nil)
					case len(data.Names) == 0:
						fileList.AddItem("(empty directory)", "", 0, nil)
					default:
						for _, f := range data.Names {
							fileList.AddItem(f, "", 0, nil)
						}
					}
				})
			}
		}
	}()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == 'c' && !appState.IsConnected() {
			appState.UpdateStatus("Connecting...")
			go func() {
				err := ftpClient.Connect("ftp.funet.fi", 21, "anonymous", "")
				if err != nil {
					appState.UpdateStatus(fmt.Sprintf("Error: %v", err))
					return
				}
				appState.SetConnected(true)
				files, err := ftpClient.ListDir("/")
				appState.UpdateFiles(state.FileListData{Names: files, Err: err})
				if err != nil {
					appState.UpdateStatus(fmt.Sprintf("Connected, but listing failed: %v", err))
					return
				}
				appState.UpdateStatus("Connected! Files listed.")
			}()
			return nil
		}
		if event.Key() == tcell.KeyRune && event.Rune() == 'd' && appState.IsConnected() {
			if err := ftpClient.Disconnect(); err != nil {
				appState.UpdateStatus(fmt.Sprintf("Disconnect error: %v", err))
			} else {
				appState.SetConnected(false)
				app.QueueUpdateDraw(func() { fileList.Clear() })
				appState.UpdateStatus("Disconnected.")
			}
			return nil
		}
		return event
	})

	appState.UpdateStatus("Ready. Press 'c' to connect, 'd' to disconnect.")

	layout := tview.NewGrid().SetRows(1, -1)
	layout.AddItem(statusBar, 0, 0, 1, 1, tview.AlignLeft, 0, false)
	layout.AddItem(fileList, 1, 0, 1, 1, tview.AlignCenter, 0, true)

	return layout
}
