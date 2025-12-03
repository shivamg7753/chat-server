package main

import (
    "chat-server/internal/handlers"
    "chat-server/internal/websocket"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/template/html/v2"
    ws "github.com/gofiber/websocket/v2"
)

func main() {

    engine := html.New("./views", ".html")

    app := fiber.New(fiber.Config{
        Views: engine,
    })

    
    app.Static("/static", "./static")

    
    h := handlers.NewAppHandler()
    app.Get("/", h.HandleGetIndex)

    
    server := websocket.NewWebSocketServer()

    app.Get("/ws", ws.New(func(c *ws.Conn) {
        server.HandleWebSocket(c)
    }))

    go server.HandleMessages()

    app.Listen(":3000")
}
