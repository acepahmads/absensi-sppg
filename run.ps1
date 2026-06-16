# run.ps1

Write-Host "=================================================="
Write-Host "[SETUP] Menjalankan Setup dan Aplikasi Absensi-SPPG"
Write-Host "=================================================="

# Pilih Mode Running
$mode = Read-Host "Pilih metode running: [D] Docker atau [L] Lokal (Ketik D atau L)"
$mode = $mode.ToUpper()

if ($mode -eq "L") {
    Write-Host "`n🚀 Memulai Setup Secara Lokal (Tanpa Docker)..." -ForegroundColor Cyan
    
    # 1. Cek dan buat folder config jika belum ada
    if (!(Test-Path -Path "config")) {
        Write-Host "[FOLDER] Membuat folder config..."
        New-Item -ItemType Directory -Path "config" | Out-Null
    }

    # 2. Cek dan konfigurasikan file config/.env
    $configEnv = "config/.env"
    $useExisting = "N"
    if (Test-Path -Path $configEnv) {
        Write-Host ""
        Write-Host "Ditemukan konfigurasi database yang sudah ada di $configEnv :" -ForegroundColor Yellow
        Get-Content $configEnv | ForEach-Object { Write-Host "  $_" -ForegroundColor Gray }
        $useExisting = Read-Host "`nApakah Anda ingin menggunakan konfigurasi database di atas? (Y/N)"
        $useExisting = $useExisting.ToUpper()
    }

    if ($useExisting -ne "Y") {
        Write-Host "`n[CONFIG] Konfigurasi Database Lokal Baru:" -ForegroundColor Cyan
        $dbHost = Read-Host "Masukkan DB Host (default: 127.0.0.1)"
        if ($dbHost -eq "") { $dbHost = "127.0.0.1" }
        
        $dbPort = Read-Host "Masukkan DB Port (default: 3306)"
        if ($dbPort -eq "") { $dbPort = "3306" }

        $dbUser = Read-Host "Masukkan DB User (default: root)"
        if ($dbUser -eq "") { $dbUser = "root" }

        $dbPassword = Read-Host "Masukkan DB Password (kosongkan jika tidak ada)"
        
        $dbName = "absensi-sppg-db"

        $configEnvContent = @"
DB_HOST=$dbHost
DB_PORT=$dbPort
DB_USER=$dbUser
DB_PASSWORD=$dbPassword
DB_NAME=$dbName
JWT_SECRET=supersecretkey
APP_PORT=8080
"@
        Set-Content -Path $configEnv -Value $configEnvContent
        Write-Host "✅ File config/.env berhasil diperbarui." -ForegroundColor Green
    } else {
        Write-Host "✅ Menggunakan file config/.env yang sudah ada." -ForegroundColor Green
    }

    # 3. Jalankan database migration dan seeder secara lokal
    Write-Host ""
    Write-Host "[INSTALLER] Menjalankan installer database lokal (GORM AutoMigrate dan Seed)..."
    go run cmd/installer/main.go
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Gagal menjalankan installer database secara lokal. Pastikan MySQL lokal Anda sudah aktif dengan kredensial yang benar." -ForegroundColor Red
        Exit 1
    }

    # 4. Jalankan aplikasi utama secara lokal
    Write-Host ""
    Write-Host "[SERVER] Menjalankan aplikasi utama..." -ForegroundColor Green
    Write-Host "➡️  Akses Web Dashboard: http://localhost:8080" -ForegroundColor Cyan
    Write-Host "Tekan Ctrl+C untuk menghentikan server."
    go run main.go
} else {
    Write-Host "`n🐳 Memulai Setup Dengan Docker Compose..." -ForegroundColor Cyan
    
    # 1. Cek dan buat folder config jika belum ada
    if (!(Test-Path -Path "config")) {
        New-Item -ItemType Directory -Path "config" | Out-Null
    }

    # 2. Cek/buat config/.env jika belum ada
    $configEnv = "config/.env"
    if (!(Test-Path -Path $configEnv)) {
        $configEnvContent = @"
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=secret
DB_NAME=absensi-sppg-db
JWT_SECRET=supersecretkey
APP_PORT=8080
"@
        Set-Content -Path $configEnv -Value $configEnvContent
    }

    # 3. Cek/buat .env di root jika belum ada
    $rootEnv = ".env"
    if (!(Test-Path -Path $rootEnv)) {
        $rootEnvContent = @"
DB_HOST=db
DB_PORT=3306
DB_USER=root
DB_PASSWORD=secret
DB_NAME=absensi-sppg-db
JWT_SECRET=supersecretkey
APP_PORT=8080
"@
        Set-Content -Path $rootEnv -Value $rootEnvContent
    }

    # 4. Jalankan docker-compose
    docker-compose up -d --build
    if ($LASTEXITCODE -ne 0) {
        Write-Host "[ERROR] Gagal menjalankan docker-compose. Pastikan Docker Desktop sudah aktif." -ForegroundColor Red
        Exit 1
    }

    # 5. Tunggu MySQL dan App siap
    Write-Host ""
    Write-Host "[WAIT] Menunggu 10 detik agar database dan backend siap..."
    for ($i = 10; $i -gt 0; $i--) {
        Write-Host -NoNewline "$i.. "
        Start-Sleep -Seconds 1
    }
    Write-Host ""

    # 6. Jalankan database migration dan seeder di dalam container
    Write-Host ""
    Write-Host "[INSTALLER] Menjalankan installer database (GORM AutoMigrate dan Seed)..."
    docker exec -i rt-app-backend go run cmd/installer/main.go
    if ($LASTEXITCODE -ne 0) {
        Write-Host "[WARN] Menunggu tambahan 5 detik untuk percobaan ulang..."
        Start-Sleep -Seconds 5
        docker exec -i rt-app-backend go run cmd/installer/main.go
    }

    if ($LASTEXITCODE -eq 0) {
        Write-Host ""
        Write-Host "[SUCCESS] Migrasi database dan seeding berhasil!"
        Write-Host "=================================================="
        Write-Host "Aplikasi siap digunakan!"
        Write-Host "Akses Web Dashboard: http://localhost:8080"
        Write-Host "=================================================="
    } else {
        Write-Host "[ERROR] Gagal menjalankan installer database di dalam container." -ForegroundColor Red
    }
}
