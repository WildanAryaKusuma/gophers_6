package cashier

import (
	"fmt"
	"github.com/fatih/color"
)

func Payment () {
	var uang int
	harga := 20000 

	fmt.Println("Harga Barang: 20000")
	fmt.Print("Uang Pembeli: ")
  fmt.Scan(&uang)

	if uang >= 20000 {
		color.Green("[SISTEM] Transaksi Berhasil!. Kembalian Anda: %d", uang - harga)
	} else {
		fmt.Println("Uang Pembeli: ", uang)
		color.Red("[SISTEM] Transaksi Gagal PAK ZAKIII!. Uang Anda Kurang: %d", harga - uang)
	}
}