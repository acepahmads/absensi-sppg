package model

// {
//   "data": [
//     {
//       "id": "item-1",
//       "qr_code": "LPT-001",
//       "nama_barang": "Laptop Dell Latitude 5420",
//       "kategori": "Aset Tetap",
//       "jenis_barang": "Perangkat Teknologi",
//       "satuan": "Unit",
//       "stok": 5,
//       "posisi": "Gudang A - Rak 3",
//       "keterangan": "Kondisi baik",
//       "tanggal": "2024-01-15",
//       "foto_barang": "",
//       "created_at": "2024-01-15T08:30:00.000Z"
//     }
//   ],
//   "total": 150,
//   "page": 1,
//   "per_page": 20
// }

type Inventory struct {
	ID           string `json:"id" db:"id"`
	QRCode       string `json:"qr_code" db:"qr_code"`
	NamaBarang   string `json:"nama_barang" db:"nama_barang"`
	Kategori     string `json:"kategori" db:"kategori_barang"`
	JenisBarang  string `json:"jenis_barang" db:"jenis_barang"`
	Satuan       string `json:"satuan" db:"satuan"`
	StokAwal     int    `json:"stok_awal" db:"stok_awal"`
	StokAkhir    int    `json:"stok_akhir" db:"stok_akhir"`
	StokMasuk    int    `json:"stok_masuk" db:"stok_masuk"`
	Posisi       string `json:"posisi" db:"posisi"`
	Gambar       string `json:"gambar" db:"gambar"`
	BarangMasuk  int    `json:"barang_masuk" db:"barang_masuk"`
	BarangKeluar int    `json:"barang_keluar" db:"barang_keluar"`
	HargaBeli    int    `json:"harga_beli" db:"harga_beli"`
	HargaJual    int    `json:"harga_jual" db:"harga_jual"`
	Keterangan   string `json:"keterangan" db:"keterangan"`
	CreatedAt    string `json:"created_at" db:"created_at"`
}

type InventoryResponse struct {
	ID           string `json:"id"`
	QRCode       string `json:"qr_code"`
	NamaBarang   string `json:"nama_barang"`
	Kategori     string `json:"kategori"`
	JenisBarang  string `json:"jenis_barang"`
	Satuan       string `json:"satuan"`
	StokAwal     int    `json:"stok_awal"`
	StokAkhir    int    `json:"stok_akhir"`
	StokMasuk    int    `json:"stok_masuk"`
	Posisi       string `json:"posisi"`
	Gambar       string `json:"gambar"`
	BarangMasuk  int    `json:"barang_masuk"`
	BarangKeluar int    `json:"barang_keluar"`
	HargaBeli    int    `json:"harga_beli"`
	HargaJual    int    `json:"harga_jual"`
	Keterangan   string `json:"keterangan"`
	CreatedAt    string `json:"created_at"`
}

type InventoryBarangMasuk struct {
	ID          string `json:"id" db:"id"`
	QRCode      string `json:"qr_code" db:"qr_code"`
	TanggalJam  string `json:"tanggal_jam" db:"tanggal_jam"`
	NamaBarang  string `json:"nama_barang" db:"nama_barang"`
	Kategori    string `json:"kategori" db:"kategori"`
	JenisBarang string `json:"jenis_barang" db:"jenis_barang"`
	Jumlah      int    `json:"jumlah" db:"jumlah"`
	Keterangan  string `json:"keterangan" db:"keterangan"`
}

type InventoryBarangKeluar struct {
	ID          string `json:"id" db:"id"`
	QRCode      string `json:"qr_code" db:"qr_code"`
	TanggalJam  string `json:"tanggal_jam" db:"tanggal_jam"`
	NamaBarang  string `json:"nama_barang" db:"nama_barang"`
	Kategori    string `json:"kategori" db:"kategori"`
	JenisBarang string `json:"jenis_barang" db:"jenis_barang"`
	Jumlah      int    `json:"jumlah" db:"jumlah"`
	Keterangan  string `json:"keterangan" db:"keterangan"`
}
