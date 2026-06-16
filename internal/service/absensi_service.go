package service

import (
	"absensi-sppg/internal/model"
	"absensi-sppg/internal/repository"
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type AbsensiService struct {
	repo repository.AbsensiRepository
}

func NewAbsensiService(repo repository.AbsensiRepository) AbsensiService {
	return AbsensiService{
		repo: repo,
	}
}

func (s *AbsensiService) GetAbsensi(page int, per_page int, date_from string, date_to string, nameSearch string, id_leader int, activeFilters string, hide_dup bool) (map[string]interface{}, error) {
	limit := per_page

	data, total, err := s.repo.GetAll(page, limit, per_page, date_from, date_to, nameSearch, id_leader, activeFilters, hide_dup)
	if err != nil {
		return nil, err
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	response := map[string]interface{}{
		"data":         data,
		"current_page": page,
		"total":        total,
		"total_pages":  totalPage,
	}

	return response, nil
}

func (s *AbsensiService) UpdateValidation(ctx context.Context, id int64, isValidated bool) error {
	return s.repo.UpdateValidation(ctx, id, isValidated)
}

func (s *AbsensiService) UpdateHide(ctx context.Context, id int64) error {
	return s.repo.UpdateHide(ctx, id)
}

func (s *AbsensiService) InputAbsensi(ctx context.Context, req model.AbsensiInputRequest) error {
	model := model.Absensi{
		Nama:           req.Name,
		AttendanceType: req.AttendanceType,
		JamMasuk:       req.Timestamp,
		GPSLatitude:    req.GPSLatitude,
		GPSLongitude:   req.GPSLongitude,
		LocationName:   req.LocationName,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		DurationDays:   req.DurationDays,
		PhotoMasuk:     req.DocumentPhotoMasuk,
		PhotoPulang:    req.DocumentPhotoPulang,
		BuktiPhoto1:    req.DocumentPhotoBukti1,
		BuktiPhoto2:    req.DocumentPhotoBukti2,
		Status:         req.Status,
		ShiftType:      req.ShiftType,
	}

	return s.repo.Insert(ctx, model)
}

func (s *AbsensiService) GetLastAbsensi(id_karyawan int64, date string) (model.AbsensiKeterlambatan, error) {
	return s.repo.GetLastAbsensi(id_karyawan, date)
}

func (s *AbsensiService) KonfirmasiAbsensi(ctx context.Context, req model.AbsensiKonfirmasi) error {
	return s.repo.KonfirmasiAbsensi(ctx, req)
}

// 🔹 Generate tanggal
func generateDates(start, end time.Time) []string {
	var dates []string
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d.Format("2006-01-02"))
	}
	return dates
}

// 🔹 Rekapan utama
func (s *AbsensiService) RekapAbsensi(start, end string, id_leader int) (*model.RekapResponse, error) {
	startDate, _ := time.Parse("2006-01-02", start)
	endDate, _ := time.Parse("2006-01-02", end)
	dates := generateDates(startDate, endDate)

	employees, err := s.repo.GetEmployees(id_leader)
	if err != nil {
		return nil, err
	}

	absensi, err := s.repo.GetAbsensi(start, end, id_leader)
	if err != nil {
		return nil, err
	}

	holidays, err := s.repo.GetIndHolidays()
	if err != nil {
		return nil, err
	}

	resp := model.RekapResponse{
		Employees:  []model.EmployeeRekap{},
		Attendance: map[string]map[string]model.DayLog{},
	}

	for _, emp := range employees {
		summary := model.SummaryAbsensi{}
		idStr := strconv.Itoa(emp.ID)
		summary.Total = 0
		summary.TotalPotongan = 0
		summary.TotalUangMakan = 0
		summary.JumlahHari = 0
		summary.JumlahHariTA = 0
		summary.JumlahLembur = 0

		resp.Attendance[idStr] = map[string]model.DayLog{}
		// fmt.Println("nama", emp.Name, "ID", emp.ID, "absensi", absensi[emp.ID])
		uangMakan, _ := s.repo.GetUangMakan(emp.Name)
		for _, tgl := range dates {
			// fmt.Println("Nama0", emp.Name, "tgl", tgl, "absensi", tgl, "uang makan", uangMakan, "AttendanceType", absensi[emp.ID][tgl].AttendanceType, "JHari", summary.Total)
			if day, ok := absensi[emp.ID][tgl]; ok {
				uangLembur, approval, err := s.repo.GetAbsensiUangLembur(emp.Name, tgl)
				if err != nil {
					uangLembur = 0
				}
				if approval == 1 && uangLembur > 0 {
					summary.JumlahLembur += uangLembur
				}
				resp.Attendance[idStr][tgl] = model.DayLog{
					Status:               day.Status,
					CheckIn:              day.CheckIn,
					CheckOut:             day.CheckOut,
					Notes:                day.Notes,
					SupervisorValidation: day.ValidasiAtasan == 1,
					JumlahPotongan:       day.JumlahPotongan,
					AttendanceType:       day.AttendanceType,
					UangLembur:           uangLembur,
				}
				incrementSummary(&summary, day.Status)
				statusNorm := strings.ToLower(strings.TrimSpace(day.Status))
				attTypeNorm := strings.ToLower(strings.TrimSpace(day.AttendanceType))
				if statusNorm == "t1" || statusNorm == "terlambat1" {

				} else if statusNorm == "t2" || statusNorm == "terlambat2" {
					summary.TotalPotongan += day.JumlahPotongan
					if attTypeNorm == "wfh" {
						summary.Terlambat2--
					}
				} else if (statusNorm == "t3" || statusNorm == "terlambat3") && day.ValidasiAtasan == 1 {
					summary.Terlambat3--
					summary.JumlahHariTA++
				} else if (statusNorm == "t4" || statusNorm == "terlambat4") && day.ValidasiAtasan == 1 {
					summary.Terlambat4--
					summary.JumlahHariTA++
				} else if (statusNorm == "t3" || statusNorm == "terlambat3" || statusNorm == "t4" || statusNorm == "terlambat4") && day.ValidasiAtasan != 1 {
					summary.TotalPotongan += day.JumlahPotongan
				}
				summary.Total++
				if attTypeNorm == "kantor" && (statusNorm == "kantor" || statusNorm == "hadir" || ((statusNorm == "t1" || statusNorm == "terlambat1") && attTypeNorm == "kantor") || ((statusNorm == "t2" || statusNorm == "terlambat2") && attTypeNorm == "kantor")) {
					summary.TotalUangMakan += uangMakan
				}
				// fmt.Println(emp.Name, day.AttendanceType, tgl, day.CheckIn, "uangMakan", uangMakan, "TotalUangMakan", summary.TotalUangMakan, "TotalPotongan", summary.TotalPotongan, "type", day.AttendanceType, "Status", day.Status, "JumlahPotongan", day.JumlahPotongan, "JumlahHariTA", summary.JumlahHariTA)
			} else {
				/*
					if emp.Name == "FIQIH ARAFAT" {
						if isFriday(tgl) || isSaturday(tgl) {
							resp.Attendance[idStr][tgl] = model.DayLog{
								Status: "",
							}
							continue
						} else {
							if isHoliday(tgl, holidays) {
								resp.Attendance[idStr][tgl] = model.DayLog{
									Status: "",
								}
								continue
							}
							resp.Attendance[idStr][tgl] = model.DayLog{
								Status: "Alpha",
								Notes:  "Tidak Ada Keterangan, Segera di Konfirmasi",
							}

						}
					} else if emp.Name == "NAJWA DESTANIA" {
						if isMonday(tgl) || isSunday(tgl) {
							resp.Attendance[idStr][tgl] = model.DayLog{
								Status: "",
							}
							continue
						} else {
							if isHoliday(tgl, holidays) {
								resp.Attendance[idStr][tgl] = model.DayLog{
									Status: "",
								}
								continue
							}
							resp.Attendance[idStr][tgl] = model.DayLog{
								Status: "Alpha",
								Notes:  "Tidak Ada Keterangan, Segera di Konfirmasi",
							}
						}
					} else {*/
				if isWeekend(tgl) {
					resp.Attendance[idStr][tgl] = model.DayLog{
						Status: "",
					}
					continue
				}
				if isHoliday(tgl, holidays) {
					resp.Attendance[idStr][tgl] = model.DayLog{
						Status: "",
					}
					continue
				}
				resp.Attendance[idStr][tgl] = model.DayLog{
					Status: "Alpha",
					Notes:  "Tidak Ada Keterangan, Segera di Konfirmasi",
				}
				//}
				summary.Alpa++
				summary.Total++
			}
			// fmt.Println("nama1", emp.Name, "tgl", tgl, "absensi", tgl, "uang makan", uangMakan, "AttendanceType", absensi[emp.ID][tgl].AttendanceType, "JHari", summary.Total, "JumlahHariTA", summary.JumlahHariTA)
		}
		summary.TotalBayar = summary.TotalUangMakan - summary.TotalPotongan

		resp.Employees = append(resp.Employees, model.EmployeeRekap{
			ID:       emp.ID,
			Name:     emp.Name,
			Division: emp.Division,
			Summary:  summary,
		})
	}

	return &resp, nil
}

func (s *AbsensiService) GetIndHolidays() (map[string]string, error) {
	return s.repo.GetIndHolidays()
}

func (s *AbsensiService) GetAbsensiAPI(nama string, fromDate string, toDate string) ([]model.AbsensiKaryawan, error) {
	rows, err := s.repo.GetAbsensiByKaryawan(nama, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	result := make([]model.AbsensiKaryawan, 0)

	for _, r := range rows {
		// fmt.Println("r.ID", *r.ID)
		JamMasuk := strPtrToString(r.JamMasuk)
		JamPulang := strPtrToString(r.JamPulang)
		// fmt.Println("r.JamMasuk", JamMasuk)
		// fmt.Println("r.JamPulang", JamPulang)
		// fmt.Println("r.PhotoMasuk", r.PhotoMasuk)
		// fmt.Println("r.PhotoPulang ", r.PhotoPulang)
		status := strPtrToString(r.Status)
		// fmt.Println("status", status)
		// fmt.Println("strToPtr(mapStatus(status)) ", strToPtr(mapStatus(status)))
		mappedStatus := mapStatus(status)
		booleanPhotoMasuk := false
		booleanPhotoPulang := false
		if r.PhotoMasuk != nil && *r.PhotoMasuk != "" {
			booleanPhotoMasuk = true
		}
		if r.PhotoPulang != nil && *r.PhotoPulang != "" {
			booleanPhotoPulang = true
		}
		statusApproval := ""
		uangLembur, approval, err := s.repo.GetAbsensiUangLembur(nama, formatDateFromString(JamMasuk))
		if err != nil {
			uangLembur = 0
		}
		if approval == 0 {
			statusApproval = "pending"
		} else if approval == 1 {
			statusApproval = "approved"
		} else if approval == 2 {
			statusApproval = "rejected"
		} else if approval == 3 {
			statusApproval = "revise"
		}
		// fmt.Println("Nama", nama, "Tanggal", formatDateFromString(JamMasuk), "Uang Lembur", uangLembur, "Approval", approval, "Status Approval", statusApproval)
		resp := model.AbsensiKaryawan{
			ID:                fmt.Sprintf("att-%d", r.ID),
			Date:              formatDateFromString(JamMasuk),
			Status:            strToPtr(mappedStatus),
			SupervisorStatus:  mapSupervisorStatus(status, r.ValidasiAtasan),
			CheckIn:           formatTimeFromString(JamMasuk),
			CheckOut:          formatTimeFromString(JamPulang),
			HasPhotoIn:        booleanPhotoMasuk,
			HasPhotoOut:       booleanPhotoPulang,
			PhotoCheckIn:      strPtrToString(r.PhotoMasuk),
			PhotoCheckOut:     strPtrToString(r.PhotoPulang),
			ProofPhoto:        strPtrToString(r.BuktiPhoto1),
			ReasonTitle:       strPtrToString(r.KeteranganYBS),
			ReasonDescription: strPtrToString(r.Keterangan),
			Deduction:         r.JumlahPotongan,
			OvertimeApproval:  statusApproval,
			OvertimePay:       uangLembur,
		}
		result = append(result, resp)
	}

	return result, nil
}

func (s *AbsensiService) GetDailyReportByID(id int64) (model.AbsensiDailyReport, error) {
	return s.repo.GetDailyReportByID(id)
}

func (s *AbsensiService) InputAbsenMesin(ctx context.Context, nama string, timestamp string, status string) error {
	return s.repo.InputAbsenMesin(ctx, nama, timestamp, status)
}

func incrementSummary(s *model.SummaryAbsensi, status string) {
	norm := strings.ToLower(strings.TrimSpace(status))
	switch norm {
	case "kantor", "hadir":
		s.Kantor++
	case "wfh":
		s.WFH++
	case "t1", "terlambat1":
		s.Terlambat1++
	case "t2", "terlambat2":
		s.Terlambat2++
	case "t3", "terlambat3":
		s.Terlambat3++
	case "t4", "terlambat4":
		s.Terlambat4++
	case "terlambat":
		s.Terlambat1++
	case "sakit":
		s.Sakit++
	case "cuti", "izin":
		s.Cuti++
	case "cuti1/2":
		s.Cuti += 0.5
	case "cuti lapangan", "cutilapangan":
		s.CutiLapangan++
	case "dinas lapangan", "dinas":
		s.Dinas++
	case "alpa", "alpha":
		s.Alpa++
	}
}

func (s *AbsensiService) InputLembur(ctx context.Context, req model.AbsensiLembur) error {
	return s.repo.InputLembur(ctx, req)
}

func (s *AbsensiService) InputAbsensiLeader(ctx context.Context, req model.AbsensiInputLeader) error {
	return s.repo.InputAbsensiLeader(ctx, req)
}

func (s *AbsensiService) DeleteAbsensiLeader(ctx context.Context, req model.AbsensiInputLeader) error {
	return s.repo.DeleteAbsensiLeader(ctx, req)
}

func (s *AbsensiService) UpdateAbsensiLeader(ctx context.Context, req model.AbsensiInputLeader) error {
	return s.repo.UpdateAbsensiLeader(ctx, req)
}

func (s *AbsensiService) GetRekapAbsensiByKaryawan(nama string, fromDate string, toDate string) (model.RekapAbsensiByKaryawan, error) {
	return s.repo.GetRekapAbsensiByKaryawan(nama, fromDate, toDate)
}

func (s *AbsensiService) GetAbsensiSaya(nama string) (model.AbsensiSaya, error) {
	return s.repo.GetAbsensiSaya(nama)
}

func (s *AbsensiService) InputSiteReport(ctx context.Context, req model.AbsensiSiteReport) error {
	return s.repo.InputSiteReport(ctx, req)
}

func (s *AbsensiService) GetSiteReports(page int, per_page int, nama string, fromDate string, toDate string) ([]model.AbsensiSiteReport, error) {
	return s.repo.GetSiteReports(page, per_page, nama, fromDate, toDate)
}

func (s *AbsensiService) GetLemburList(page int, per_page int, nama string, fromDate string, toDate string, status string) ([]model.AbsensiLembur, int, error) {
	return s.repo.GetLemburList(page, per_page, nama, fromDate, toDate, status)
}

func (s *AbsensiService) ApproveLembur(ctx context.Context, id int64, nama string, tanggal string) error {
	return s.repo.ApproveLembur(ctx, id, nama, tanggal)
}

func (s *AbsensiService) RejectLembur(ctx context.Context, id int64, nama string, tanggal string, catatan string) error {
	return s.repo.RejectLembur(ctx, id, nama, tanggal, catatan)
}

func (s *AbsensiService) ReviseLembur(ctx context.Context, id int64, nama string, tanggal string, catatan string) error {
	return s.repo.ReviseLembur(ctx, id, nama, tanggal, catatan)
}

func (s *AbsensiService) GetLemburDetail(nama string, tanggal string) (model.AbsensiLembur, error) {
	return s.repo.GetLemburDetail(nama, tanggal)
}

func (s *AbsensiService) InputDailyReport(ctx context.Context, req model.AbsensiDailyReport) error {
	return s.repo.InputDailyReport(ctx, req)
}

func (s *AbsensiService) GetDailyReports(page int, per_page int, nama string, fromDate string, toDate string, role string) ([]model.AbsensiDailyReport, error) {
	return s.repo.GetDailyReports(page, per_page, nama, fromDate, toDate, role)
}

func isWeekend(dateStr string) bool {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false
	}

	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

func isFriday(dateStr string) bool {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false
	}

	weekday := t.Weekday()
	return weekday == time.Friday
}

func isMonday(dateStr string) bool {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false
	}

	weekday := t.Weekday()
	return weekday == time.Monday
}

func isSaturday(dateStr string) bool {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false
	}

	weekday := t.Weekday()
	return weekday == time.Saturday
}

func isSunday(dateStr string) bool {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false
	}

	weekday := t.Weekday()
	return weekday == time.Sunday
}

func isHoliday(dateStr string, holidays map[string]string) bool {
	_, ok := holidays[dateStr]
	return ok
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("15:04")
}

func formatTimePtr(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format("15:04")
}

func mapSupervisorStatus(status string, v *string) string {
	if v == nil || *v == "" {
		if strings.ToLower(status) == "hadir" || strings.ToLower(status) == "wfh" {
			return "approved"
		}
		return "pending"
	}

	switch strings.ToLower(*v) {
	case "1":
		return "approved"
	case "2":
		return "rejected"
	default:
		return "pending"
	}
}

func mapStatus(status string) string {
	switch strings.ToLower(status) {
	case "t3":
		return "late"
	case "t4":
		return "late"
	case "hadir":
		return "present"
	case "alpha":
		return "alpha"
	case "wfh":
		return "wfh"
	case "sakit":
		return "sick"
	case "cuti":
		return "leave"
	case "cuti1/2":
		return "leave1_2"
	case "dinas lapangan":
		return "field"
	case "cuti lapangan":
		return "sleave"
	default:
		return status
	}
}

func mapStatus1(status string) string {
	switch strings.ToLower(status) {
	case "kantor":
		return "Kantor"
	case "wfh":
		return "WFH"
	case "sakit":
		return "Sakit"
	case "cuti":
		return "Cuti"
	case "dinas_lapangan":
		return "Dinas Lapangan"
	case "cuti_lapangan":
		return "Cuti Lapangan"
	default:
		return status
	}
}

func formatTimeFromString(s string) string {
	if s == "" {
		return ""
	}

	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return ""
	}

	return t.Format("15:04")
}

func formatDateFromString(s string) string {
	if s == "" {
		return ""
	}

	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return ""
	}

	return t.Format("2006-01-02")
}

func parseMySQLDateTime(s *string) (*time.Time, error) {
	if s == nil || *s == "" {
		return nil, nil
	}

	// Support DATETIME & DATETIME(6)
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04:05.999999",
	}

	var err error
	for _, layout := range layouts {
		var t time.Time
		t, err = time.ParseInLocation(layout, *s, time.Local)
		if err == nil {
			return &t, nil
		}
	}

	return nil, err
}

func strPtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func strToPtr(s string) *string {
	return &s
}

func (s *AbsensiService) GetDashboardStats(ctx context.Context) (model.DashboardStats, error) {
	return s.repo.GetDashboardStats(ctx)
}
