package state

import "sync"

type StateEvent struct {
	Type string // "status", "files"
	Data interface{}
}

// FileListData is sent with Type "files". Err is nil on success (Names may be empty).
type FileListData struct {
	Names []string
	Err   error
}

type AppState struct {
	mu        sync.RWMutex
	connected bool
	Path      string
	StatusMsg string
	Files     []string
	UpdateCh  chan StateEvent
}

func (s *AppState) SetConnected(v bool) {
	s.mu.Lock()
	s.connected = v
	s.mu.Unlock()
}

func (s *AppState) IsConnected() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.connected
}

func (s *AppState) UpdateStatus(msg string) {
	s.mu.Lock()
	s.StatusMsg = msg
	s.mu.Unlock()

	s.UpdateCh <- StateEvent{Type: "status", Data: msg}
}

// UpdateFiles pushes a directory listing to the UI. Err distinguishes failure from an empty directory.
func (s *AppState) UpdateFiles(data FileListData) {
	s.mu.Lock()
	if data.Err != nil {
		s.Files = nil
	} else {
		s.Files = data.Names
	}
	s.mu.Unlock()

	s.UpdateCh <- StateEvent{Type: "files", Data: data}
}
