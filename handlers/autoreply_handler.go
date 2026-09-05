package handlers

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"

	"gowavuelite/database"
	"gowavuelite/models"
)

type AutoReplyHandler struct {
}

func NewAutoReplyHandler() *AutoReplyHandler {
	return &AutoReplyHandler{}
}

func (h *AutoReplyHandler) List(c *fiber.Ctx) error {
	var rules []models.AutoReply
	if err := database.DB.Order("id asc").Find(&rules).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "db_error"})
	}
	return c.JSON(rules)
}

func (h *AutoReplyHandler) Create(c *fiber.Ctx) error {
	type req struct {
		Keyword   string `json:"keyword"`
		ReplyText string `json:"reply_text"`
		MatchType string `json:"match_type"`
		IsActive  bool   `json:"is_active"`
	}
	var body req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid_body"})
	}
	body.Keyword = strings.TrimSpace(body.Keyword)
	body.ReplyText = strings.TrimSpace(body.ReplyText)
	if body.Keyword == "" || body.ReplyText == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "missing_fields"})
	}
	if body.MatchType == "" {
		body.MatchType = "contains"
	}
	rule := models.AutoReply{
		Keyword:   body.Keyword,
		ReplyText: body.ReplyText,
		MatchType: body.MatchType,
		IsActive:  body.IsActive,
	}
	if err := database.DB.Create(&rule).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "db_error"})
	}
	return c.JSON(rule)
}

func (h *AutoReplyHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var rule models.AutoReply
	if err := database.DB.First(&rule, id).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "not_found"})
	}
	type req struct {
		Keyword   *string `json:"keyword"`
		ReplyText *string `json:"reply_text"`
		MatchType *string `json:"match_type"`
		IsActive  *bool   `json:"is_active"`
	}
	var body req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid_body"})
	}
	if body.Keyword != nil {
		rule.Keyword = strings.TrimSpace(*body.Keyword)
	}
	if body.ReplyText != nil {
		rule.ReplyText = strings.TrimSpace(*body.ReplyText)
	}
	if body.MatchType != nil {
		rule.MatchType = *body.MatchType
	}
	if body.IsActive != nil {
		rule.IsActive = *body.IsActive
	}
	if err := database.DB.Save(&rule).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "db_error"})
	}
	return c.JSON(rule)
}

func (h *AutoReplyHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := database.DB.Delete(&models.AutoReply{}, id).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "db_error"})
	}
	return c.JSON(fiber.Map{"success": true})
}

