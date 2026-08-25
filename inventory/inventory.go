package inventory

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"gorm.io/gorm"
)

var reader = bufio.NewReader(os.Stdin)

func bacaBaris() string {
	teks, _ := reader.ReadString('\n')
	return strings.TrimSpace(teks)
}

func bacaAngka() (int, error) {
	teks := bacaBaris()
	return strconv.Atoi(teks)
}

type Inventory struct {
	gorm.Model
	NamaBarang string `gorm:"column:nama_barang"`
	Harga      int    `gorm:"column:harga"`
	Stok       int    `gorm:"column:stok"`
}

var dataAwal = []Inventory{
	{NamaBarang: "Indomie Goreng", Harga: 3000, Stok: 50},
	{NamaBarang: "Beras 5kg", Harga: 75000, Stok: 10},
	{NamaBarang: "Minyak Goreng 2L", Harga: 14000, Stok: 100},
}

const jumlahPenyerbuFlashSale = 1000

func seedJikaKosong(db *gorm.DB) error {
	var jumlah int64
	if err := db.Model(&Inventory{}).Count(&jumlah).Error; err != nil {
		return err
	}
	if jumlah == 0 {
		if err := db.Create(&dataAwal).Error; err != nil {
			return err
		}
		color.Yellow("[SISTEM]: Data awal berhasil di-seed ke database.")
	}
	return nil
}

func TokoKelontong(db *gorm.DB) {
	if err := db.AutoMigrate(&Inventory{}); err != nil {
		color.Red("[SISTEM]: Gagal migrasi tabel inventory: %v", err)
		return
	}

	if err := seedJikaKosong(db); err != nil {
		color.Red("[SISTEM]: Gagal seeding data awal: %v", err)
		return
	}

	for {
		fmt.Print("\n")
		fmt.Println("=== Toko Kelontong Inventory ===")
		fmt.Print("\n")
		fmt.Println("1. Tambah Barang ke Gudang")
		fmt.Println("2. Lihat Semua Stok Barang")
		fmt.Println("3. Beli Barang")
		fmt.Println("4. Flash Sale 12.12")
		fmt.Println("5. Keluar")
		fmt.Print("\n")
		fmt.Print("Pilih Menu (1-5): ")
		opsi, err := bacaAngka()
		if err != nil {
			fmt.Println("Opsi tidak valid!")
			continue
		}

		switch opsi {
		case 1:
			tambahBarang(db)
		case 2:
			lihatSemuaBarang(db)
		case 3:
			beliBarang(db)
		case 4:
			flashSale(db)
		case 5:
			color.Green("Sampai jumpa! gawe gawe gawe")
			return
		default:
			fmt.Println("Opsi tidak valid!")
		}
	}
}

func tambahBarang(db *gorm.DB) {
	var barang Inventory

	fmt.Print("Nama barang  : ")
	barang.NamaBarang = bacaBaris()
	if barang.NamaBarang == "" {
		color.Red("[SISTEM]: Nama barang tidak boleh kosong!")
		return
	}

	fmt.Print("Harga barang : ")
	harga, err := bacaAngka()
	if err != nil {
		color.Red("[SISTEM]: Harga harus berupa angka!")
		return
	}
	barang.Harga = harga

	fmt.Print("Stok barang  : ")
	stok, err := bacaAngka()
	if err != nil {
		color.Red("[SISTEM]: Stok harus berupa angka!")
		return
	}
	barang.Stok = stok

	if err := db.Create(&barang).Error; err != nil {
		color.Red("[SISTEM]: Gagal menambahkan barang: %v", err)
		return
	}
	color.Green("Barang berhasil ditambahkan!")
}

func lihatSemuaBarang(db *gorm.DB) {
	var daftarBarang []Inventory
	if err := db.Order("id asc").Find(&daftarBarang).Error; err != nil {
		color.Red("[SISTEM]: Gagal mengambil data barang: %v", err)
		return
	}

	if len(daftarBarang) == 0 {
		color.Red("Inventory masih kosong!")
		return
	}

	fmt.Println("\n--- Daftar Barang ---")
	fmt.Printf("%-5s %-20s %-15s %-10s\n", "ID", "Nama Barang", "Harga", "Stok")
	fmt.Println("----------------------------------------------------")
	for _, barang := range daftarBarang {
		fmt.Printf("%-5d %-20s Rp%-13d %-10d\n",
			barang.ID,
			barang.NamaBarang,
			barang.Harga,
			barang.Stok,
		)
	}
}

func beliBarang(db *gorm.DB) {
	var daftarBarang []Inventory
	if err := db.Order("id asc").Find(&daftarBarang).Error; err != nil {
		color.Red("[SISTEM]: Gagal mengambil data barang: %v", err)
		return
	}

	if len(daftarBarang) == 0 {
		color.Red("Inventory masih kosong!")
		return
	}

	fmt.Println("\n=== MENU PEMBELIAN ===")
	for _, barang := range daftarBarang {
		fmt.Printf("ID: %d | Nama: %-15s | Stok: %d | Harga: Rp %d\n",
			barang.ID,
			barang.NamaBarang,
			barang.Stok,
			barang.Harga,
		)
	}

	fmt.Print("\nPilih ID Barang yang mau dibeli: ")
	idAngka, err := bacaAngka()
	if err != nil || idAngka <= 0 {
		color.Red("[SISTEM]: ID barang tidak valid!")
		return
	}
	idPilihan := uint(idAngka)

	fmt.Print("Jumlah yang mau dibeli: ")
	jumlah, err := bacaAngka()
	if err != nil || jumlah <= 0 {
		color.Red("[SISTEM]: Jumlah tidak valid!")
		return
	}

	var barang Inventory
	if err := db.First(&barang, idPilihan).Error; err != nil {
		color.Red("[SISTEM]: ID barang tidak valid!")
		return
	}

	totalHarga := barang.Harga * jumlah

	fmt.Printf("\n[SISTEM]: Total Harga: Rp %d\n", totalHarga)
	fmt.Print("Masukkan Uang Anda: ")
	uangMasuk, err := bacaAngka()
	if err != nil || uangMasuk < 0 {
		color.Red("[SISTEM]: Input uang tidak valid!")
		return
	}

	if uangMasuk < totalHarga {
		color.Red("[SISTEM]: Uang tidak cukup! Kurang Rp %d", totalHarga-uangMasuk)
		return
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&Inventory{}).
			Where("id = ? AND stok >= ?", idPilihan, jumlah).
			Update("stok", gorm.Expr("stok - ?", jumlah))

		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("stok tidak cukup")
		}
		return nil
	})

	if err != nil {
		color.Red("[SISTEM]: Transaksi gagal: %v", err)
		return
	}

	db.First(&barang, idPilihan)
	kembalian := uangMasuk - totalHarga

	fmt.Print("\n")
	color.Green("[SISTEM]: Transaksi Berhasil!")
	color.Green("[SISTEM]: Kembalian Anda: Rp %d", kembalian)
	color.Green("[SISTEM]: Stok %s sekarang: %d pcs", barang.NamaBarang, barang.Stok)
	fmt.Println("\n--------------------------------------------")
}

func flashSale(db *gorm.DB) {
	var daftarBarang []Inventory
	if err := db.Order("id asc").Find(&daftarBarang).Error; err != nil {
		color.Red("[SISTEM]: Gagal mengambil data barang: %v", err)
		return
	}

	if len(daftarBarang) == 0 {
		color.Red("Inventory masih kosong!")
		return
	}

	fmt.Println("\n= DAFTAR BARANG =")
	for _, barang := range daftarBarang {
		fmt.Printf("ID: %d | Nama: %-18s | Stok: %-5d | Harga: Rp %d\n",
			barang.ID, barang.NamaBarang, barang.Stok, barang.Harga)
	}

	fmt.Print("\nPilih ID barang untuk FLASH SALE 12.12: ")
	idAngka, err := bacaAngka()
	if err != nil || idAngka <= 0 {
		color.Red("[SISTEM]: ID barang tidak valid!")
		return
	}
	idPilihan := uint(idAngka)

	var barangAwal Inventory
	if err := db.First(&barangAwal, idPilihan).Error; err != nil {
		color.Red("[SISTEM]: ID barang tidak valid!")
		return
	}

	jumlahPenyerbu := jumlahPenyerbuFlashSale

	fmt.Printf("\nSERBUAN 12.12 untuk \"%s\"\n", barangAwal.NamaBarang)
	fmt.Printf("Stok %d, penyerbu %d orang.\n", barangAwal.Stok, jumlahPenyerbu)

	var wg sync.WaitGroup
	var mu sync.Mutex
	terjual := 0
	gagal := 0

	mulai := time.Now()

	for i := 0; i < jumlahPenyerbu; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			var berhasil bool
			err := db.Transaction(func(tx *gorm.DB) error {
				result := tx.Model(&Inventory{}).
					Where("id = ? AND stok > 0", idPilihan).
					Update("stok", gorm.Expr("stok - 1"))

				if result.Error != nil {
					return result.Error
				}
				berhasil = result.RowsAffected > 0
				return nil
			})

			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				gagal++
				return
			}
			if berhasil {
				terjual++
			} else {
				gagal++
			}
		}()
	}

	wg.Wait()
	durasi := time.Since(mulai)

	var barangAkhir Inventory
	db.First(&barangAkhir, idPilihan)

	fmt.Println("\n= FLASH SALE SELESAI =")
	fmt.Printf("Waktu proses : %.2fms\n", float64(durasi.Microseconds())/1000)
	fmt.Printf("Berhasil beli: %d unit\n", terjual)
	fmt.Printf("Gagal (kehabisan): %d pembeli\n", gagal)
	fmt.Printf("Sisa stok    : %d unit\n", barangAkhir.Stok)

	if barangAkhir.Stok < 0 {
		color.Red("BAHAYA! Stok %s MINUS %d -- Toko menjual barang yang tidak ada!",
			barangAkhir.NamaBarang, -barangAkhir.Stok)
	} else {
		color.Green("AMAN! Terjual pas %d unit, stok %s tidak pernah minus.",
			terjual, barangAkhir.NamaBarang)
	}
}