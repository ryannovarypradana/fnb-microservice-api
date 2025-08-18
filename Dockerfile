# --- 1. Build Stage ---
# Menggunakan image Go yang resmi sebagai dasar untuk build
FROM golang:1.21-alpine AS builder

# Menentukan argumen yang bisa kita berikan saat build,
# yaitu nama service yang ingin kita build.
ARG SERVICE_NAME

# Set direktori kerja di dalam container
WORKDIR /app

# Copy file go.mod dan go.sum terlebih dahulu untuk caching
COPY go.mod go.sum ./
# Download semua dependensi
RUN go mod download

# Copy seluruh source code proyek
COPY . .

# Build binary Go untuk service yang ditentukan oleh SERVICE_NAME.
# CGO_ENABLED=0 membuat binary yang statis (tidak butuh C library).
# -o /app/server akan menyimpan hasil build dengan nama "server".
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/${SERVICE_NAME}

# --- 2. Final Stage ---
# Menggunakan image Alpine yang sangat kecil sebagai dasar image akhir
FROM alpine:latest

# Set direktori kerja
WORKDIR /root/

# Hanya copy binary yang sudah di-build dari stage sebelumnya.
# Ini membuat image akhir kita sangat kecil dan aman.
COPY --from=builder /app/server .

# Expose port yang akan digunakan oleh aplikasi (opsional, tapi best practice)
# Port spesifik akan ditentukan saat menjalankan container.

# Command untuk menjalankan binary saat container dimulai.
CMD ["./server"]