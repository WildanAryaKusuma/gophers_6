package handler

import (
	"belajar-golang/internal/model"
	"belajar-golang/internal/service"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type InventoryHandler struct {
	svc service.InventoryService
}

func NewInventoryHandler(svc service.InventoryService) *InventoryHandler {
	return &InventoryHandler{svc: svc}
}

// errorResponse menerjemahkan error dari service ke HTTP status.
func errorResponse(c fiber.Ctx, err error, pesanInvalid string) error {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		return c.Status(400).JSON(fiber.Map{"error": pesanInvalid})
	case errors.Is(err, service.ErrStokTidakCukup):
		return c.Status(400).JSON(fiber.Map{"error": "Stok tidak cukup"})
	case errors.Is(err, service.ErrUangKurang):
		return c.Status(400).JSON(fiber.Map{"error": "Uang tidak cukup"})
	case errors.Is(err, service.ErrNotFound):
		return c.Status(404).JSON(fiber.Map{"error": "Barang tidak ditemukan"})
	default:
		return c.Status(500).JSON(fiber.Map{"error": "Kesalahan server!"})
	}
}

func parseID(c fiber.Ctx) (uint, bool) {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil || id == 0 {
		return 0, false
	}
	return uint(id), true
}

func (h *InventoryHandler) GetAll(c fiber.Ctx) error {
	daftarBarang, err := h.svc.GetAll()
	if err != nil {
		return errorResponse(c, err, "")
	}
	return c.JSON(daftarBarang)
}

func (h *InventoryHandler) GetByID(c fiber.Ctx) error {
	id, ok := parseID(c)
	if !ok {
		return c.Status(400).JSON(fiber.Map{"error": "ID barang tidak valid"})
	}

	barang, err := h.svc.GetByID(id)
	if err != nil {
		return errorResponse(c, err, "")
	}
	return c.JSON(barang)
}

func (h *InventoryHandler) Create(c fiber.Ctx) error {
	var req model.CreateInventoryRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Data tidak valid"})
	}

	barang, err := h.svc.Create(req)
	if err != nil {
		return errorResponse(c, err, "Nama wajib diisi, harga harus > 0, dan stok tidak boleh minus")
	}
	return c.Status(201).JSON(barang)
}

func (h *InventoryHandler) Beli(c fiber.Ctx) error {
	var req model.BeliRequest
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Data tidak valid"})
	}

	barang, trx, err := h.svc.Beli(req)
	if err != nil {
		return errorResponse(c, err, "Jumlah beli harus > 0 dan uang tidak boleh minus")
	}

	return c.JSON(fiber.Map{
		"message":     "Transaksi berhasil",
		"nama_barang": barang.NamaBarang,
		"jumlah":      trx.Jumlah,
		"total":       trx.Total,
		"uang":        trx.Uang,
		"kembalian":   trx.Kembalian,
		"sisa_stok":   barang.Stok,
	})
}

func (h *InventoryHandler) FlashSale(c fiber.Ctx) error {
	id, ok := parseID(c)
	if !ok {
		return c.Status(400).JSON(fiber.Map{"error": "ID barang tidak valid"})
	}

	// Body opsional: kosong berarti pakai jumlah penyerbu default.
	var req model.FlashSaleRequest
	if body := c.Body(); len(body) > 0 {
		if err := json.Unmarshal(body, &req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Data tidak valid"})
		}
	}

	hasil, err := h.svc.FlashSale(id, req.Penyerbu)
	if err != nil {
		return errorResponse(c, err, "Jumlah penyerbu harus 1 - 10000")
	}
	return c.JSON(hasil)
}
