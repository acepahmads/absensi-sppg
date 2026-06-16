package repository

import (
	"absensi-sppg/internal/model"
	"absensi-sppg/pkg/utils"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type InventoryRepository interface {
	GetAll(page int, limit int, per_page int, date_from string, date_to string, nameSearch string, id_leader int, category string, jenis string) ([]model.Inventory, int, error)
	GetBarangMasuk(page int, limit int, per_page int, date_from string, date_to string, nameSearch string, id_leader int, category string, jenis string) ([]model.InventoryBarangMasuk, int, error)
	GetBarangKeluar(page int, limit int, per_page int, date_from string, date_to string, nameSearch string, id_leader int, category string, jenis string) ([]model.InventoryBarangKeluar, int, error)
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

type inventoryRepository struct {
	db *sqlx.DB
}

func NewInventoryRepository(db *sqlx.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) GetAll(page int, limit int, per_page int, date_from string, date_to string, nameSearch string, id_leader int, category string, jenis string) ([]model.Inventory, int, error) {
	if per_page > 0 {
		limit = per_page
	}
	offset := (page - 1) * limit

	var results []model.Inventory

	query := `
        SELECT 
			COALESCE(i.id, 0)           AS id, 
			COALESCE(i.qr_code, '')     AS qr_code,
			COALESCE(i.nama_barang, '') AS nama_barang,
			COALESCE(i.kategori_barang, '') AS kategori_barang,
			COALESCE(i.jenis_barang, '')    AS jenis_barang,
			COALESCE(i.satuan, '')          AS satuan,
			COALESCE(i.stok_awal, 0)       AS stok_awal,
			COALESCE(i.stok_akhir, 0)      AS stok_akhir,
			COALESCE(i.stok_masuk, 0)      AS stok_masuk,
			COALESCE(i.posisi, '')         AS posisi,
			COALESCE(i.gambar, '')         AS gambar,
			COALESCE(i.barang_masuk, 0)    AS barang_masuk,
			COALESCE(i.barang_keluar, 0)   AS barang_keluar,
			COALESCE(i.harga_beli, 0)      AS harga_beli,
			COALESCE(i.harga_jual, 0)      AS harga_jual,
			COALESCE(i.keterangan, '')     AS keterangan,
			COALESCE(i.created_at, '')    AS created_at
        FROM inventory i
		WHERE 1 = 1
    `

	params := []interface{}{}
	if date_from != "" && date_to != "" {
		if len(date_from) == 10 {
			date_from = date_from + " 00:00:00"
		}
		if len(date_to) == 10 {
			date_to = date_to + " 23:59:59"
		}
		query += ` 
            AND (
				(i.created_at IS NOT NULL AND i.created_at BETWEEN ? AND ?)
            )
        `
		params = append(params, date_from, date_to, date_from, date_to)
	}

	if nameSearch != "" {
		query += `
			AND i.nama_barang LIKE ?
		`
		params = append(params, "%"+nameSearch+"%")
	}

	query += `
		ORDER BY
			COALESCE(
				i.created_at
			) DESC
	    LIMIT ? OFFSET ?
    `
	params = append(params, limit, offset)

	fmt.Println("Query", query, params)

	err := r.db.Select(&results, query, params...)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, 0, err
	}

	fmt.Println("Results", results)

	// QUERY COUNT
	countQuery := `
        SELECT COUNT(*)
        FROM inventory
		WHERE 1 = 1
    `
	countParams := []interface{}{}

	if date_from != "" && date_to != "" {
		if len(date_from) == 10 {
			date_from = date_from + " 00:00:00"
		}
		if len(date_to) == 10 {
			date_to = date_to + " 23:59:59"
		}
		countQuery += `
            AND (
				(i.created_at IS NOT NULL AND i.created_at BETWEEN ? AND ?)
            )
        `
		countParams = append(countParams, date_from, date_to)
	}

	if nameSearch != "" {
		countQuery += `
			AND nama_barang LIKE ?
		`
		countParams = append(countParams, "%"+nameSearch+"%")
	}

	var total int
	err = r.db.Get(&total, countQuery, countParams...)
	if err != nil {
		fmt.Println("Error executing count query:", err)
		return nil, 0, err
	}

	return results, total, nil
}

func (r *inventoryRepository) GetBarangMasuk(page int, limit int, per_page int, date_from string, date_to string, nameSearch string, id_leader int, category string, jenis string) ([]model.InventoryBarangMasuk, int, error) {
	if per_page > 0 {
		limit = per_page
	}
	offset := (page - 1) * limit

	var results []model.InventoryBarangMasuk

	query := `
		SELECT 
			COALESCE(im.id, 0)           AS id,
			COALESCE(im.qr_code, '')     AS qr_code,
			COALESCE(im.tanggal_jam, '') AS tanggal_jam,
			COALESCE(im.nama_barang, '') AS nama_barang,
			COALESCE(im.kategori, '')    AS kategori,
			COALESCE(im.jenis_barang, '') AS jenis_barang,
			COALESCE(im.jumlah, 0)       AS jumlah,
			COALESCE(im.keterangan, '')  AS keterangan
		FROM inventory_masuk im
		WHERE 1 = 1
	`

	params := []interface{}{}
	if date_from != "" && date_to != "" {
		if len(date_from) == 10 {
			date_from = date_from + " 00:00:00"
		}
		if len(date_to) == 10 {
			date_to = date_to + " 23:59:59"
		}
		query += `
			AND (
				(im.tanggal_jam IS NOT NULL AND im.tanggal_jam BETWEEN ? AND ?)
			)
		`
		params = append(params, date_from, date_to)
	}

	if nameSearch != "" {
		query += `
			AND im.nama_barang LIKE ?
		`
		params = append(params, "%"+nameSearch+"%")
	}

	if category != "" {
		query += `
			AND im.kategori = ?
		`
		params = append(params, category)
	}

	if jenis != "" {
		query += `
			AND im.jenis_barang = ?
		`
		params = append(params, jenis)
	}

	query += `
		ORDER BY
			COALESCE(
				im.tanggal_jam,
				im.jenis_barang
			) DESC
	    LIMIT ? OFFSET ? 
	`
	params = append(params, limit, offset)

	err := r.db.Select(&results, query, params...)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, 0, err
	}

	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM inventory_masuk im
		WHERE 1 = 1
	`
	countParams := []interface{}{}
	if date_from != "" && date_to != "" {
		if len(date_from) == 10 {
			date_from = date_from + " 00:00:00"
		}
		if len(date_to) == 10 {
			date_to = date_to + " 23:59:59"
		}
		countQuery += `
			AND (
				(im.tanggal_jam IS NOT NULL AND im.tanggal_jam BETWEEN ? AND ?)
			)
		`
		countParams = append(countParams, date_from, date_to)
	}

	if nameSearch != "" {
		countQuery += `
			AND im.nama_barang LIKE ?
		`
		countParams = append(countParams, "%"+nameSearch+"%")
	}

	if category != "" {
		countQuery += `
			AND im.kategori = ?
		`
		countParams = append(countParams, category)
	}

	if jenis != "" {
		countQuery += `
			AND im.jenis_barang = ?
		`
		countParams = append(countParams, jenis)
	}

	err = r.db.Get(&total, countQuery, countParams...)
	if err != nil {
		fmt.Println("Error executing count query:", err)
		return nil, 0, err
	}

	return results, total, nil
}

func (r *inventoryRepository) GetBarangKeluar(page int, limit int, per_page int, date_from string, date_to string, nameSearch string, id_leader int, category string, jenis string) ([]model.InventoryBarangKeluar, int, error) {
	if per_page > 0 {
		limit = per_page
	}
	offset := (page - 1) * limit

	var results []model.InventoryBarangKeluar

	query := `
		SELECT 
			COALESCE(im.id, 0)           AS id,
			COALESCE(im.qr_code, '')     AS qr_code,
			COALESCE(im.tanggal_jam, '') AS tanggal_jam,
			COALESCE(im.nama_barang, '') AS nama_barang,
			COALESCE(im.kategori, '')    AS kategori,
			COALESCE(im.jenis_barang, '') AS jenis_barang,
			COALESCE(im.jumlah, 0)       AS jumlah,
			COALESCE(im.keterangan, '')  AS keterangan
		FROM inventory_keluar im
		WHERE 1 = 1
	`

	params := []interface{}{}
	if date_from != "" && date_to != "" {
		if len(date_from) == 10 {
			date_from = date_from + " 00:00:00"
		}
		if len(date_to) == 10 {
			date_to = date_to + " 23:59:59"
		}
		query += `
			AND (
				(im.tanggal_jam IS NOT NULL AND im.tanggal_jam BETWEEN ? AND ?)
			)
		`
		params = append(params, date_from, date_to)
	}

	if nameSearch != "" {
		query += `
			AND im.nama_barang LIKE ?
		`
		params = append(params, "%"+nameSearch+"%")
	}

	if category != "" {
		query += `
			AND im.kategori = ?
		`
		params = append(params, category)
	}

	if jenis != "" {
		query += `
			AND im.jenis_barang = ?
		`
		params = append(params, jenis)
	}

	query += `
		ORDER BY
			COALESCE(
				im.tanggal_jam,
				im.jenis_barang
			) DESC
	    LIMIT ? OFFSET ? 
	`
	params = append(params, limit, offset)

	err := r.db.Select(&results, query, params...)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, 0, err
	}

	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM inventory_keluar im
		WHERE 1 = 1
	`
	countParams := []interface{}{}
	if date_from != "" && date_to != "" {
		if len(date_from) == 10 {
			date_from = date_from + " 00:00:00"
		}
		if len(date_to) == 10 {
			date_to = date_to + " 23:59:59"
		}
		countQuery += `
			AND (
				(im.tanggal_jam IS NOT NULL AND im.tanggal_jam BETWEEN ? AND ?)
			)
		`
		countParams = append(countParams, date_from, date_to)
	}

	if nameSearch != "" {
		countQuery += `
			AND im.nama_barang LIKE ?
		`
		countParams = append(countParams, "%"+nameSearch+"%")
	}

	if category != "" {
		countQuery += `
			AND im.kategori = ?
		`
		countParams = append(countParams, category)
	}

	if jenis != "" {
		countQuery += `
			AND im.jenis_barang = ?
		`
		countParams = append(countParams, jenis)
	}

	err = r.db.Get(&total, countQuery, countParams...)
	if err != nil {
		fmt.Println("Error executing count query:", err)
		return nil, 0, err
	}

	return results, total, nil
}

func (r *inventoryRepository) InputBarang(inventory model.Inventory) error {
	gambar, err := utils.SaveBase64ImageInv(
		inventory.Gambar,
		strings.ReplaceAll(inventory.QRCode, " ", "_"),
	)
	if err != nil {
		return err
	}
	query := `
		INSERT INTO inventory (
			qr_code,
			nama_barang,
			kategori_barang,
			jenis_barang,
			satuan,
			stok_awal,
			stok_akhir,
			stok_masuk,
			posisi,
			gambar,
			barang_masuk,
			barang_keluar,
			harga_beli,
			harga_jual,
			keterangan,
			created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = r.db.Exec(query,
		inventory.QRCode,
		inventory.NamaBarang,
		inventory.Kategori,
		inventory.JenisBarang,
		inventory.Satuan,
		inventory.StokAwal,
		inventory.StokAkhir,
		inventory.StokMasuk,
		inventory.Posisi,
		gambar,
		inventory.BarangMasuk,
		inventory.BarangKeluar,
		inventory.HargaBeli,
		inventory.HargaJual,
		inventory.Keterangan,
		inventory.CreatedAt,
	)
	if err != nil {
		fmt.Println("Error executing insert query:", err)
		return err
	}

	return nil
}

func (r *inventoryRepository) UpdateBarang(inventory model.Inventory) error {
	gambar, err := utils.SaveBase64ImageInv(
		inventory.Gambar,
		strings.ReplaceAll(inventory.QRCode, " ", "_"),
	)
	if err != nil {
		return err
	}
	query := `
		UPDATE inventory SET
			qr_code = ?,
			nama_barang = ?,
			kategori_barang = ?,
			jenis_barang = ?,
			satuan = ?,
			stok_awal = ?,
			stok_akhir = ?,
			stok_masuk = ?,
			posisi = ?,`
	if gambar != "" {
		query += `gambar = ?,`
	}
	query += `
			barang_masuk = ?,
			barang_keluar = ?,
			harga_beli = ?,
			harga_jual = ?,
			keterangan = ?
		WHERE id = ?
	`

	_, err = r.db.Exec(query,
		inventory.QRCode,
		inventory.NamaBarang,
		inventory.Kategori,
		inventory.JenisBarang,
		inventory.Satuan,
		inventory.StokAwal,
		inventory.StokAkhir,
		inventory.StokMasuk,
		inventory.Posisi,
		gambar,
		inventory.BarangMasuk,
		inventory.BarangKeluar,
		inventory.HargaBeli,
		inventory.HargaJual,
		inventory.Keterangan,
		inventory.ID,
	)
	if err != nil {
		fmt.Println("Error executing update query:", err)
		return err
	}

	return nil
}

func (r *inventoryRepository) DeleteBarang(id int) error {
	query := `
		DELETE FROM inventory
		WHERE id = ?
	`

	_, err := r.db.Exec(query, id)
	if err != nil {
		fmt.Println("Error executing delete query:", err)
		return err
	}

	return nil
}

func (r *inventoryRepository) InputBarangMasuk(inventory model.InventoryBarangMasuk) error {
	// ID          string `json:"id" db:"id"`
	// QRCode      string `json:"qr_code" db:"qr_code"`
	// TanggalJam  string `json:"tanggal_jam" db:"tanggal_jam"`
	// NamaBarang  string `json:"nama_barang" db:"nama_barang"`
	// Kategori    string `json:"kategori" db:"kategori"`
	// JenisBarang string `json:"jenis_barang" db:"jenis_barang"`
	// Jumlah      int    `json:"jumlah" db:"jumlah"`
	// Keterangan  string `json:"keterangan" db:"keterangan"`

	query := `
		INSERT INTO inventory_masuk (
			qr_code,
			tanggal_jam,
			nama_barang,
			kategori,
			jenis_barang,
			jumlah,
			keterangan
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		inventory.QRCode,
		inventory.TanggalJam,
		inventory.NamaBarang,
		inventory.Kategori,
		inventory.JenisBarang,
		inventory.Jumlah,
		inventory.Keterangan,
	)
	if err != nil {
		fmt.Println("Error executing insert query:", err)
		return err
	}

	return nil
}

func (r *inventoryRepository) UpdateBarangMasuk(inventory model.InventoryBarangMasuk) error {
	query := `
		UPDATE inventory_masuk SET
			qr_code = ?,
			tanggal_jam = ?,
			nama_barang = ?,
			kategori = ?,
			jenis_barang = ?,
			jumlah = ?,
			keterangan = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query,
		inventory.QRCode,
		inventory.TanggalJam,
		inventory.NamaBarang,
		inventory.Kategori,
		inventory.JenisBarang,
		inventory.Jumlah,
		inventory.Keterangan,
		inventory.ID,
	)
	if err != nil {
		fmt.Println("Error executing update query:", err)
		return err
	}

	return nil
}

func (r *inventoryRepository) DeleteBarangMasuk(id int) error {
	query := `
		DELETE FROM inventory_masuk
		WHERE id = ?
	`

	_, err := r.db.Exec(query, id)
	if err != nil {
		fmt.Println("Error executing delete query:", err)
		return err
	}

	return nil
}

func (r *inventoryRepository) InputBarangKeluar(inventory model.InventoryBarangKeluar) error {
	// ID          string `json:"id" db:"id"`
	// QRCode      string `json:"qr_code" db:"qr_code"`
	// TanggalJam  string `json:"tanggal_jam" db:"tanggal_jam"`
	// NamaBarang  string `json:"nama_barang" db:"nama_barang"`
	// Kategori    string `json:"kategori" db:"kategori"`
	// JenisBarang string `json:"jenis_barang" db:"jenis_barang"`
	// Jumlah      int    `json:"jumlah" db:"jumlah"`
	// Keterangan  string `json:"keterangan" db:"keterangan"`

	query := `
		INSERT INTO inventory_keluar (
			qr_code,
			tanggal_jam,
			nama_barang,
			kategori,
			jenis_barang,
			jumlah,
			keterangan
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		inventory.QRCode,
		inventory.TanggalJam,
		inventory.NamaBarang,
		inventory.Kategori,
		inventory.JenisBarang,
		inventory.Jumlah,
		inventory.Keterangan,
	)
	if err != nil {
		fmt.Println("Error executing insert query:", err)
		return err
	}

	return nil
}

func (r *inventoryRepository) UpdateBarangKeluar(inventory model.InventoryBarangKeluar) error {
	query := `
		UPDATE inventory_keluar SET
			qr_code = ?,
			tanggal_jam = ?,
			nama_barang = ?,
			kategori = ?,
			jenis_barang = ?,
			jumlah = ?,
			keterangan = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query,
		inventory.QRCode,
		inventory.TanggalJam,
		inventory.NamaBarang,
		inventory.Kategori,
		inventory.JenisBarang,
		inventory.Jumlah,
		inventory.Keterangan,
		inventory.ID,
	)
	if err != nil {
		fmt.Println("Error executing update query:", err)
		return err
	}

	return nil
}

func (r *inventoryRepository) DeleteBarangKeluar(id int) error {
	query := `
		DELETE FROM inventory_keluar
		WHERE id = ?
	`

	_, err := r.db.Exec(query, id)
	if err != nil {
		fmt.Println("Error executing delete query:", err)
		return err
	}

	return nil
}
