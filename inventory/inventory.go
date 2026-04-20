package inventory

import (
	"fmt"
	"github.com/fatih/color"
)

type Inventory struct {
	namaBarang string
	harga      int
	stok       int
}

func TokoKelontong() {
	var daftarBarang []Inventory

	for {
		var opsi int
		fmt.Print("\n")
		fmt.Println("=== Toko Kelontong Inventory ===")
		fmt.Print("\n")
		fmt.Println("1. Tambah Barang ke Gudang")
		fmt.Println("2. Lihat Semua Stok Barang")
		fmt.Println("3. Beli Barang")
		fmt.Println("4. Keluar")
		fmt.Print("\n")
		fmt.Print("Pilih Menu (1-4): ")
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
			color.Green("Sampai jumpa! gawe gawe gawe")
			return

		} else {
			fmt.Println("Opsi tidak valid!")
		}
	}
}