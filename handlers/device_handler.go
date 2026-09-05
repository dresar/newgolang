package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"gowavuelite/services"
)

type DeviceHandler struct {
	wa *services.WhatsAppService
}

func NewDeviceHandler(wa *services.WhatsAppService) *DeviceHandler {
	return &DeviceHandler{wa: wa}
}

func (h *DeviceHandler) GetStatus(c *fiber.Ctx) error {
	status := h.wa.Status()
	return c.JSON(status)
}

func (h *DeviceHandler) GetQR(c *fiber.Ctx) error {
	status := h.wa.Status()
	if status.LastQRCode == "" {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "no_qr"})
	}
	return c.JSON(fiber.Map{"qr": status.LastQRCode})
}

func (h *DeviceHandler) Reconnect(c *fiber.Ctx) error {
	if err := h.wa.Reconnect(); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true})
}

