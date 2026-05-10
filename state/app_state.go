package state

import (
	"sync"
)

// StateEvent represents UI updates from background goroutines
type StateEvent struct {
	Type string // "status", "files", "progress", "error"
	Data interface{}
}

// AppState holds centralized application data
type AppState struct {
	mu        sync.RWMutex
	Connected bool
	Path      string
	StatusMsg string
	UpdateCh  chan StateEvent
}

// UpdateStatus safely pushes a status change to the UI channel
func (s *AppState) UpdateStatus(msg string) {
	s.mu.Lock()
	s.StatusMsg = msg
	s.mu.Unlock()

	select {
	case s.UpdateCh <- StateEvent{Type: "status", Data: msg}:
	default: // Drop if channel full to avoid blocking goroutines
	}
}

// GetState returns a safe snapshot of current state
func (s *AppState) GetState() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return map[string]interface{}{
		"connected": s.Connected,
		"path":      s.Path,
		"status":    s.StatusMsg,
	}
}
