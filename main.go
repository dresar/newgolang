package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

	"gowavuelite/database"
	"gowavuelite/handlers"
	"gowavuelite/models"
	"gowavuelite/routes"
	"gowavuelite/services"
)

func ensureAdminUser() {
	var count int64
	database.DB.Model(&models.User{}).Count(&count)
	if count == 0 {
		user := models.User{
			Username:     "admin",
			PasswordHash: handlers.HashPassword("admin123"),
		}
		database.DB.Create(&user)
	}
}

func main() {
	if err := database.ConnectDatabase(); err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	ensureAdminUser()

	waService, err := services.NewWhatsAppService("whatsmeow.db")
	if err != nil {
		log.Fatalf("failed to init whatsapp service: %v", err)
	}
	go func() {
		err := waService.Start()
		if err != nil {
			log.Printf("whatsapp start error: %v", err)
		}
	}()

	app := fiber.New()
	store := session.New()

	routes.Setup(app, store, waService)

	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
