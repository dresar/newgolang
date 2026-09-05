package handlers

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"

	"gowavuelite/database"
	"gowavuelite/models"
	"gowavuelite/services"
)

type MessageHandler struct {
	wa *services.WhatsAppService
}

func NewMessageHandler(wa *services.WhatsAppService) *MessageHandler {
	return &MessageHandler{wa: wa}
}

func (h *MessageHandler) SendMessage(c *fiber.Ctx) error {
	type req struct {
		To      string `json:"to"`
		Message string `json:"message"`
	}
	var body req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid_body"})
	}
	body.To = strings.TrimSpace(body.To)
	body.Message = strings.TrimSpace(body.Message)
	if body.To == "" || body.Message == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "missing_fields"})
	}
	log := models.MessageLog{
		To:      body.To,
		Message: body.Message,
		Status:  "pending",
	}
	if err := database.DB.Create(&log).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "db_error"})
	}
	err := h.wa.SendTextMessage(body.To, body.Message)
	if err != nil {
		log.Status = "failed"
		log.Error = err.Error()
		database.DB.Save(&log)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "send_failed"})
	}
	log.Status = "sent"
	database.DB.Save(&log)
	return c.JSON(fiber.Map{"success": true, "log_id": log.ID})
}

func (h *MessageHandler) ListLogs(c *fiber.Ctx) error {
	var logs []models.MessageLog
	if err := database.DB.Order("created_at desc").Limit(100).Find(&logs).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "db_error"})
	}
	return c.JSON(logs)
}

