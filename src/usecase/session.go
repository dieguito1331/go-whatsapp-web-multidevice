package usecase

import (
	"fmt"
	"sync"

	"github.com/aldinokemal/go-whatsapp-web-multidevice/config"
	domainSession "github.com/aldinokemal/go-whatsapp-web-multidevice/domains/session"
	"github.com/aldinokemal/go-whatsapp-web-multidevice/infrastructure/whatsapp"
)

type sessionManager struct {
	sessions map[string]*whatsapp.WaCli
	mu       sync.Mutex
}

func NewSessionManager() domainSession.ISessionManager {
	return &sessionManager{
		sessions: make(map[string]*whatsapp.WaCli),
	}
}

func (sm *sessionManager) CreateSession(sessionID string) (*whatsapp.WaCli, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, exists := sm.sessions[sessionID]; exists {
		return nil, fmt.Errorf("session with ID '%s' already exists", sessionID)
	}

	// For multi-session, each session needs its own database.
	dbPath := fmt.Sprintf("%s/session-%s.db", config.PathData, sessionID)
	cli, err := whatsapp.NewWaCLI(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create new whatsapp client: %w", err)
	}

	sm.sessions[sessionID] = cli
	return cli, nil
}

func (sm *sessionManager) GetSession(sessionID string) (*whatsapp.WaCli, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, exists := sm.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session with ID '%s' not found", sessionID)
	}
	return session, nil
}

func (sm *sessionManager) DeleteSession(sessionID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, exists := sm.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session with ID '%s' not found", sessionID)
	}

	session.Disconnect()
	delete(sm.sessions, sessionID)
	return nil
}

func (sm *sessionManager) GetAllSessions() map[string]*whatsapp.WaCli {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Return a copy to prevent race conditions on the original map
	sessionsCopy := make(map[string]*whatsapp.WaCli)
	for id, cli := range sm.sessions {
		sessionsCopy[id] = cli
	}
	return sessionsCopy
} 