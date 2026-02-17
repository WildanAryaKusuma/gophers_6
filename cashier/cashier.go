package cashier

import (
	"fmt"
	"github.com/fatih/color"
)

func Payment () {
	var uang int

	fmt.Println("Harga Barang: 20000")
	fmt.Print("Uang Pembeli: ")
  fmt.Scan(&uang)

	if uang >= 20000 {
		color.Green("[SISTEM] Transaksi Berhasil!. Kembalian Anda: %d", uang - 20000)
	} else {
		fmt.Println("Uang Pembeli: ", uang)
		color.Red("[SISTEM] Transaksi Gagal PAK ZAKIII!. Uang Anda Kurang: %d", 20000 - uang)
	}
}