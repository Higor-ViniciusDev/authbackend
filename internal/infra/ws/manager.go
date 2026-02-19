package ws

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/Higor-ViniciusDev/auth/configuration/logger"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // In production: validate origin
	},
}

// Manager keeps a map of active WebSocket connections on this pod.
// key: user email → list of active connections in this pod
// thread-safe using RWMutex.
type Manager struct {
	mu          sync.RWMutex
	connections map[string][]*websocket.Conn
}

func NewManager() *Manager {
	return &Manager{
		connections: make(map[string][]*websocket.Conn),
	}
}

// HandleConnection is the HTTP handler for GET /ws/verify-status?email=x
// upgrades to WebSocket, registers the connection and waits.
// when the browser closes the tab, removes the connection from the map.
func (m *Manager) HandleConnection(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "email required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("WS Manager: upgrade error", err)
		return
	}

	m.register(email, conn)
	logger.Info(fmt.Sprintf("WS Manager: connection registered for %s", email))

	// blocks until the connection closes (browser closed tab or navigated away)
	defer func() {
		m.unregister(email, conn)
		conn.Close()
		logger.Info(fmt.Sprintf("WS Manager: connection removed for %s", email))
	}()

	// read loop needed to detect client disconnect
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

// Notify envia { "status": "verified" } para todas as conexões ativas do email neste pod.
// Chamado pelo EmailVerifiedConsumer quando recebe o evento do RabbitMQ.
func (m *Manager) Notify(email string) {
	m.mu.RLock()
	conns := m.connections[email]
	m.mu.RUnlock()

	if len(conns) == 0 {
		// this pod has no connection for this email — another pod will notify
		return
	}

	message := map[string]string{"status": "verified"}

	for _, conn := range conns {
		if err := conn.WriteJSON(message); err != nil {
			logger.Error(fmt.Sprintf("WS Manager: error notifying %s", email), err)
		}
	}
}

func (m *Manager) register(email string, conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connections[email] = append(m.connections[email], conn)
}

func (m *Manager) unregister(email string, conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	conns := m.connections[email]
	for i, c := range conns {
		if c == conn {
			m.connections[email] = append(conns[:i], conns[i+1:]...)
			break
		}
	}

	if len(m.connections[email]) == 0 {
		delete(m.connections, email)
	}
}
