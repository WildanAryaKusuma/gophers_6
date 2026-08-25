package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"belajar-golang/inventory"
)

type Products struct {
	gorm.Model
	Nama  string
	Harga int
	Stok  int
}

func main() {
	dsn := "host=localhost port=5432 user=postgres password=Mcrazor28 dbname=tokoanu sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal membuka koneksi ke database: ", err)
	}
	fmt.Println("Berhasil konek ke database")

	// if err := db.AutoMigrate(&Products{}); err != nil {
	// 	log.Fatal("Gagal migrasi: ", err)
	// }
	// fmt.Println("Berhasil migrasi produk")

	// ===== CREATE =====
	// produkBaru := []Products{
	// 	{Nama: "Beras Petruk 5kg", Harga: 120000, Stok: 10},
	// 	{Nama: "Gula 500g", Harga: 10000, Stok: 50},
	// }
	// if err := db.Create(&produkBaru).Error; err != nil {
	// 	log.Fatal("Gagal insert produk: ", err)
	// }
	// fmt.Println("Berhasil insert produk baru")

	// ===== READ (semua data) =====
	// var semuaProduk []Products
	// db.Find(&semuaProduk)

	// fmt.Println("\nDaftar Semua Produk")
	// fmt.Println("ID | Nama | Harga | Stok")
	// for _, p := range semuaProduk {
	// 	fmt.Println(p.ID, p.Nama, p.Harga, p.Stok)
	// }

	// ===== READ (satu data spesifik) =====
	// var berasPetruk Products
	// if err := db.First(&berasPetruk, "nama = ?", "Beras Petruk 5kg").Error; err != nil {
	// 	log.Println("Produk tidak ditemukan: ", err)
	// 	return
	// }

	// ===== UPDATE =====
	// jumlahBeli := 5
	// if berasPetruk.Stok < jumlahBeli {
	// 	fmt.Println("Stok beras tidak cukup")
	// } else {
	// 	berasPetruk.Stok -= jumlahBeli
	// 	db.Save(&berasPetruk)
	// 	fmt.Println("Stok berhasil diupdate, sisa:", berasPetruk.Stok)
	// }

	// ===== DELETE =====
	// db.Delete(&berasPetruk)
	// fmt.Println("Produk berhasil dihapus (soft delete)")

	inventory.TokoKelontong(db)
}