package handler

import (
	"absensi-sppg/internal/model"
	"absensi-sppg/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"absensi-sppg/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AbsensiResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Date           string `json:"date"`
	DateOut        string `json:"date_out"`
	TimeIn         string `json:"time_in"`
	TimeOut        string `json:"time_out"`
	Status         string `json:"status"`
	LateMinutes    int    `json:"LateMinutes"`
	Deduction      string `json:"deduction"`
	Notes          string `json:"notes"`
	PICNotes       string `json:"pic_notes"`
	PhotoMasuk     string `json:"photo_masuk"`
	PhotoPulang    string `json:"photo_pulang"`
	PhotoBukti1    string `json:"photo_bukti1"`
	PhotoBukti2    string `json:"photo_bukti2"`
	Supervisor     string `json:"supervisor"`
	IsValidated    bool   `json:"is_validated"`
	AttendanceType string `json:"attendance_type"`
}

type AbsensiHandler struct {
	AbsensiService service.AbsensiService
	UserService    service.UserService
	lastDeviceSync sync.Map
}

func NewAbsensiHandler(service service.AbsensiService, userService service.UserService) *AbsensiHandler {
	return &AbsensiHandler{
		AbsensiService: service,
		UserService:    userService,
	}
}

func (h *AbsensiHandler) GetAll1(c *gin.Context) {
	fmt.Println("AbsensiHandler GetAll")
	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	// status := c.Query("status")
	perPageStr := c.Query("per_page")
	per_page, err := strconv.Atoi(perPageStr)
	if err != nil {
		per_page = 0 // default value
	}

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	nameSearch := c.Query("name_search")

	idLeaderStr := c.Query("id_leader")
	idLeader, _ := strconv.Atoi(idLeaderStr)

	activeFilters := c.Query("active_filters")

	showDup := c.Query("show_dup")
	showDupBool, err := strconv.ParseBool(showDup)
	if err != nil {
		showDupBool = false
	}

	result, err := h.AbsensiService.GetAbsensi(c.Request.Context(), page, per_page, dateFrom, dateTo, nameSearch, idLeader, activeFilters, showDupBool)
	// b, _ := json.MarshalIndent(result, "", "  ")
	// fmt.Println(string(b))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed load Absensi data",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *AbsensiHandler) GetAll(c *gin.Context) {
	fmt.Println("AbsensiHandler GetAll")

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

	idLeaderStr := c.Query("id_leader")
	idLeader, _ := strconv.Atoi(idLeaderStr)

	activeFiltersTemp := c.Query("active_filters")
	activeFilters := ""
	hideDup := false
	for _, filter := range strings.Split(activeFiltersTemp, ",") {
		if filter != "Hide Duplicate" {
			if activeFilters != "" {
				activeFilters += ","
			}
			activeFilters += filter
		} else {
			hideDup = true
		}
	}

	// LOAD
	result, err := h.AbsensiService.GetAbsensi(c.Request.Context(), page, perPage, dateFrom, dateTo, nameSearch, idLeader, activeFilters, hideDup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed load Absensi data",
		})
		return
	}

	// PRINT RAW (debug)
	// b, _ := json.MarshalIndent(result, "", "  ")
	// fmt.Println(string(b))

	// GET "data"
	// ambil field "data"
	rawData, ok := result["data"]
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid data format",
		})
		return
	}

	// rawData = []model.Absensi
	records, ok := rawData.([]model.Absensi)
	if !ok {
		fmt.Printf("Invalid type: %T\n", rawData)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid data structure (expected []model.Absensi)",
		})
		return
	}

	var responseData []AbsensiResponse

	for _, r := range records {

		// parse jam_masuk
		var dateStr, dateStrOut, timeIn, timeOut string
		// fmt.Print("Nama:", safePtr(r.Nama))
		late := 0
		if r.JamMasuk != nil && *r.JamMasuk != "" {
			t, _ := time.Parse(time.RFC3339, *r.JamMasuk)
			dateStr = t.Format("02 January 2006")
			timeIn = t.Format("15:04")
			// fmt.Print("  Masuk Date:", dateStr)
			if keterlambatan := safePtr(r.Keterlambatan); keterlambatan != "" {
				if f, err := strconv.ParseFloat(keterlambatan, 64); err == nil {
					late = int(math.Round(f)) // dibulatkan
					// fmt.Print(" Terlambat:", late)
				}
			}

		} else if r.LemburMasuk != nil && *r.LemburMasuk != "" {
			t, _ := time.Parse(time.RFC3339, *r.LemburMasuk)
			dateStr = t.Format("02 January 2006")
			timeIn = t.Format("15:04")
			// fmt.Print(" Masuk Lembur Date:", dateStr)
		}

		if r.JamPulang != nil && *r.JamPulang != "" {
			t, _ := time.Parse(time.RFC3339, *r.JamPulang)
			dateStrOut = t.Format("02 January 2006")
			timeOut = t.Format("15:04")
			// fmt.Print(" Pulang Date:", dateStr)
		} else if r.LemburPulang != nil && *r.LemburPulang != "" {
			t, _ := time.Parse(time.RFC3339, *r.LemburPulang)
			dateStrOut = t.Format("02 January 2006")
			timeOut = t.Format("15:04")
			// fmt.Print(" Pulang Lembur Date:", dateStr)
		}
		// fmt.Println("")

		responseData = append(responseData, AbsensiResponse{
			ID:             safePtr(r.ID),
			Name:           safePtr(r.Nama),
			Date:           dateStr,
			DateOut:        dateStrOut,
			TimeIn:         timeIn,
			TimeOut:        timeOut,
			Status:         safePtr(r.Status),
			LateMinutes:    late,
			Deduction:      safePtr(r.JumlahPotongan),
			Notes:          safePtr(r.Keterangan),
			PICNotes:       safePtr(r.KeteranganYBS),
			Supervisor:     safePtr(r.Atasan),
			PhotoMasuk:     safePtr(r.PhotoMasuk),
			PhotoPulang:    safePtr(r.PhotoPulang),
			PhotoBukti1:    safePtr(r.BuktiPhoto1),
			PhotoBukti2:    safePtr(r.BuktiPhoto2),
			IsValidated:    safePtr(r.ValidasiAtasan) == "1",
			AttendanceType: mapStatus1(safePtr(r.AttendanceType)),
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

func (h *AbsensiHandler) UpdateStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid id",
		})
		return
	}

	var req struct {
		IsValidated bool `json:"is_validated"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
		})
		return
	}

	if err := h.AbsensiService.UpdateValidation(c.Request.Context(), id, req.IsValidated); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to update validation",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"id":           id,
		"is_validated": req.IsValidated,
		"message":      "status validasi berhasil diperbarui",
	})
}

func (h *AbsensiHandler) UpdateHide(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid id",
		})
		return
	}

	if err := h.AbsensiService.UpdateHide(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to update hide row",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"id":      id,
		"message": "Baris berhasil disembunyikan",
	})
}

func (h *AbsensiHandler) InputAbsensi(c *gin.Context) {
	var req model.AbsensiInputRequest
	fmt.Println("InputAbsensi")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request payload",
			"error":   err.Error(),
		})
		fmt.Println("Error", err.Error())
		return
	}

	b, _ := json.MarshalIndent(req, "", "  ")
	fmt.Println(string(b))

	err := h.AbsensiService.InputAbsensi(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to input absensi",
			"error":   err.Error(),
		})
		fmt.Println("Error ", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Absensi berhasil disimpan",
	})
}

func (h *AbsensiHandler) InputAbsenMesin(c *gin.Context) {
	var req struct {
		Nama      string `json:"nama" binding:"required"`
		Timestamp string `json:"timestamp" binding:"required"`
		Status    string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request payload",
			"error":   err.Error(),
		})
		return
	}

	err := h.AbsensiService.InputAbsenMesin(c.Request.Context(), req.Nama, req.Timestamp, req.Status)
	if err != nil {
		if strings.Contains(err.Error(), "user_karyawan not found") {
			log.Printf("[AbsensiHandler] Karyawan tidak ditemukan: '%s' (error: %v)", req.Nama, err)
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": fmt.Sprintf("Karyawan dengan nama '%s' tidak ditemukan di database", req.Nama),
				"error":   err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to save attendance from machine",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Absensi dari mesin berhasil disimpan",
	})
}

func (h *AbsensiHandler) GetLast(c *gin.Context) {
	idParam := c.Query("id_karyawan")
	fmt.Println("id_karyawan", idParam)
	id_karyawan, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid id karyawan ",
		})
		return
	}

	date := c.Query("date_konfirmasi_absen")

	absensi, err := h.AbsensiService.GetLastAbsensi(c.Request.Context(), id_karyawan, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get last absensi",
			"error":   err.Error(),
		})
		return
	}

	fmt.Println("Last Absensi", absensi)

	c.JSON(http.StatusOK, absensi)
}

func (h *AbsensiHandler) KonfirmasiAbsensi(c *gin.Context) {
	var req model.AbsensiKonfirmasi

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request payload",
			"error":   err.Error(),
		})
		return
	}

	err := h.AbsensiService.KonfirmasiAbsensi(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to konfirmasi absensi",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Absensi berhasil dikonfirmasi",
	})
}

func (h *AbsensiHandler) GetAllPerhitungan(c *gin.Context) {
	start := c.Query("start_date")
	end := c.Query("end_date")
	idLeader := c.Query("id_leader")
	idLeaderInt := 0
	//convert to int idLeader
	if idLeader != "" {
		idLeaderInt, _ = strconv.Atoi(idLeader)
	} else {
		idLeaderInt = 2
	}

	if start == "" || end == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "start & end date required",
		})
		return
	}

	data, err := h.AbsensiService.RekapAbsensi(c.Request.Context(), start, end, idLeaderInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	// jsonBytes, _ := json.MarshalIndent(data, "", "  ")
	// fmt.Println("Data send:\n", string(jsonBytes))
	c.JSON(http.StatusOK, data)
}

func (h *AbsensiHandler) GetIndHolidays(c *gin.Context) {
	data, err := h.AbsensiService.GetIndHolidays(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *AbsensiHandler) GetAbsensiByKaryawan(c *gin.Context) {
	namaKaryawan := c.Query("nama")
	//buang + string
	namaKaryawan = strings.ReplaceAll(namaKaryawan, "+", " ")
	fromDate := c.Query("dateFrom")
	toDate := c.Query("dateTo")

	data, err := h.AbsensiService.GetAbsensiAPI(c.Request.Context(), namaKaryawan, fromDate, toDate)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// jsonBytes, _ := json.MarshalIndent(data, "", "  ")
	// fmt.Println("Data send:\n", string(jsonBytes))

	c.JSON(200, data)
}

func (h *AbsensiHandler) InputLembur(c *gin.Context) {
	var req model.AbsensiLembur
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request payload",
			"error":   err.Error(),
		})
		return
	}

	err := h.AbsensiService.InputLembur(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to input lembur",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Lembur berhasil diinput",
	})
}

func (h *AbsensiHandler) InputAbsensiLeader(c *gin.Context) {
	var req model.AbsensiInputLeader
	fmt.Println("InputAbsensiLeader")
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request payload",
			"error":   err.Error(),
		})
		fmt.Println("Error", err.Error())
		return
	}

	err := h.AbsensiService.InputAbsensiLeader(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to input absensi",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Absensi Leader berhasil disimpan",
	})
}

func (h *AbsensiHandler) DeleteAbsensiLeader(c *gin.Context) {
	var req model.AbsensiInputLeader
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request payload",
			"error":   err.Error(),
		})
		fmt.Println("Error", err.Error())
		return
	}

	err := h.AbsensiService.DeleteAbsensiLeader(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to delete absensi",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Absensi Leader berhasil dihapus",
	})
}

func (h *AbsensiHandler) UpdateAbsensiLeader(c *gin.Context) {
	var req model.AbsensiInputLeader
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request payload",
			"error":   err.Error(),
		})
		fmt.Println("Error", err.Error())
		return
	}

	err := h.AbsensiService.UpdateAbsensiLeader(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to update absensi",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Absensi Leader berhasil diupdate",
	})
}

func (h *AbsensiHandler) GetRekapAbsensiByKaryawan(c *gin.Context) {
	namaKaryawan := c.Query("nama")
	//buang + string
	namaKaryawan = strings.ReplaceAll(namaKaryawan, "+", " ")
	fromDate := c.Query("dateFrom")
	toDate := c.Query("dateTo")

	data, err := h.AbsensiService.GetRekapAbsensiByKaryawan(c.Request.Context(), namaKaryawan, fromDate, toDate)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	jsonBytes, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println("Data send:\n", string(jsonBytes))

	c.JSON(200, data)
}

func (h *AbsensiHandler) GetAbsensiSaya(c *gin.Context) {
	namaKaryawan := c.Query("nama")
	fmt.Println("GetAbsensiSaya for", namaKaryawan)
	if namaKaryawan == "" {
		c.JSON(400, gin.H{"error": "nama parameter is required"})
		return
	}
	//buang + string
	namaKaryawan = strings.ReplaceAll(namaKaryawan, "+", " ")

	data, err := h.AbsensiService.GetAbsensiSaya(c.Request.Context(), namaKaryawan)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// jsonBytes, _ := json.MarshalIndent(data, "", "  ")
	// fmt.Println("Data send:\n", string(jsonBytes))

	c.JSON(200, data)
}

func (h *AbsensiHandler) InputSiteReport(c *gin.Context) {
	var req model.AbsensiSiteReport
	fmt.Println("InputSiteReport")
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request payload",
			"error":   err.Error(),
		})
		fmt.Println("Error", err.Error())
		return
	}
	//printout req
	// b, _ := json.MarshalIndent(req, "", "  ")
	// fmt.Println(string(b))

	err := h.AbsensiService.InputSiteReport(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to input site report",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Site Report berhasil disimpan",
	})
}

func (h *AbsensiHandler) GetSiteReports(c *gin.Context) {
	page := c.Query("page")
	perPage := c.Query("per_page")
	pageInt, _ := strconv.Atoi(page)
	perPageInt, _ := strconv.Atoi(perPage)
	if pageInt < 1 {
		pageInt = 1
	}
	if perPageInt < 1 {
		perPageInt = 0
	}
	nameSearch := c.Query("nama")
	dateFrom := c.Query("start_date")
	dateTo := c.Query("end_date")

	data, err := h.AbsensiService.GetSiteReports(c.Request.Context(), pageInt, perPageInt, nameSearch, dateFrom, dateTo)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	jsonBytes, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println("Data send:\n", string(jsonBytes))

	c.JSON(200, data)
}

func (h *AbsensiHandler) GetLemburList(c *gin.Context) {
	page := c.Query("page")
	perPage := c.Query("per_page")

	pageInt, _ := strconv.Atoi(page)
	perPageInt, _ := strconv.Atoi(perPage)

	if pageInt < 1 {
		pageInt = 1
	}
	if perPageInt < 1 {
		perPageInt = 50
	}

	nameSearch := c.Query("nama")
	dateFrom := c.Query("start_date")
	dateTo := c.Query("end_date")
	status := c.Query("status")

	data, total, err := h.AbsensiService.GetLemburList(c.Request.Context(), pageInt, perPageInt, nameSearch, dateFrom, dateTo, status)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data":     data,
		"total":    total,
		"page":     pageInt,
		"per_page": perPageInt,
	})
}

func (h *AbsensiHandler) ApproveLembur(c *gin.Context) {
	var req struct {
		ID      int64  `json:"id" db:"id"`
		Nama    string `json:"nama" db:"nama"`
		Tanggal string `json:"tanggal" db:"tanggal_lembur"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request payload",
			"error":   err.Error(),
		})
		return
	}
	err := h.AbsensiService.ApproveLembur(c.Request.Context(), req.ID, req.Nama, req.Tanggal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to approve lembur",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Lembur berhasil diapprove",
	})
}

func (h *AbsensiHandler) RejectLembur(c *gin.Context) {
	var req struct {
		ID      int64  `json:"id"`
		Nama    string `json:"nama"`
		Tanggal string `json:"tanggal"`
		Catatan string `json:"catatan"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request payload",
			"error":   err.Error(),
		})
		return
	}
	err := h.AbsensiService.RejectLembur(c.Request.Context(), req.ID, req.Nama, req.Tanggal, req.Catatan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to reject lembur",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Lembur berhasil ditolak",
	})
}

func (h *AbsensiHandler) ReviseLembur(c *gin.Context) {
	var req struct {
		ID      int64  `json:"id" db:"id"`
		Nama    string `json:"nama" db:"nama"`
		Tanggal string `json:"tanggal" db:"tanggal_lembur"`
		Catatan string `json:"catatan" db:"keterangan"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request payload",
			"error":   err.Error(),
		})
		return
	}
	err := h.AbsensiService.ReviseLembur(c.Request.Context(), req.ID, req.Nama, req.Tanggal, req.Catatan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to revise lembur",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Lembur berhasil direvisi",
	})
}

func (h *AbsensiHandler) GetLemburDetail(c *gin.Context) {
	nama := c.Query("nama")
	tanggal := c.Query("tanggal")

	if nama == "" || tanggal == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "nama dan tanggal wajib diisi",
		})
		return
	}

	data, err := h.AbsensiService.GetLemburDetail(c.Request.Context(), nama, tanggal)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Data lembur tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *AbsensiHandler) InputDailyReport(c *gin.Context) {
	var req model.AbsensiDailyReport
	fmt.Println("InputDailyReport")
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request payload",
			"error":   err.Error(),
		})
		fmt.Println("Error", err.Error())
		return
	}
	//printout req
	b, _ := json.MarshalIndent(req, "", "  ")
	fmt.Println(string(b))
	err := h.AbsensiService.InputDailyReport(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to input daily report",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Daily Report berhasil disimpan",
	})

	type Job struct {
		ProjectName   string      `json:"project_name"`
		JobType       string      `json:"job_type"`
		JobTitle      string      `json:"job_title"`
		Description   string      `json:"description"`
		Progress      int         `json:"progress"`
		Kendala       string      `json:"kendala"`
		Solusi        string      `json:"solusi"`
		KendalaDoc    interface{} `json:"kendala_doc"`
		SolusiDoc     interface{} `json:"solusi_doc"`
		SaveDocument1 interface{} `json:"save_document_1"`
		SaveDocument2 interface{} `json:"save_document_2"`
		SaveDocument3 interface{} `json:"save_document_3"`
	}
	var jobs []Job

	err1 := json.Unmarshal([]byte(req.PekerjaanList), &jobs)
	if err1 != nil {
		fmt.Println("Error  :", err1)
		return
	}

	var pekerjaanText strings.Builder
	for i, job := range jobs {
		pekerjaanText.WriteString(fmt.Sprintf(
			"\n%d. Project       : %s\n"+
				"   Tipe          : %s\n"+
				"   Judul         : %s\n"+
				"   %s\n"+
				"   Progress      : %d%%\n"+
				"   %s\n"+
				"   %s\n",
			i+1,
			cleanText(job.ProjectName),
			cleanText(job.JobType),
			cleanText(job.JobTitle),

			formatMultiline("Deskripsi", job.Description),
			job.Progress,
			formatMultilineKendala("Kendala", cleanKendala(job.Kendala)),
			formatMultiline("Solusi", job.Solusi),
		))
	}
	// req.Nama = "X"
	message := fmt.Sprintf(
		`📋 Daily Stand Up

👤 Nama    : %s
📍 Lokasi  : %s
📅 Tanggal : %s
⏰ Jam     : %s - %s

🔧 Pekerjaan:%s

%s
`,
		cleanText(req.Nama),
		cleanText(req.LokasiKerja),
		cleanText(req.Tanggal),
		cleanText(req.JamMulai),
		cleanText(req.JamSelesai),
		pekerjaanText.String(),
		formatMultiline("📅 Rencana Besok", req.RencanaBesok),
	)
	message = fmt.Sprintf("<pre>%s</pre>", message)
	fmt.Println(message)
	utils.SendTelegram("7073260553:AAHOps0arqbCXPTeOCMrv4KO7TGpWisMQfs", "-1002131837764", message)
}

func (h *AbsensiHandler) GetDailyReportByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}
	result, err := h.AbsensiService.GetDailyReportByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	jsonBytes, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println("Data send:\n", string(jsonBytes))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, result)
}

func (h *AbsensiHandler) GetDailyReports(c *gin.Context) {
	nama := c.Query("nama")
	fmt.Println("nama", nama)
	page := c.Query("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}
	perPage := c.Query("per_page")
	perPageInt, err := strconv.Atoi(perPage)
	if err != nil {
		perPageInt = 20
	}
	if perPageInt < 1 {
		perPageInt = 0
	}
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	role := c.Query("role")
	fmt.Println("dateFrom", dateFrom)
	fmt.Println("dateTo", dateTo)
	fmt.Println("perPageInt", perPageInt)
	fmt.Println("pageInt", pageInt)
	fmt.Println("role", role)

	result, err := h.AbsensiService.GetDailyReports(c.Request.Context(), pageInt, perPageInt, nama, dateFrom, dateTo, role)
	jsonBytes, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println("Data send:\n", string(jsonBytes))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, result)
}

func cleanText(s string) string {
	s = strings.TrimSpace(s)

	// ubah literal \n jadi newline asli
	s = strings.ReplaceAll(s, "\\n", "\n")

	// normalize newline
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")

	// optional: rapihin jadi 1 baris (kalau mau flat)
	s = strings.ReplaceAll(s, "\n", " ")

	// hilangkan double spasi
	s = strings.Join(strings.Fields(s), " ")

	return s
}
func formatMultiline(label, value string) string {
	const labelWidth = 13 // ⬅️ ini kunci biar rata

	value = strings.TrimSpace(value)

	// convert \n literal → newline asli
	value = strings.ReplaceAll(value, "\\n", "\n")
	value = strings.ReplaceAll(value, "\r\n", "\n")
	value = strings.ReplaceAll(value, "\r", "\n")

	lines := strings.Split(value, "\n")

	// rapihin tiap baris
	for i, line := range lines {
		lines[i] = strings.Join(strings.Fields(line), " ")
	}

	// format label fix width (kiri rata)
	labelFormatted := fmt.Sprintf("%-*s", labelWidth, label)

	// baris pertama
	result := fmt.Sprintf("%s : %s", labelFormatted, lines[0])

	// indent = labelWidth + " : "
	indent := strings.Repeat(" ", labelWidth+6)

	// baris lanjutan
	for i := 1; i < len(lines); i++ {
		if lines[i] != "" {
			result += "\n" + indent + lines[i]
		}
	}

	return result
}
func cleanKendala(s string) string {
	s = strings.TrimSpace(s)

	// ubah literal \n jadi newline asli
	s = strings.ReplaceAll(s, "\\n", "\n")
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")

	lines := strings.Split(s, "\n")
	var cleaned []string

	for _, line := range lines {
		line = strings.Join(strings.Fields(line), " ")
		if line != "" { // ⬅️ buang baris kosong
			cleaned = append(cleaned, line)
		}
	}

	return strings.Join(cleaned, "\n")
}
func formatMultilineKendala(label, value string) string {
	const labelWidth = 13 // ⬅️ ini kunci biar rata

	value = strings.TrimSpace(value)

	// convert \n literal → newline asli
	value = strings.ReplaceAll(value, "\\n", "\n")
	value = strings.ReplaceAll(value, "\r\n", "\n")
	value = strings.ReplaceAll(value, "\r", "\n")

	lines := strings.Split(value, "\n")

	// rapihin tiap baris
	for i, line := range lines {
		lines[i] = strings.Join(strings.Fields(line), " ")
	}

	// format label fix width (kiri rata)
	labelFormatted := fmt.Sprintf("%-*s", labelWidth, label)

	// baris pertama
	result := fmt.Sprintf("%s : %s", labelFormatted, lines[0])

	// indent = labelWidth + " : "
	indent := strings.Repeat(" ", labelWidth+6)

	// baris lanjutan
	for i := 1; i < len(lines); i++ {
		if lines[i] != "" {
			result += "\n" + indent + lines[i]
		}
	}

	return result
}
func cleanTextKeepNewline(s string) string {
	s = strings.TrimSpace(s)

	// ubah literal \n jadi newline asli
	s = strings.ReplaceAll(s, "\\n", "\n")

	// normalize
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")

	// hilangkan spasi berlebih per baris (tanpa hapus newline)
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.Join(strings.Fields(line), " ")
	}

	return strings.Join(lines, "\n")
}

func safe(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	if s, ok := v.(*string); ok && s != nil {
		return *s
	}
	return fmt.Sprintf("%v", v)
}

func safePtr(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

type UpdateAbsensiStatusRequest struct {
	IsValidated bool `json:"is_validated"`
}

func mapStatus1(status string) string {
	switch strings.ToLower(status) {
	case "kantor":
		return "Kantor"
	case "wfh":
		return "WFH"
	case "sakit":
		return "SAKIT"
	case "cuti":
		return "CUTI"
	case "dinas_lapangan":
		return "Dinas Lapangan"
	case "cuti_lapangan":
		return "Cuti Lapangan"
	default:
		return status
	}
}

func (h *AbsensiHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.AbsensiService.GetDashboardStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch dashboard statistics",
		})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *AbsensiHandler) GetAttendanceStats(c *gin.Context) {
	stats, err := h.AbsensiService.GetAttendanceStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch attendance statistics",
		})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *AbsensiHandler) GetIndividualStats(c *gin.Context) {
	karyawanIDStr := c.Query("karyawanID")
	if karyawanIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing karyawanID parameter",
		})
		return
	}
	id, err := strconv.Atoi(karyawanIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid karyawanID parameter",
		})
		return
	}
	stats, err := h.AbsensiService.GetIndividualStats(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch individual statistics",
		})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *AbsensiHandler) HandleADMSHandshake(c *gin.Context) {
	sn := c.Query("SN")
	options := c.Query("options")
	log.Printf("[ADMS] GET handshake/ping from SN: %s, options: %s", sn, options)

	c.Header("Content-Type", "text/plain")
	if options == "all" {
		nowStr := time.Now().Format("2006-01-02 15:04:05")
		// Respond with standard ZK configuration options
		response := "RegistryCode=\r\n" +
			"RequestDelay=30\r\n" +
			"ResponseDelay=30\r\n" +
			"TransInterval=10\r\n" +
			"TransFlag=1111111111\r\n" +
			"Realtime=1\r\n" +
			"SessionID=1\r\n" +
			"ServerTime=" + nowStr + "\r\n"
		c.String(http.StatusOK, response)
		return
	}

	c.String(http.StatusOK, "OK")
}

func (h *AbsensiHandler) HandleADMSUpload(c *gin.Context) {
	sn := c.Query("SN")
	table := c.Query("table")
	log.Printf("[ADMS] POST upload for SN: %s, table: %s", sn, table)

	reqCtx := c.Request.Context()
	if sn != "" {
		tenantID, err := h.UserService.GetTenantIDByDeviceSN(reqCtx, sn)
		if err == nil && tenantID > 0 {
			reqCtx = context.WithValue(reqCtx, "tenantID", tenantID)
			log.Printf("[ADMS] Mapped SN %s to Tenant ID: %d", sn, tenantID)
		} else {
			log.Printf("[ADMS] SN %s not explicitly mapped to any tenant, using default context", sn)
		}
	}

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("[ADMS] Failed to read body: %v", err)
		c.Header("Content-Type", "text/plain")
		c.String(http.StatusBadRequest, "ERROR: failed to read body")
		return
	}

	bodyStr := string(bodyBytes)
	log.Printf("[ADMS] Payload:\n%s", bodyStr)

	if table == "ATTLOG" {
		lines := strings.Split(bodyStr, "\n")
		successCount := 0
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// Fields are separated by tab \t
			fields := strings.Split(line, "\t")
			if len(fields) < 3 {
				log.Printf("[ADMS] Invalid log line: %s", line)
				continue
			}

			pin := strings.TrimSpace(fields[0])
			timestamp := strings.TrimSpace(fields[1])
			statusCode := strings.TrimSpace(fields[2])

			// ZK status codes: 0 = Masuk, 1 = Pulang, 4 = Lembur-Masuk, 5 = Lembur-Pulang
			var status string
			switch statusCode {
			case "0":
				status = "Masuk"
			case "1":
				status = "Pulang"
			case "4":
				status = "Lembur-Masuk"
			case "5":
				status = "Lembur-Pulang"
			default:
				log.Printf("[ADMS] Skipping status code %s (line: %s)", statusCode, line)
				continue
			}

			// Look up employee name by pin_mesin
			employeeName, err := h.AbsensiService.GetKaryawanNameByPin(reqCtx, pin)
			if err != nil {
				log.Printf("[ADMS] Employee not found for PIN: %s (error: %v)", pin, err)
				continue
			}

			// Call InputAbsenMesin using service
			err = h.AbsensiService.InputAbsenMesin(reqCtx, employeeName, timestamp, status)
			if err != nil {
				log.Printf("[ADMS] Failed to save attendance for %s: %v", employeeName, err)
				continue
			}

			successCount++
			log.Printf("[ADMS] Success check-in for %s (PIN: %s) at %s with status %s", employeeName, pin, timestamp, status)
		}
		log.Printf("[ADMS] Processed %d records successfully", successCount)
	}

	c.Header("Content-Type", "text/plain")
	c.String(http.StatusOK, "OK")
}

func (h *AbsensiHandler) HandleADMSGetRequest(c *gin.Context) {
	c.Header("Content-Type", "text/plain")
	sn := c.Query("SN")
	info := c.Query("INFO")

	shouldSync := false
	if info != "" {
		shouldSync = true
	} else if sn != "" {
		if val, ok := h.lastDeviceSync.Load(sn); ok {
			if lastTime, ok := val.(time.Time); ok {
				if time.Since(lastTime) >= 1*time.Minute {
					shouldSync = true
				}
			} else {
				shouldSync = true
			}
		} else {
			shouldSync = true
		}
	}

	if shouldSync {
		if sn != "" {
			h.lastDeviceSync.Store(sn, time.Now())
		}
		// ZK firmware adds +1 hour offset to CONTROL DEVICE SetTime value, so send serverTime - 1 hour
		adjustedTime := time.Now().Add(-1 * time.Hour)
		nowStr := adjustedTime.Format("2006-01-02 15:04:05")
		cmd := fmt.Sprintf("C:101:CONTROL DEVICE SetTime %s\r\n", nowStr)
		log.Printf("[ADMS] Sending periodic 1-min CONTROL DEVICE SetTime (-1h compensated: %s) to SN %s", nowStr, sn)
		c.String(http.StatusOK, cmd)
		return
	}

	c.String(http.StatusOK, "OK")
}

func (h *AbsensiHandler) HandleADMSDeviceCmd(c *gin.Context) {
	sn := c.Query("SN")
	bodyBytes, _ := io.ReadAll(c.Request.Body)
	log.Printf("[ADMS] DeviceCmd Result from SN %s: %s", sn, strings.TrimSpace(string(bodyBytes)))
	c.Header("Content-Type", "text/plain")
	c.String(http.StatusOK, "OK")
}

func (h *AbsensiHandler) GetLateRulesAPI(c *gin.Context) {
	tid, exists := utils.GetTenantID(c)
	if !exists || tid == 0 {
		tid = 1
	}

	rules, err := h.AbsensiService.GetLateRules(c.Request.Context(), tid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve late rules: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, rules)
}

func (h *AbsensiHandler) UpdateLateRuleAPI(c *gin.Context) {
	tid, exists := utils.GetTenantID(c)
	if !exists || tid == 0 {
		tid = 1
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rule ID"})
		return
	}

	var req struct {
		MinMinutes       int     `json:"min_minutes"`
		MaxMinutes       int     `json:"max_minutes"`
		DeductionBase    string  `json:"deduction_base"`
		DeductionPercent float64 `json:"deduction_percent"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Basic validation
	if req.MinMinutes < 0 || req.MaxMinutes < req.MinMinutes {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Batas minimal/maksimal menit tidak valid"})
		return
	}
	if req.DeductionBase != "none" && req.DeductionBase != "uang_makan" && req.DeductionBase != "uang_harian" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dasar potongan tidak valid"})
		return
	}
	if req.DeductionPercent < 0 || req.DeductionPercent > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Persentase potongan harus berada di antara 0 dan 100"})
		return
	}

	rule := model.AbsensiLateRule{
		ID:               id,
		MinMinutes:       req.MinMinutes,
		MaxMinutes:       req.MaxMinutes,
		DeductionBase:    req.DeductionBase,
		DeductionPercent: req.DeductionPercent,
		TenantID:         tid,
	}

	err = h.AbsensiService.UpdateLateRule(c.Request.Context(), rule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update late rule: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Late rule updated successfully"})
}

