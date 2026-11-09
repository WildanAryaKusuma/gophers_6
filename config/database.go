package config

import (
	"belajar-golang/internal/model"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

var dataAwal = []model.Inventory{
	{NamaBarang: "Indomie Goreng", Harga: 3000, Stok: 50},
	{NamaBarang: "Beras 5kg", Harga: 75000, Stok: 10},
	{NamaBarang: "Minyak Goreng 2L", Harga: 14000, Stok: 100},
}

func ConnectPostgres() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		GetEnv("DB_HOST", "localhost"),
		GetEnv("DB_PORT", "5432"),
		GetEnv("DB_USER", "postgres"),
		GetEnv("DB_PASSWORD", "Mcrazor28"),
		GetEnv("DB_NAME", "tokoanu"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gagal terhubung ke database: %w", err)
	}

	if err := db.AutoMigrate(&model.Inventory{}, &model.Transaksi{}); err != nil {
		return nil, fmt.Errorf("gagal migrasi: %w", err)
	}

	if err := seedJikaKosong(db); err != nil {
		return nil, fmt.Errorf("gagal seeding data awal: %w", err)
	}

	log.Println("Berhasil terhubung ke database!")
	return db, nil
}

func seedJikaKosong(db *gorm.DB) error {
	var jumlah int64
	if err := db.Model(&model.Inventory{}).Count(&jumlah).Error; err != nil {
		return err
	}
	if jumlah > 0 {
		return nil
	}

	if err := db.Create(&dataAwal).Error; err != nil {
		return err
	}
	log.Println("Data awal berhasil di-seed ke database.")
	return nil
}
