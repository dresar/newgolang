package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"

	"gowavuelite/handlers"
	"gowavuelite/services"
)

func Setup(app *fiber.App, store *session.Store, wa *services.WhatsAppService) {
	app.Use(logger.New())

	app.Static("/", "./public")

	authHandler := handlers.NewAuthHandler(store)
	deviceHandler := handlers.NewDeviceHandler(wa)
	messageHandler := handlers.NewMessageHandler(wa)
	autoReplyHandler := handlers.NewAutoReplyHandler()

	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)
	auth.Get("/me", authHandler.Me)

	protected := api.Use(handlers.RequireAuth(store))

	device := protected.Group("/device")
	device.Get("/status", deviceHandler.GetStatus)
	device.Get("/qr", deviceHandler.GetQR)
	device.Post("/reconnect", deviceHandler.Reconnect)

	msg := protected.Group("/messages")
	msg.Post("/send", messageHandler.SendMessage)
	msg.Get("/logs", messageHandler.ListLogs)

	ar := protected.Group("/autoreply")
	ar.Get("/", autoReplyHandler.List)
	ar.Post("/", autoReplyHandler.Create)
	ar.Put("/:id", autoReplyHandler.Update)
	ar.Delete("/:id", autoReplyHandler.Delete)
}
