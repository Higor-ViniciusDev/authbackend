package events

import (
	"sync"
	"time"
)

type EventInterface interface {
	GetName() string
	GetDateTime() time.Time
	GetPayload() any
	SetPayload(payload any)
}

type EventHandlerInterface interface {
	Handle(event EventInterface, wg *sync.WaitGroup)
}

type EventDispachtInterface interface {
	RegistrarHandler(eventoNome string, handler EventHandlerInterface) error
	Dispatch(evento EventInterface) error
	Remove(eventoNome string, handler EventHandlerInterface) error
	HasHandlers(eventoNome string, handle EventHandlerInterface) bool
	Clear()
}
