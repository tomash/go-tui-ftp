package state

import "sync"

type StateEvent struct {
	Type string // "status", "files", "error"
	Data interface{}
}

type AppState struct {
	mu        sync.RWMutex
	Connected bool
	Path      string
	StatusMsg string
	Files     []string // Added for Phase 2
	UpdateCh  chan StateEvent
}

func (s *AppState) UpdateStatus(msg string) {
	s.mu.Lock()
	s.StatusMsg = msg
	s.mu.Unlock()

	select {
	case s.UpdateCh <- StateEvent{Type: "status", Data: msg}:
	default:
	}
}

// UpdateFiles pushes a new file list to the UI safely
func (s *AppState) UpdateFiles(files []string) {
	s.mu.Lock()
	s.Files = files // Copy slice reference
	s.mu.Unlock()

	select {
	case s.UpdateCh <- StateEvent{Type: "files", Data: files}:
	default:
	}
}
