package model

import "gorm.io/gorm"

// Inventory adalah tabel barang di gudang toko kelontong.
type Inventory struct {
	gorm.Model
	NamaBarang string `gorm:"column:nama_barang" json:"nama_barang"`
	Harga      int    `gorm:"column:harga" json:"harga"`
	Stok       int    `gorm:"column:stok" json:"stok"`
}

// Transaksi mencatat setiap pembelian yang berhasil.
type Transaksi struct {
	gorm.Model
	InventoryID uint `json:"id_barang"`
	Jumlah      int  `json:"jumlah"`
	Total       int  `json:"total"`
	Uang        int  `json:"uang"`
	Kembalian   int  `json:"kembalian"`
}

type CreateInventoryRequest struct {
	NamaBarang string `json:"nama_barang"`
	Harga      int    `json:"harga"`
	Stok       int    `json:"stok"`
}

type BeliRequest struct {
	IDBarang uint `json:"id_barang"`
	Jumlah   int  `json:"jumlah"`
	Uang     int  `json:"uang"`
}

type FlashSaleRequest struct {
	Penyerbu int `json:"penyerbu"`
}

type FlashSaleResult struct {
	NamaBarang string  `json:"nama_barang"`
	StokAwal   int     `json:"stok_awal"`
	Penyerbu   int     `json:"penyerbu"`
	Terjual    int     `json:"terjual"`
	Gagal      int     `json:"gagal"`
	SisaStok   int     `json:"sisa_stok"`
	DurasiMs   float64 `json:"durasi_ms"`
	Aman       bool    `json:"aman"`
}
