package model

type Absensi struct {
	//ka.nama, ka.jam_masuk, ka.jam_pulang, ka.lembur_masuk, ka.lembur_pulang, ka.status, ka.keterlambatan, ka.atasan, ka.validasi_atasan, ka.jumlah_potongan, ka.bukti_photo, ka.keterangan
	ID             *string `gorm:"column:id" json:"id" db:"id"`
	Nama           *string `gorm:"column:nama" json:"nama" db:"nama"`
	JamMasuk       *string `gorm:"column:jam_masuk" json:"jam_masuk" db:"jam_masuk"`
	JamPulang      *string `gorm:"column:jam_pulang" json:"jam_pulang" db:"jam_pulang"`
	LemburMasuk    *string `gorm:"column:lembur_masuk" json:"lembur_masuk" db:"lembur_masuk"`
	LemburPulang   *string `gorm:"column:lembur_pulang" json:"lembur_pulang" db:"lembur_pulang"`
	Status         *string `gorm:"column:status" json:"status" db:"status"`
	Keterlambatan  *string `gorm:"column:keterlambatan" json:"keterlambatan" db:"keterlambatan"`
	JumlahPotongan *string `gorm:"column:jumlah_potongan" json:"jumlah_potongan" db:"jumlah_potongan"`
	Keterangan     *string `gorm:"column:keterangan" json:"keterangan" db:"keterangan"`
	Atasan         *string `gorm:"column:atasan" json:"atasan" db:"atasan"`
	KeteranganYBS  *string `gorm:"column:keterangan_ybs" json:"keterangan_ybs" db:"keterangan_ybs"`
	ValidasiAtasan *string `gorm:"column:validasi_atasan" json:"validasi_atasan" db:"validasi_atasan"`
	BuktiPhoto1    *string `gorm:"column:bukti_photo1" json:"bukti_photo1" db:"bukti_photo1"`
	BuktiPhoto2    *string `gorm:"column:bukti_photo2" json:"bukti_photo2" db:"bukti_photo2"`
	PhotoMasuk     *string `gorm:"column:photo_masuk" json:"photo_masuk" db:"photo_masuk"`
	PhotoPulang    *string `gorm:"column:photo_pulang" json:"photo_pulang" db:"photo_pulang"`

	AttendanceType *string `gorm:"column:attendance_type" json:"attendance_type" binding:"required" db:"attendance_type"`
	GPSLatitude    float64 `gorm:"column:gps_latitude" json:"gps_latitude"`
	GPSLongitude   float64 `gorm:"column:gps_longitude" json:"gps_longitude"`
	LocationName   *string `gorm:"column:location_name" json:"location_name"`
	StartDate      *string `gorm:"column:start_date" json:"start_date"`
	EndDate        *string `gorm:"column:end_date" json:"end_date"`
	DurationDays   float64 `gotm:"column:duration_days" json:"duration_days"`
	ShiftType      *string `gorm:"column:shift_type" json:"shift_type"`
}

type AbsensiInputRequest struct {
	Name                *string `json:"name" binding:"required"`
	AttendanceType      *string `json:"attendance_type" binding:"required"`
	Timestamp           *string `json:"timestamp" binding:"required"`
	GPSLatitude         float64 `json:"gps_latitude"`
	GPSLongitude        float64 `json:"gps_longitude"`
	LocationName        *string `json:"location_name"`
	StartDate           *string `json:"start_date"`
	EndDate             *string `json:"end_date"`
	DurationDays        float64 `json:"duration_days"`
	DocumentPhotoMasuk  *string `json:"document_photo_masuk"`
	DocumentPhotoPulang *string `json:"document_photo_pulang"`
	DocumentPhotoBukti1 *string `json:"document_photo_bukti1"`
	DocumentPhotoBukti2 *string `json:"document_photo_bukti2"`
	Status              *string `json:"status"`
	ShiftType           *string `json:"shift_type"`
}

type AbsensiInputLeader struct {
	Name   string `json:"name"`
	Date   string `json:"date"`
	Status string `json:"status"`
	Notes  string `json:"notes"`
}

type AbsensiKeterlambatan struct {
	ID       string `json:"id" db:"id"`
	Nama     string `json:"nama" db:"nama"`
	Tanggal  string `json:"tanggal" db:"tanggal"`
	JamMasuk string `json:"jam_masuk" db:"jam_masuk"`
	Status   string `json:"status" db:"status"`
	Lokasi   string `json:"lokasi" db:"lokasi"`
}

type AbsensiKonfirmasi struct {
	ID         string `json:"id_absensi" db:"id"`
	Alasan     string `json:"alasan" db:"keterangan_ybs"`
	Keterangan string `json:"keterangan" db:"keterangan"`
	FotoBukti  string `json:"foto_bukti" db:"bukti_photo1"`
}

// rekap
type RekapResponse struct {
	Employees  []EmployeeRekap              `json:"employees"`
	Attendance map[string]map[string]DayLog `json:"attendance"`
}

type EmployeeRekap struct {
	ID       int            `json:"id"`
	Name     string         `json:"name"`
	Division string         `json:"division"`
	Summary  SummaryAbsensi `json:"summary"`
}

type SummaryAbsensi struct {
	Total          int     `json:"total"`
	JumlahHariTA   int     `json:"jumlahHariTA"`
	Kantor         int     `json:"kantor"`
	WFH            int     `json:"wfh"`
	Terlambat1     int     `json:"terlambat1"`
	Terlambat2     int     `json:"terlambat2"`
	Terlambat3     int     `json:"terlambat3"`
	Terlambat4     int     `json:"terlambat4"`
	Alpa           int     `json:"alpa"`
	Sakit          int     `json:"sakit"`
	Cuti           float64 `json:"cuti"`
	CutiLapangan   int     `json:"cutiLapangan"`
	Dinas          int     `json:"dinas"`
	TotalPotongan  float64 `json:"totalPotongan"`
	TotalUangMakan float64 `json:"totalUangMakan"`
	JumlahHari     float64 `json:"jumlahHari"`
	TotalBayar     float64 `json:"totalBayar"`
	JumlahLembur   float64 `json:"jumlahLembur"`
}

type DayLog struct {
	Status               string  `json:"status"`
	CheckIn              string  `json:"checkIn"`
	CheckOut             string  `json:"checkOut"`
	Notes                string  `json:"notes"`
	JumlahPotongan       float64 `json:"jumlahPotongan"`
	ValidasiAtasan       int     `json:"validasiAtasan"`
	SupervisorValidation bool    `json:"supervisorValidation"`
	AttendanceType       string  `json:"attendanceType"`
	UangLembur           float64 `json:"uangLembur"`
}

//end rekap

type AbsensiKaryawan struct {
	ID                string  `json:"id"`
	Date              string  `json:"date"`
	Status            *string `json:"status"`
	SupervisorStatus  string  `json:"supervisor_status"`
	CheckIn           string  `json:"check_in"`
	CheckOut          string  `json:"check_out"`
	HasPhotoIn        bool    `json:"has_photo_in"`
	HasPhotoOut       bool    `json:"has_photo_out"`
	PhotoCheckIn      string  `json:"photo_check_in"`
	PhotoCheckOut     string  `json:"photo_check_out"`
	ProofPhoto        string  `json:"proof_photo"`
	ReasonTitle       string  `json:"reason_title"`
	ReasonDescription string  `json:"reason_description"`
	Deduction         *string `json:"deduction"`
	OvertimeApproval  string  `json:"overtime_approval"`
	OvertimePay       float64 `json:"overtime_pay"`
}

type AbsensiLembur struct {
	ID               int      `json:"id" db:"id"`
	Nama             string   `json:"nama" db:"nama"`
	TanggalLembur    string   `json:"tanggal_lembur" db:"tanggal_lembur"`
	DurationWeekday1 float64  `json:"lembur_weekday_1" db:"lembur_weekday_1"`
	DurationWeekday2 float64  `json:"lembur_weekday_2" db:"lembur_weekday_2"`
	DurationWeekend1 float64  `json:"lembur_weekend_1" db:"lembur_weekend_1"`
	DurationWeekend2 float64  `json:"lembur_weekend_2" db:"lembur_weekend_2"`
	DurationWeekend3 float64  `json:"lembur_weekend_3" db:"lembur_weekend_3"`
	DaftarPekerjaan  string   `json:"daftar_pekerjaan" db:"daftar_pekerjaan"`
	BuktiPersetujuan string   `json:"bukti_persetujuan_atasan" db:"bukti_persetujuan_atasan"`
	BuktiPekerjaan   []string `json:"bukti_pekerjaan" db:"bukti_pekerjaan"`
	Approval         string   `json:"approval" db:"approval"`
	Keterangan       string   `json:"keterangan" db:"keterangan"`
	JumlahBayar      float64  `json:"Jumlah_bayar" db:"Jumlah_bayar"`
}

type RekapAbsensiByKaryawan struct {
	Hadir         int `json:"hadir" db:"hadir"`
	WFH           int `json:"wfh" db:"wfh"`
	CutiTahunan   int `json:"cuti_tahunan" db:"cuti_tahunan"`
	Cuti1_2       int `json:"cuti_1_2" db:"cuti_1_2"`
	DinasLapangan int `json:"dinas_lapangan" db:"dinas_lapangan"`
	Sakit         int `json:"sakit" db:"sakit"`
	CutiLapangan  int `json:"cuti_lapangan" db:"cuti_lapangan"`
}

//	{
//	  "status": "Hadir",
//	  "waktu_checkin": "08:15 WIB",
//	  "lokasi": "Kantor Pusat - Jakarta"
//	}
type AbsensiSaya struct {
	Status      string `json:"status" db:"status"`
	CheckInTime string `json:"waktu_checkin" db:"jam_masuk"`
	Location    string `json:"lokasi" db:"location_name"`
}

type AbsensiSiteReport struct {
	ID                      int    `json:"id" db:"id"`
	IDKaryawan              string `json:"id_karyawan" db:"id_karyawan"`
	Nama                    string `json:"nama" db:"nama"`
	JenisReport             string `json:"jenis_report" db:"jenis_report"`
	NamaSystem              string `json:"nama_system" db:"nama_system"`
	Site                    string `json:"site" db:"site"`
	MaintenanceDay          int    `json:"maintenance_day" db:"maintenance_day"`
	JamMasuk                string `json:"jam_masuk" db:"jam_masuk"`
	JamPulang               string `json:"jam_pulang" db:"jam_pulang"`
	PekerjaanHariIni        string `json:"pekerjaan_hari_ini" db:"pekerjaan_hari_ini"`
	YangDikerjakanEsok      string `json:"yang_dikerjakan_esok" db:"yang_dikerjakan_esok"`
	HasilPekerjaan          string `json:"hasil_pekerjaan" db:"hasil_pekerjaan"`
	Kendala                 string `json:"kendala" db:"kendala"`
	BuktiFoto1              string `json:"bukti_foto_1" db:"bukti_foto_1"`
	BuktiFoto2              string `json:"bukti_foto_2" db:"bukti_foto_2"`
	BuktiFoto3              string `json:"bukti_foto_3" db:"bukti_foto_3"`
	AdaPenggantianSparepart string `json:"ada_penggantian_sparepart" db:"ada_penggantian_sparepart"`
	FotoSparepartSebelum    string `json:"foto_sparepart_sebelum" db:"foto_sparepart_sebelum"`
	FotoSparepartSesudah    string `json:"foto_sparepart_sesudah" db:"foto_sparepart_sesudah"`
	CalibrationAttachment   string `json:"calibration_attachment" db:"calibration_attachment"`
	SubmittedAt             string `json:"submitted_at" db:"submitted_at"`
}

type AbsensiDailyReport struct {
	ID            int    `json:"id" db:"id"`
	Nama          string `json:"nama" db:"nama"`
	LokasiKerja   string `json:"lokasi_kerja" db:"lokasi_kerja"`
	Tanggal       string `json:"tanggal" db:"tanggal"`
	JamMulai      string `json:"jam_mulai" db:"jam_mulai"`
	JamSelesai    string `json:"jam_selesai" db:"jam_selesai"`
	PekerjaanList string `json:"pekerjaan_list" db:"pekerjaan_list"`
	RencanaBesok  string `json:"rencana_besok" db:"rencana_besok"`
	CreatedAt     string `json:"created_at" db:"created_at"`
}

type PekerjaanItem struct {
	ID           int         `json:"id" db:"id"`
	Nama         string      `json:"nama" db:"nama"`
	ProjectName  string      `json:"project_name" db:"project_name"`
	JobType      string      `json:"job_type" db:"job_type"`
	JobTitle     string      `json:"job_title" db:"job_title"`
	Description  string      `json:"description"`
	Progress     int         `json:"progress"`
	Kendala      string      `json:"kendala"`
	Solusi       string      `json:"solusi"`
	KendalaDoc   *FileUpload `json:"kendala_doc"`
	SolusiDoc    *FileUpload `json:"solusi_doc"`
	SaveDocument *FileUpload `json:"save_document"`
}

type FileUpload struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size int    `json:"size"`
	Data string `json:"data"` // base64
}

type DashboardStats struct {
	TotalKaryawan         int      `json:"total_karyawan"`
	TotalLeader           int      `json:"total_leader"`
	AkunTerdaftar         int      `json:"akun_terdaftar"`
	AttendanceRateToday   float64  `json:"attendance_rate_today"`
	TepatWaktuPercent     float64  `json:"tepat_waktu_percent"`
	TerlambatPercent      float64  `json:"terlambat_percent"`
	LemburPercent         float64  `json:"lembur_percent"`
	AbsentPercent         float64  `json:"absent_percent"`
	TotalAbsenBulanIni    int      `json:"total_absen_bulan_ini"`
	TotalLemburBulanIni   int      `json:"total_lembur_bulan_ini"`
	AttendanceRateAverage float64  `json:"attendance_rate_average"`
	RecentActivities      []string `json:"recent_activities"`
}

type AbsensiStatistik struct {
	TotalHadirToday     int               `json:"total_hadir_today"`
	TotalKaryawan       int               `json:"total_karyawan"`
	AttendanceRateToday float64           `json:"attendance_rate_today"`
	TepatWaktuRateToday float64           `json:"tepat_waktu_rate_today"`
	GeofenceRateToday   float64           `json:"geofence_rate_today"`
	TepatWaktuCount     int               `json:"tepat_waktu_count"`
	TerlambatCount      int               `json:"terlambat_count"`
	SafeGeofenceCount   int               `json:"safe_geofence_count"`
	UnsafeGeofenceCount int               `json:"unsafe_geofence_count"`
	Tren7Hari           []TrenHari        `json:"tren_7_hari"`
	DistribusiPeran     []DistribusiPeran `json:"distribusi_peran"`
}

type TrenHari struct {
	Label string `json:"label"`
	Count int    `json:"count"`
}

type DistribusiPeran struct {
	Peran      string  `json:"peran"`
	HadirCount int     `json:"hadir_count"`
	TotalCount int     `json:"total_count"`
	Percentage float64 `json:"percentage"`
}

type KaryawanKehadiranIndividu struct {
	Hadir int `json:"hadir"`
	Sakit int `json:"sakit"`
	Izin  int `json:"izin"`
	Alpha int `json:"alpha"`
	Libur int `json:"libur"`
}

