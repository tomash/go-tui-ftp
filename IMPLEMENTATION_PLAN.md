# Simple FTP TUI Client - Implementation Plan

## Libraries
- UI: github.com/rivo/tview + github.com/gdamore/tcell/v2
- FTP: github.com/jlumbroso/golang-ftp
- State/Concurrency: sync, context, channels (stdlib)

## Architecture
- State-centric TUI: UI reacts to state changes via channels/mutex. Never blocks tview loop.
- Non-blocking FTP: All network calls run in background goroutines; progress/status routed back through channels.
- Layout: Top status bar | Middle split (file list + details/command) | Bottom command/progress

## Phases
1. Scaffolding & Dependencies (Current)
2. FTP Core Integration
3. TUI Layout & State Binding
4. Interactivity & Polish
5. Testing & Hardening
