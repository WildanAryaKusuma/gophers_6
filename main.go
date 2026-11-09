package main

import (
	"belajar-golang/config"
	"belajar-golang/internal/handler"
	"belajar-golang/internal/repository"
	"belajar-golang/internal/service"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
)

func main() {
	db, err := config.ConnectPostgres()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Gagal membuka database: ", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	inventoryRepo := repository.NewInventoryRepository(db)
	inventoryService := service.NewInventoryService(inventoryRepo)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)

	app := fiber.New()

	api := app.Group("/api")
	api.Get("/inventory", inventoryHandler.GetAll)
	api.Get("/inventory/:id", inventoryHandler.GetByID)
	api.Post("/inventory", inventoryHandler.Create)
	api.Post("/inventory/:id/flash-sale", inventoryHandler.FlashSale)
	api.Post("/beli", inventoryHandler.Beli)

	port := config.GetEnv("APP_PORT", "3000")
	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Fatal(err)
		}
	}()
	log.Println("Server sudah jalan di http://localhost:" + port)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Mematikan server dengan bersih...")
	if err := app.Shutdown(); err != nil {
		log.Fatal("Gagal shutdown: ", err)
	}
	if err := sqlDB.Close(); err != nil {
		log.Println("Gagal menutup koneksi database: ", err)
	}

	log.Println("Server berhasil dimatikan dengan bersih")
}
