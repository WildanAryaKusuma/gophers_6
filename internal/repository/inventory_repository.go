package repository

import (
	"belajar-golang/internal/model"
	"errors"

	"gorm.io/gorm"
)

var ErrStokTidakCukup = errors.New("Stok tidak cukup!")

type InventoryRepository interface {
	FindAll() ([]model.Inventory, error)
	FindByID(id uint) (model.Inventory, error)
	Create(barang *model.Inventory) error
	Beli(id uint, jumlah int, trx *model.Transaksi) error
	KurangiSatu(id uint) (bool, error)
}

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) FindAll() ([]model.Inventory, error) {
	var daftarBarang []model.Inventory
	err := r.db.Order("id asc").Find(&daftarBarang).Error
	return daftarBarang, err
}

func (r *inventoryRepository) FindByID(id uint) (model.Inventory, error) {
	var barang model.Inventory
	err := r.db.First(&barang, id).Error
	return barang, err
}

func (r *inventoryRepository) Create(barang *model.Inventory) error {
	return r.db.Create(barang).Error
}

// Beli mengurangi stok dan mencatat transaksi dalam satu database transaction.
// Pengecekan stok ada di klausa WHERE, jadi aman walau ratusan request datang
// bersamaan: database yang menjamin stok tidak pernah minus.
func (r *inventoryRepository) Beli(id uint, jumlah int, trx *model.Transaksi) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&model.Inventory{}).
			Where("id = ? AND stok >= ?", id, jumlah).
			Update("stok", gorm.Expr("stok - ?", jumlah))

		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrStokTidakCukup
		}

		return tx.Create(trx).Error
	})
}

// KurangiSatu dipakai flash sale: true kalau berhasil dapat barang.
func (r *inventoryRepository) KurangiSatu(id uint) (bool, error) {
	result := r.db.Model(&model.Inventory{}).
		Where("id = ? AND stok > 0", id).
		Update("stok", gorm.Expr("stok - 1"))

	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}
