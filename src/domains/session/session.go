package session

import (
	"github.com/aldinokemal/go-whatsapp-web-multidevice/infrastructure/whatsapp"
)

// ISessionManager defines the interface for managing WhatsApp sessions.
type ISessionManager interface {
	CreateSession(sessionID string) (*whatsapp.WaCli, error)
	GetSession(sessionID string) (*whatsapp.WaCli, error)
	DeleteSession(sessionID string) error
	GetAllSessions() map[string]*whatsapp.WaCli
} 