package main

import (
	"chat-server/internal/config"
	"chat-server/internal/database"
	"chat-server/internal/handlers"
	"chat-server/internal/websocket"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
	ws "github.com/gofiber/websocket/v2"
)

func main() {
	config.Load()

	log.Println("Starting Chat Server...")

	db, err := database.NewDB(config.DBPath)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	engine := html.New("./internal/views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	app.Static("/static", "./internal/static")

	h := handlers.NewAppHandler(db)
	app.Get("/", h.HandleGetIndex)

	app.Post("/api/register", h.HandleRegister)
	app.Post("/api/login", h.HandleLogin)
	app.Get("/api/messages", h.HandleGetMessages)

	server := websocket.NewWebSocketServer(db)

	app.Get("/ws", ws.New(func(c *ws.Conn) {
		server.HandleWebSocket(c)
	}))

	go server.HandleMessages()

	log.Println("Server running on http://localhost:" + config.Port)
	log.Println("Endpoints: POST /api/register, POST /api/login, GET /api/messages, WS /ws?token=<jwt>")

	if err := app.Listen(":" + config.Port); err != nil {
		log.Fatal("Server failed:", err)
	}
}
