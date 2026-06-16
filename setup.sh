// setup.sh placeholder
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
echo "🐳 Menjalankan docker-compose..."
docker-compose up -d --build

# 3. Tunggu MySQL siap
echo "⏳ Menunggu MySQL siap..."
until docker exec rt-app-db mysql -uroot -psecret -e "SELECT 1" &> /dev/null
do
  echo -n "."
  sleep 2
done
echo ""
echo "✅ MySQL siap digunakan."

# 4. Jalankan seeder SQL
echo "🌱 Menjalankan seeder..."
docker cp seed/seeder.sql rt-app-db:/seed.sql
docker exec -i rt-app-db sh -c 'mysql -uroot -p"$MYSQL_ROOT_PASSWORD" "$MYSQL_DATABASE"' < seed/seeder.sql

echo "✅ Seeder selesai."

# 5. Info akses
echo "🌐 Aplikasi siap dijalankan!"
echo "➡️ Endpoint API: http://localhost:8080"
echo "📂 Struktur direktori: lihat README.md"

exit 0
