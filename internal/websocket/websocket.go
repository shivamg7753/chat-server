package websocket

import (
    "bytes"
    "chat-server/internal"
    "encoding/json"
    "html/template"
    "log"

    "github.com/gofiber/websocket/v2"
)

type WebSocketServer struct {
    Clients   map[*websocket.Conn]bool
    Broadcast chan *internal.Message
}

func NewWebSocketServer() *WebSocketServer {
    return &WebSocketServer{
        Clients:   make(map[*websocket.Conn]bool),
        Broadcast: make(chan *internal.Message),
    }
}

func (s *WebSocketServer) HandleWebSocket(c *websocket.Conn) {

    s.Clients[c] = true

    defer func() {
        delete(s.Clients, c)
        c.Close()
    }()

    for {
        _, msg, err := c.ReadMessage()
        if err != nil {
            log.Println("Socket read error:", err)
            break
        }

        var m internal.Message
        json.Unmarshal(msg, &m)

        s.Broadcast <- &m
    }
}

func (s *WebSocketServer) HandleMessages() {
    for {
        msg := <-s.Broadcast

        data := renderTemplate(msg)

        for client := range s.Clients {
            client.WriteMessage(websocket.TextMessage, data)
        }
    }
}

func renderTemplate(msg *internal.Message) []byte {
    t, err := template.ParseFiles("./internal/views/message.html")
    if err != nil {
        log.Println("template error:", err)
    }

    var buf bytes.Buffer
    t.Execute(&buf, msg)
    return buf.Bytes()
}
