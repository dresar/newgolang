package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

	"gowavuelite/database"
	"gowavuelite/models"
)

type AuthHandler struct {
	store *session.Store
}

func NewAuthHandler(store *session.Store) *AuthHandler {
	return &AuthHandler{store: store}
}

func HashPassword(password string) string {
	h := sha256.Sum256([]byte(password))
	return hex.EncodeToString(h[:])
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	type req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var body req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid_body"})
	}
	body.Username = strings.TrimSpace(body.Username)
	if body.Username == "" || body.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "missing_credentials"})
	}
	var user models.User
	if err := database.DB.Where("username = ?", body.Username).First(&user).Error; err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid_credentials"})
	}
	if user.PasswordHash != HashPassword(body.Password) {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid_credentials"})
	}
	sess, err := h.store.Get(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "session_error"})
	}
	sess.Set("user_id", user.ID)
	if err := sess.Save(); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "session_error"})
	}
	return c.JSON(fiber.Map{
		"user": fiber.Map{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	sess, err := h.store.Get(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "session_error"})
	}
	sess.Destroy()
	return c.JSON(fiber.Map{"success": true})
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	sess, err := h.store.Get(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	userID := sess.Get("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	return c.JSON(fiber.Map{
		"id":       user.ID,
		"username": user.Username,
	})
}

func RequireAuth(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		if sess.Get("user_id") == nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		return c.Next()
	}
}
