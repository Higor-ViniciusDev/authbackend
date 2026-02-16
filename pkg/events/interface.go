package events

import (
	"sync"
	"time"
)

// Exemplos: usuário fez login, pedido foi criado, arquivo foi enviado, etc.
type EventInterface interface {
	// GetNome retorna o nome do evento, por exemplo: "UsuarioLogado".
	GetName() string
	// GetDateTime retorna a data e hora em que o evento ocorreu.
	GetDateTime() time.Time
	// GetValues retorna os dados associados ao evento.
	// Pode ser qualquer informação relevante sobre o evento.
	GetPayload() any
	SetPayload(payload any)
}

// EventoHandlerInterface representa um manipulador de eventos.
// Um handler é responsável por executar alguma ação quando um evento ocorre.
// Exemplo: enviar um e-mail quando um pedido é criado.
type EventHandlerInterface interface {
	// Handle executa a ação desejada ao receber um evento.
	Handle(event EventInterface, wg *sync.WaitGroup)
}

// EventDispachtInterface gerencia o registro e disparo de eventos e seus handlers.
// Ele permite adicionar, remover e acionar handlers para eventos específicos.
type EventDispachtInterface interface {
	// RegistrarHandler associa um handler a um evento específico pelo nome.
	// Assim, quando esse evento ocorrer, o handler será chamado.
	RegistrarHandler(eventoNome string, handler EventHandlerInterface) error
	// Dispatch dispara um evento, chamando todos os handlers registrados para ele.
	Dispatch(evento EventInterface) error
	// Remove remove um handler específico de um evento.
	Remove(eventoNome string, handler EventHandlerInterface) error
	// HasHandlers verifica se existem handlers registrados para um evento.
	HasHandlers(eventoNome string, handle EventHandlerInterface) bool
	// Clear remove todos os handlers de todos os eventos.
	Clear()
}
