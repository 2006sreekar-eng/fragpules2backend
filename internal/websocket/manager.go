package websocket

import (
	"encoding/json"
	"net/http"
	"sync"
	"github.com/gorilla/websocket"
)

type Client struct {
	manager *Manager
	conn    *websocket.Conn
	send    chan []byte
}

type Manager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true },
}

func NewManager() *Manager {
	return &Manager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (m *Manager) Run() {
	for {
		select {
		case client := <-m.register:
			m.mutex.Lock()
			m.clients[client] = true
			m.mutex.Unlock()
		case client := <-m.unregister:
			m.mutex.Lock()
			if _, ok := m.clients[client]; ok {
				delete(m.clients, client)
				close(client.send)
			}
			m.mutex.Unlock()
		case message := <-m.broadcast:
			m.mutex.RLock()
			for client := range m.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(m.clients, client)
				}
			}
			m.mutex.RUnlock()
		}
	}
}

func (m *Manager) BroadcastEvent(eventType string, data interface{}) {
	payload := map[string]interface{}{"event": eventType, "data": data}
	msg, err := json.Marshal(payload)
	if err != nil {
		return
	}
	m.broadcast <- msg
}

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	client := &Client{manager: m, conn: conn, send: make(chan []byte, 256)}
	m.register <- client

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.manager.unregister <- c
		c.conn.Close()
	}()
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (c *Client) writePump() {
	defer func() { c.conn.Close() }()
	for {
		msg, ok := <-c.send
		if !ok {
			_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		_ = c.conn.WriteMessage(websocket.TextMessage, msg)
	}
}