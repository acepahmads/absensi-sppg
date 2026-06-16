package service

import (
	"absensi-sppg/internal/model"
	"absensi-sppg/internal/repository"
	"math"
)

type InventoryService interface {
	GetAll(page, perPage int, dateFrom, dateTo, nameSearch, category, jenis string) (map[string]interface{}, error)
	GetBarangMasuk(page, perPage int, dateFrom, dateTo, nameSearch, category, jenis string) (map[string]interface{}, error)
	GetBarangKeluar(page, perPage int, dateFrom, dateTo, nameSearch, category, jenis string) (map[string]interface{}, error)
	InputBarang(Inventory model.Inventory) error
	DeleteBarang(id int) error
	UpdateBarang(Inventory model.Inventory) error
	InputBarangMasuk(Inventory model.InventoryBarangMasuk) error
	UpdateBarangMasuk(Inventory model.InventoryBarangMasuk) error
	DeleteBarangMasuk(id int) error
	InputBarangKeluar(Inventory model.InventoryBarangKeluar) error
	UpdateBarangKeluar(Inventory model.InventoryBarangKeluar) error
	DeleteBarangKeluar(id int) error
}

type inventoryService struct {
	inventoryRepo repository.InventoryRepository
}

func NewInventoryService(inventoryRepo repository.InventoryRepository) InventoryService {
	return &inventoryService{inventoryRepo}
}

func (s *inventoryService) GetAll(page, perPage int, dateFrom, dateTo, nameSearch, category, jenis string) (map[string]interface{}, error) {
	limit := 20

	data, total, err := s.inventoryRepo.GetAll(page, perPage, 0, dateFrom, dateTo, nameSearch, 0, category, jenis)
	if err != nil {
		return nil, err
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	response := map[string]interface{}{
		"data":         data,
		"current_page": 1,
		"total":        total,
		"total_pages":  totalPage,
	}

	return response, nil
}

func (s *inventoryService) GetBarangMasuk(page, perPage int, dateFrom, dateTo, nameSearch, category, jenis string) (map[string]interface{}, error) {
	limit := 20

	data, total, err := s.inventoryRepo.GetBarangMasuk(page, perPage, 0, dateFrom, dateTo, nameSearch, 0, category, jenis)
	if err != nil {
		return nil, err
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	response := map[string]interface{}{
		"data":         data,
		"current_page": 1,
		"total":        total,
		"total_pages":  totalPage,
	}

	return response, nil
}

func (s *inventoryService) GetBarangKeluar(page, perPage int, dateFrom, dateTo, nameSearch, category, jenis string) (map[string]interface{}, error) {
	limit := 20

	data, total, err := s.inventoryRepo.GetBarangKeluar(page, perPage, 0, dateFrom, dateTo, nameSearch, 0, category, jenis)
	if err != nil {
		return nil, err
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	response := map[string]interface{}{
		"data":         data,
		"current_page": 1,
		"total":        total,
		"total_pages":  totalPage,
	}

	return response, nil
}

func (s *inventoryService) InputBarang(Inventory model.Inventory) error {
	err := s.inventoryRepo.InputBarang(Inventory)
	if err != nil {
		return err
	}
	return nil
}

func (s *inventoryService) UpdateBarang(Inventory model.Inventory) error {
	err := s.inventoryRepo.UpdateBarang(Inventory)
	if err != nil {
		return err
	}
	return nil
}

func (s *inventoryService) DeleteBarang(id int) error {
	err := s.inventoryRepo.DeleteBarang(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *inventoryService) InputBarangMasuk(Inventory model.InventoryBarangMasuk) error {
	err := s.inventoryRepo.InputBarangMasuk(Inventory)
	if err != nil {
		return err
	}
	return nil
}

func (s *inventoryService) UpdateBarangMasuk(Inventory model.InventoryBarangMasuk) error {
	err := s.inventoryRepo.UpdateBarangMasuk(Inventory)
	if err != nil {
		return err
	}
	return nil
}

func (s *inventoryService) DeleteBarangMasuk(id int) error {
	err := s.inventoryRepo.DeleteBarangMasuk(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *inventoryService) InputBarangKeluar(Inventory model.InventoryBarangKeluar) error {
	err := s.inventoryRepo.InputBarangKeluar(Inventory)
	if err != nil {
		return err
	}
	return nil
}

func (s *inventoryService) UpdateBarangKeluar(Inventory model.InventoryBarangKeluar) error {
	err := s.inventoryRepo.UpdateBarangKeluar(Inventory)
	if err != nil {
		return err
	}
	return nil
}

func (s *inventoryService) DeleteBarangKeluar(id int) error {
	err := s.inventoryRepo.DeleteBarangKeluar(id)
	if err != nil {
		return err
	}
	return nil
}
