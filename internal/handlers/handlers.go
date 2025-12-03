package handlers

import "github.com/gofiber/fiber/v2"

type AppHandler struct{}

func NewAppHandler() *AppHandler {
    return &AppHandler{}
}

func (a *AppHandler) HandleGetIndex(c *fiber.Ctx) error {
    return c.Render("index", fiber.Map{})
}
