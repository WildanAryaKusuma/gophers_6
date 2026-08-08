package inventory

import (
	"fmt"
	"sync"
	"time"

	"github.com/fatih/color"
)

type Inventory struct {
	namaBarang string
	harga      int
	stok       int
}

var dataAwal = []Inventory{
	{namaBarang: "Indomie Goreng", harga: 3000, stok: 50},
	{namaBarang: "Beras 5kg", harga: 75000, stok: 10},
	{namaBarang: "Minyak Goreng 2L", harga: 14000, stok: 100},
}

const jumlahPenyerbuFlashSale = 1000

func TokoKelontong() {
	daftarBarang := make([]Inventory, len(dataAwal))
	copy(daftarBarang, dataAwal)

	for {
		var opsi int
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
		fmt.Scan(&opsi)

		if opsi == 1 {
			var barang Inventory
			fmt.Print("Nama barang  : ")
			fmt.Scan(&barang.namaBarang)
			fmt.Print("Harga barang : ")
			fmt.Scan(&barang.harga)
			fmt.Print("Stok barang  : ")
			fmt.Scan(&barang.stok)
			daftarBarang = append(daftarBarang, barang)
			color.Green("Barang berhasil ditambahkan!")

		} else if opsi == 2 {
			if len(daftarBarang) == 0 {
				color.Red("Inventory masih kosong!")
				continue
			}
			fmt.Println("\n--- Daftar Barang ---")
			fmt.Printf("%-5s %-20s %-15s %-10s\n", "No", "Nama Barang", "Harga", "Stok")
			fmt.Println("----------------------------------------------------")
			for i, barang := range daftarBarang {
				fmt.Printf("%-5d %-20s Rp%-13d %-10d\n",
					i+1,
					barang.namaBarang,
					barang.harga,
					barang.stok,
				)
			}

		} else if opsi == 3 {
			if len(daftarBarang) == 0 {
				color.Red("Inventory masih kosong!")
				continue
			}

			fmt.Println("\n=== MENU PEMBELIAN ===")
			for i, barang := range daftarBarang {
				fmt.Printf("ID: %d | Nama: %-15s | Stok: %d | Harga: Rp %d\n",
					i+1,
					barang.namaBarang,
					barang.stok,
					barang.harga,
				)
			}

			var idPilihan, jumlah int
			fmt.Print("\nPilih ID Barang yang mau dibeli: ")
			fmt.Scan(&idPilihan)

			if idPilihan < 1 || idPilihan > len(daftarBarang) {
				color.Red("[SISTEM]: ID barang tidak valid!")
				continue
			}

			fmt.Print("Jumlah yang mau dibeli: ")
			fmt.Scan(&jumlah)

			if jumlah <= 0 {
				color.Red("[SISTEM]: Jumlah tidak valid!")
				continue
			}

			barang := &daftarBarang[idPilihan-1]

			if jumlah > barang.stok {
				color.Red("[SISTEM]: Stok tidak cukup! Stok tersedia: %d", barang.stok)
				continue
			}

			totalHarga := barang.harga * jumlah
			fmt.Printf("\n[SISTEM]: Total Harga: Rp %d\n", totalHarga)

			var uangMasuk int
			fmt.Print("Masukkan Uang Anda: ")
			fmt.Scan(&uangMasuk)

			if uangMasuk < totalHarga {
				color.Red("[SISTEM]: Uang tidak cukup! Kurang Rp %d", totalHarga-uangMasuk)
				continue
			}

			barang.stok -= jumlah
			kembalian := uangMasuk - totalHarga

			fmt.Print("\n")
			color.Green("[SISTEM]: Transaksi Berhasil!")
			color.Green("[SISTEM]: Kembalian Anda: Rp %d", kembalian)
			color.Green("[SISTEM]: Stok %s sekarang: %d pcs", barang.namaBarang, barang.stok)
			fmt.Println("\n--------------------------------------------")

		} else if opsi == 4 {
			if len(daftarBarang) == 0 {
				color.Red("Inventory masih kosong!")
				continue
			}

			fmt.Println("\n= DAFTAR BARANG =")
			for i, barang := range daftarBarang {
				fmt.Printf("ID: %d | Nama: %-18s | Stok: %-5d | Harga: Rp %d\n",
					i+1, barang.namaBarang, barang.stok, barang.harga)
			}

			var pilihan int
			fmt.Print("\nPilih ID barang untuk FLASH SALE 12.12: ")
			fmt.Scan(&pilihan)

			if pilihan < 1 || pilihan > len(daftarBarang) {
				color.Red("[SISTEM]: ID barang tidak valid!")
				continue
			}

			barang := &daftarBarang[pilihan-1]
			jumlahPenyerbu := jumlahPenyerbuFlashSale

			fmt.Printf("\nSERBUAN 12.12 untuk \"%s\"\n", barang.namaBarang)
			fmt.Printf("Stok %d, penyerbu %d orang.\n", barang.stok, jumlahPenyerbu)

			var wg sync.WaitGroup
			var mu sync.Mutex // <- Mesin Transaksi Aman: gerbang berkunci
			terjual := 0
			gagal := 0

			mulai := time.Now()

			for i := 0; i < jumlahPenyerbu; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()

					mu.Lock()
					defer mu.Unlock()

					if barang.stok > 0 {
						time.Sleep(time.Microsecond)
						barang.stok--
						terjual++
					} else {
						gagal++
					}
				}()
			}

			wg.Wait()
			durasi := time.Since(mulai)

			fmt.Println("\n= FLASH SALE SELESAI =")
			fmt.Printf("Waktu proses : %.2fms\n", float64(durasi.Microseconds())/1000)
			fmt.Printf("Berhasil beli: %d unit\n", terjual)
			fmt.Printf("Gagal (kehabisan): %d pembeli\n", gagal)
			fmt.Printf("Sisa stok    : %d unit\n", barang.stok)

			if barang.stok < 0 {
				color.Red("BAHAYA! Stok %s MINUS %d -- Toko menjual barang yang tidak ada!",
					barang.namaBarang, -barang.stok)
			} else {
				color.Green("AMAN! Terjual pas %d unit, stok %s tidak pernah minus.",
					terjual, barang.namaBarang)
			}

		} else if opsi == 5 {
			color.Green("Sampai jumpa! gawe gawe gawe")
			return

		} else {
			fmt.Println("Opsi tidak valid!")
		}
	}
}
