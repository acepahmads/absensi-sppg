package handler

import (
	"absensi-sppg/internal/model"
	"absensi-sppg/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	inventoryService service.InventoryService
}

func NewInventoryHandler(inventoryService service.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
	}
}

func (h *InventoryHandler) GetAll(c *gin.Context) {
	pageStr := c.Query("page")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	perPageStr := c.Query("per_page")
	perPage, _ := strconv.Atoi(perPageStr)
	if perPage < 1 {
		perPage = 0
	}

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	nameSearch := c.Query("search")
	category := c.Query("kategori")
	jenis := c.Query("jenis")

	result, err := h.inventoryService.GetAll(page, perPage, dateFrom, dateTo, nameSearch, category, jenis)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed load Inventory data",
		})
		return
	}

	// PRINT RAW (debug)
	b, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(b))

	rawData, ok := result["data"]
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid data format",
		})
		return
	}

	records, ok := rawData.([]model.Inventory)
	if !ok {
		fmt.Printf("Invalid type: %T\n", rawData)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid data structure (expected []model.Inventory)",
		})
		return
	}

	var responseData []model.InventoryResponse
	for _, r := range records {
		if r.Gambar != "" {
			if !strings.HasPrefix(r.Gambar, "static/uploads/inventory/") {
				r.Gambar = "static/uploads/inventory/" + r.Gambar
			}
		}
		responseData = append(responseData, model.InventoryResponse{
			ID:           r.ID,
			QRCode:       r.QRCode,
			NamaBarang:   r.NamaBarang,
			Kategori:     r.Kategori,
			JenisBarang:  r.JenisBarang,
			Satuan:       r.Satuan,
			StokAwal:     r.StokAwal,
			StokMasuk:    r.StokMasuk,
			Posisi:       r.Posisi,
			Gambar:       r.Gambar,
			StokAkhir:    r.StokAkhir,
			BarangMasuk:  r.BarangMasuk,
			BarangKeluar: r.BarangKeluar,
			HargaBeli:    r.HargaBeli,
			HargaJual:    r.HargaJual,
			Keterangan:   r.Keterangan,
			CreatedAt:    r.CreatedAt,
		})
	}

	// fmt.Println(responseData)

	// PAGINATION FIELDS
	total := 0
	if v, ok := result["total"].(int); ok {
		total = v
	}
	fmt.Println("total", total)

	totalPages := 0
	if v, ok := result["total_pages"].(int); ok {
		totalPages = v
	}
	fmt.Println("total_pages", totalPages)
	fmt.Println("current_page", page)
	// fmt.Println("data", responseData)

	// Final Response
	c.JSON(http.StatusOK, gin.H{
		"current_page": page,
		"total":        total,
		"total_pages":  totalPages,
		"data":         responseData,
	})
}

func (h *InventoryHandler) GetBarangMasuk(c *gin.Context) {
	pageStr := c.Query("page")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	perPageStr := c.Query("per_page")
	perPage, _ := strconv.Atoi(perPageStr)
	if perPage < 1 {
		perPage = 0
	}

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	nameSearch := c.Query("search")
	category := c.Query("kategori")
	jenis := c.Query("jenis")
	result, err := h.inventoryService.GetBarangMasuk(page, perPage, dateFrom, dateTo, nameSearch, category, jenis)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed load Inventory data",
		})
		return
	}

	b, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(b))

	rawData, ok := result["data"]
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid data format",
		})
		return
	}

	records, ok := rawData.([]model.InventoryBarangMasuk)
	if !ok {
		fmt.Printf("Invalid type: %T\n", rawData)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid data structure (expected []model.InventoryBarangMasuk)",
		})
		return
	}

	var responseData []model.InventoryBarangMasuk

	for _, r := range records {
		responseData = append(responseData, model.InventoryBarangMasuk{
			ID:          r.ID,
			QRCode:      r.QRCode,
			TanggalJam:  r.TanggalJam,
			NamaBarang:  r.NamaBarang,
			Kategori:    r.Kategori,
			JenisBarang: r.JenisBarang,
			Jumlah:      r.Jumlah,
			Keterangan:  r.Keterangan,
		})
	}

	// fmt.Println(responseData)

	// PAGINATION FIELDS
	total := 0
	if v, ok := result["total"].(int); ok {
		total = v
	}
	fmt.Println("total", total)

	totalPages := 0
	if v, ok := result["total_pages"].(int); ok {
		totalPages = v
	}
	fmt.Println("total_pages", totalPages)
	fmt.Println("current_page", page)
	// fmt.Println("data", responseData)

	// Final Response
	c.JSON(http.StatusOK, gin.H{
		"current_page": page,
		"total":        total,
		"total_pages":  totalPages,
		"data":         responseData,
	})
}

func (h *InventoryHandler) GetBarangKeluar(c *gin.Context) {
	pageStr := c.Query("page")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	perPageStr := c.Query("per_page")
	perPage, _ := strconv.Atoi(perPageStr)
	if perPage < 1 {
		perPage = 0
	}

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	nameSearch := c.Query("search")
	category := c.Query("kategori")
	jenis := c.Query("jenis")
	result, err := h.inventoryService.GetBarangKeluar(page, perPage, dateFrom, dateTo, nameSearch, category, jenis)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed load Inventory data",
		})
		return
	}

	b, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(b))

	rawData, ok := result["data"]
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid data format",
		})
		return
	}

	records, ok := rawData.([]model.InventoryBarangKeluar)
	if !ok {
		fmt.Printf("Invalid type: %T\n", rawData)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid data structure (expected []model.InventoryBarangKeluar)",
		})
		return
	}

	var responseData []model.InventoryBarangKeluar

	for _, r := range records {
		responseData = append(responseData, model.InventoryBarangKeluar{
			ID:          r.ID,
			QRCode:      r.QRCode,
			TanggalJam:  r.TanggalJam,
			NamaBarang:  r.NamaBarang,
			Kategori:    r.Kategori,
			JenisBarang: r.JenisBarang,
			Jumlah:      r.Jumlah,
			Keterangan:  r.Keterangan,
		})
	}

	// fmt.Println(responseData)

	// PAGINATION FIELDS
	total := 0
	if v, ok := result["total"].(int); ok {
		total = v
	}
	fmt.Println("total", total)

	totalPages := 0
	if v, ok := result["total_pages"].(int); ok {
		totalPages = v
	}
	fmt.Println("total_pages", totalPages)
	fmt.Println("current_page", page)
	// fmt.Println("data", responseData)

	// Final Response
	c.JSON(http.StatusOK, gin.H{
		"current_page": page,
		"total":        total,
		"total_pages":  totalPages,
		"data":         responseData,
	})
}

func (h *InventoryHandler) InputBarang(c *gin.Context) {
	var Inventory model.Inventory
	if err := c.ShouldBindJSON(&Inventory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := h.inventoryService.InputBarang(Inventory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to input inventory data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory data input successful"})
}

func (h *InventoryHandler) UpdateBarang(c *gin.Context) {
	idStr := c.Param("id")
	// id, err := strconv.Atoi(idStr)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inventory ID"})
	// 	return
	// }

	var Inventory model.Inventory
	if err := c.ShouldBindJSON(&Inventory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	Inventory.ID = idStr // Set the ID from the URL parameter

	err := h.inventoryService.UpdateBarang(Inventory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update inventory data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory data updated successfully"})
}

func (h *InventoryHandler) DeleteBarang(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inventory ID"})
		return
	}

	err = h.inventoryService.DeleteBarang(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete inventory data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory data deleted successfully"})
}

func (h *InventoryHandler) InputBarangMasuk(c *gin.Context) {
	var Inventory model.InventoryBarangMasuk
	if err := c.ShouldBindJSON(&Inventory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := h.inventoryService.InputBarangMasuk(Inventory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to input inventory data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory data input successful"})
}

func (h *InventoryHandler) UpdateBarangMasuk(c *gin.Context) {
	idStr := c.Param("id")
	// id, err := strconv.Atoi(idStr)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inventory ID"})
	// 	return
	// }

	var Inventory model.InventoryBarangMasuk
	if err := c.ShouldBindJSON(&Inventory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	Inventory.ID = idStr // Set the ID from the URL parameter

	err := h.inventoryService.UpdateBarangMasuk(Inventory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update inventory data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory data updated successfully"})
}

func (h *InventoryHandler) DeleteBarangMasuk(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inventory ID"})
		return
	}

	err = h.inventoryService.DeleteBarangMasuk(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete inventory data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory data deleted successfully"})
}

func (h *InventoryHandler) InputBarangKeluar(c *gin.Context) {
	var Inventory model.InventoryBarangKeluar
	if err := c.ShouldBindJSON(&Inventory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := h.inventoryService.InputBarangKeluar(Inventory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to input inventory data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory data input successful"})
}

func (h *InventoryHandler) UpdateBarangKeluar(c *gin.Context) {
	idStr := c.Param("id")
	// id, err := strconv.Atoi(idStr)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inventory ID"})
	// 	return
	// }

	var Inventory model.InventoryBarangKeluar
	if err := c.ShouldBindJSON(&Inventory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	Inventory.ID = idStr // Set the ID from the URL parameter

	err := h.inventoryService.UpdateBarangKeluar(Inventory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update inventory data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory data updated successfully"})
}

func (h *InventoryHandler) DeleteBarangKeluar(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inventory ID"})
		return
	}

	err = h.inventoryService.DeleteBarangKeluar(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete inventory data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory data deleted successfully"})
}
