package service

import (
	"belajar-golang/internal/model"
	"belajar-golang/internal/repository"
	"errors"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
)

var (
	ErrInvalidInput   = errors.New("Input tidak valid!")
	ErrNotFound       = errors.New("Barang tidak ditemukan!")
	ErrStokTidakCukup = repository.ErrStokTidakCukup
	ErrUangKurang     = errors.New("Uang tidak cukup!")
)

const (
	defaultPenyerbu = 1000
	maxPenyerbu     = 10000
)

type InventoryService interface {
	GetAll() ([]model.Inventory, error)
	GetByID(id uint) (model.Inventory, error)
	Create(req model.CreateInventoryRequest) (model.Inventory, error)
	Beli(req model.BeliRequest) (model.Inventory, model.Transaksi, error)
	FlashSale(id uint, penyerbu int) (model.FlashSaleResult, error)
}

type inventoryService struct {
	repo repository.InventoryRepository
}

func NewInventoryService(repo repository.InventoryRepository) InventoryService {
	return &inventoryService{repo: repo}
}

func (s *inventoryService) GetAll() ([]model.Inventory, error) {
	return s.repo.FindAll()
}

func (s *inventoryService) GetByID(id uint) (model.Inventory, error) {
	barang, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Inventory{}, ErrNotFound
		}
		return model.Inventory{}, err
	}
	return barang, nil
}

func (s *inventoryService) Create(req model.CreateInventoryRequest) (model.Inventory, error) {
	nama := strings.TrimSpace(req.NamaBarang)
	if nama == "" || req.Harga <= 0 || req.Stok < 0 {
		return model.Inventory{}, ErrInvalidInput
	}

	barang := model.Inventory{NamaBarang: nama, Harga: req.Harga, Stok: req.Stok}
	if err := s.repo.Create(&barang); err != nil {
		return model.Inventory{}, err
	}
	return barang, nil
}

func (s *inventoryService) Beli(req model.BeliRequest) (model.Inventory, model.Transaksi, error) {
	if req.Jumlah <= 0 || req.Uang < 0 {
		return model.Inventory{}, model.Transaksi{}, ErrInvalidInput
	}

	barang, err := s.GetByID(req.IDBarang)
	if err != nil {
		return model.Inventory{}, model.Transaksi{}, err
	}

	total := barang.Harga * req.Jumlah
	if req.Uang < total {
		return model.Inventory{}, model.Transaksi{}, ErrUangKurang
	}

	trx := model.Transaksi{
		InventoryID: barang.ID,
		Jumlah:      req.Jumlah,
		Total:       total,
		Uang:        req.Uang,
		Kembalian:   req.Uang - total,
	}

	// Stok dicek ulang secara atomik di repository, bukan dari data yang dibaca di atas.
	if err := s.repo.Beli(barang.ID, req.Jumlah, &trx); err != nil {
		return model.Inventory{}, model.Transaksi{}, err
	}

	barangTerbaru, err := s.GetByID(barang.ID)
	if err != nil {
		return model.Inventory{}, model.Transaksi{}, err
	}
	return barangTerbaru, trx, nil
}

// FlashSale mensimulasikan banyak pembeli (goroutine) yang berebut 1 unit barang
// secara bersamaan.
func (s *inventoryService) FlashSale(id uint, penyerbu int) (model.FlashSaleResult, error) {
	if penyerbu == 0 {
		penyerbu = defaultPenyerbu
	}
	if penyerbu < 0 || penyerbu > maxPenyerbu {
		return model.FlashSaleResult{}, ErrInvalidInput
	}

	barangAwal, err := s.GetByID(id)
	if err != nil {
		return model.FlashSaleResult{}, err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	terjual := 0
	gagal := 0

	mulai := time.Now()

	for i := 0; i < penyerbu; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			berhasil, err := s.repo.KurangiSatu(id)

			mu.Lock()
			defer mu.Unlock()
			if err == nil && berhasil {
				terjual++
			} else {
				gagal++
			}
		}()
	}

	wg.Wait()
	durasi := time.Since(mulai)

	barangAkhir, err := s.GetByID(id)
	if err != nil {
		return model.FlashSaleResult{}, err
	}

	return model.FlashSaleResult{
		NamaBarang: barangAkhir.NamaBarang,
		StokAwal:   barangAwal.Stok,
		Penyerbu:   penyerbu,
		Terjual:    terjual,
		Gagal:      gagal,
		SisaStok:   barangAkhir.Stok,
		DurasiMs:   float64(durasi.Microseconds()) / 1000,
		Aman:       barangAkhir.Stok >= 0,
	}, nil
}
