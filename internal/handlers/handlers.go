package handlers

import (
	"chat-server/internal/auth"
	"chat-server/internal/database"
	"chat-server/internal/models"
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
)

type AppHandler struct {
	db *database.DB
}

func NewAppHandler(db *database.DB) *AppHandler {
	return &AppHandler{db: db}
}

func (h *AppHandler) HandleGetIndex(c *fiber.Ctx) error {
	return c.Render("index", nil)
}

func (h *AppHandler) HandleRegister(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid request"})
	}

	if req.Username == "" || req.Password == "" {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Username and password required"})
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return c.Status(500).JSON(models.ErrorResponse{Error: "Internal server error"})
	}

	userID, err := h.db.CreateUser(req.Username, hashedPassword)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return c.Status(400).JSON(models.ErrorResponse{Error: "Username already exists"})
	}

	token, err := auth.GenerateToken(userID, req.Username)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return c.Status(500).JSON(models.ErrorResponse{Error: "Internal server error"})
	}

	return c.JSON(models.AuthResponse{
		Token:    token,
		Username: req.Username,
		UserID:   userID,
	})
}

func (h *AppHandler) HandleLogin(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Invalid request"})
	}

	if req.Username == "" || req.Password == "" {
		return c.Status(400).JSON(models.ErrorResponse{Error: "Username and password required"})
	}

	userID, hashedPassword, err := h.db.GetUserByUsername(req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(401).JSON(models.ErrorResponse{Error: "Invalid credentials"})
		}
		log.Printf("Error getting user: %v", err)
		return c.Status(500).JSON(models.ErrorResponse{Error: "Internal server error"})
	}

	if !auth.CheckPasswordHash(req.Password, hashedPassword) {
		return c.Status(401).JSON(models.ErrorResponse{Error: "Invalid credentials"})
	}

	token, err := auth.GenerateToken(userID, req.Username)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return c.Status(500).JSON(models.ErrorResponse{Error: "Internal server error"})
	}

	return c.JSON(models.AuthResponse{
		Token:    token,
		Username: req.Username,
		UserID:   userID,
	})
}

func (h *AppHandler) HandleGetMessages(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(401).JSON(models.ErrorResponse{Error: "Authorization required"})
	}

	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	_, err := auth.ValidateToken(tokenString)
	if err != nil {
		return c.Status(401).JSON(models.ErrorResponse{Error: "Invalid token"})
	}

	room := c.Query("room", "")
	limit := c.QueryInt("limit", 50)

	var messages []database.Message
	if room != "" {
		messages, err = h.db.GetMessagesByRoom(room, limit)
	} else {
		messages, err = h.db.GetRecentMessages(limit)
	}

	if err != nil {
		log.Printf("Error getting messages: %v", err)
		return c.Status(500).JSON(models.ErrorResponse{Error: "Internal server error"})
	}

	return c.JSON(messages)
}
