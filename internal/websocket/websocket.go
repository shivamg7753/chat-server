package websocket

import (
	"bytes"
	"chat-server/internal/auth"
	"chat-server/internal/database"
	"chat-server/internal/models"
	"encoding/json"
	"html/template"
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
)

type WebSocketServer struct {
	Clients    map[*websocket.Conn]*ClientInfo
	Broadcast  chan *models.Message
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
	mu         sync.RWMutex
	db         *database.DB
}

type ClientInfo struct {
	UserID   int64
	Username string
}

func NewWebSocketServer(db *database.DB) *WebSocketServer {
	return &WebSocketServer{
		Clients:    make(map[*websocket.Conn]*ClientInfo),
		Broadcast:  make(chan *models.Message, 256),
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),
		db:         db,
	}
}

func (s *WebSocketServer) HandleWebSocket(c *websocket.Conn) {
	tokenString := c.Query("token")
	if tokenString == "" {
		log.Println("WebSocket connection rejected: no token provided")
		c.WriteMessage(websocket.TextMessage, []byte("Authentication required"))
		c.Close()
		return
	}

	claims, err := auth.ValidateToken(tokenString)
	if err != nil {
		log.Printf("WebSocket connection rejected: invalid token - %v", err)
		c.WriteMessage(websocket.TextMessage, []byte("Invalid token"))
		c.Close()
		return
	}

	s.mu.Lock()
	s.Clients[c] = &ClientInfo{
		UserID:   claims.UserID,
		Username: claims.Username,
	}
	s.mu.Unlock()

	log.Printf("User %s connected. Total clients: %d", claims.Username, len(s.Clients))

	defer func() {
		s.mu.Lock()
		delete(s.Clients, c)
		clientCount := len(s.Clients)
		s.mu.Unlock()
		c.Close()
		log.Printf("User %s disconnected. Total clients: %d", claims.Username, clientCount)
	}()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var m models.Message
		if err := json.Unmarshal(msg, &m); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		m.User = claims.Username

		if err := s.db.SaveMessage(claims.UserID, m.User, m.Text, m.Room); err != nil {
			log.Printf("Error saving message: %v", err)
		}

		log.Printf("Message from %s: %s", m.User, m.Text)

		s.Broadcast <- &m
	}
}

func (s *WebSocketServer) HandleMessages() {
	for {
		msg := <-s.Broadcast

		data := renderTemplate(msg)

		s.mu.RLock()
		for client := range s.Clients {
			err := client.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Printf("Error sending message to client: %v", err)
			}
		}
		s.mu.RUnlock()

		log.Printf("Broadcasted message to %d clients", len(s.Clients))
	}
}

func renderTemplate(msg *models.Message) []byte {
	t, err := template.ParseFiles("./internal/views/message.html")
	if err != nil {
		log.Printf("Template error: %v", err)
		return []byte{}
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, msg); err != nil {
		log.Printf("Template execution error: %v", err)
		return []byte{}
	}

	return buf.Bytes()
}
