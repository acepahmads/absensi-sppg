# Dockerfile placeholder
# Gunakan image Go resmi
FROM golang:1.21-alpine

# Set environment variable
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Install git dan mysql-client (opsional untuk debug koneksi)
RUN apk add --no-cache git mysql-client

# Set direktori kerja dalam container
WORKDIR /app

# Salin file go.mod dan go.sum terlebih dahulu (agar cache efisien)
COPY go.mod go.sum ./
RUN go mod download

# Salin semua kode
COPY . .

# Build aplikasi Go
RUN go build -o main ./main.go

# Expose port aplikasi
EXPOSE 8080

# Jalankan binary yang sudah di-build
CMD ["./main"]
