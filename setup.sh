#!/bin/bash

# setup.sh
# Script untuk inisialisasi dan menjalankan aplikasi RT App Backend

echo "🚀 Menjalankan setup aplikasi RT App Backend..."

# 1. Cek dan buat file .env jika belum ada
if [ ! -f .env ]; then
  echo "🔧 Membuat file .env default..."
  cat <<EOF > .env
DB_HOST=db
DB_PORT=3306
DB_USER=root
DB_PASSWORD=secret
DB_NAME=absensi-sppg-db
JWT_SECRET=supersecretkey
JWT_EXPIRE_HOURS=72
EOF
else
  echo "✅ File .env sudah ada."
fi

# 2. Build dan start container dengan Docker Compose
echo "🐳 Menjalankan docker compose..."
if docker compose version >/dev/null 2>&1; then
  docker compose up -d --build
elif docker-compose version >/dev/null 2>&1; then
  docker-compose up -d --build
else
  echo "❌ Error: docker compose atau docker-compose tidak ditemukan!"
  echo "Silakan instal Docker Compose terlebih dahulu."
  exit 1
fi

# 3. Tunggu MySQL siap
echo "⏳ Menunggu MySQL siap..."
until docker exec rt-app-db mysql -uroot -psecret -e "SELECT 1" &> /dev/null
do
  echo -n "."
  sleep 2
done
echo ""
echo "✅ MySQL siap digunakan."

# 4. Jalankan seeder & migrasi
echo "🌱 Menjalankan migrasi dan seeder database..."
docker exec -i rt-app-backend go run cmd/installer/main.go

echo "✅ Migrasi dan Seeder selesai."

# 5. Info akses
echo "🌐 Aplikasi siap dijalankan!"
echo "➡️ Endpoint API: http://localhost:8080"
echo "📂 Struktur direktori: lihat README.md"

exit 0
