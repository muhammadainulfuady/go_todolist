# Todo List REST API

RESTful API berbasis **Go** untuk mengelola tugas sehari-hari berdasarkan skala prioritas **Matriks Eisenhower** (4 skala). Autentikasi *passwordless* via **OTP email (SMTP)** + token **JWT**, database **MySQL**. Kontrak API lengkap ada di [`api/api-spec.json`](api/api-spec.json) dan detail teknis di [`PLANT.md`](PLANT.md).

## Prasyarat

- Go 1.26+
- MySQL (server aktif di `127.0.0.1:3306`)

## Setup

1. Salin `.env.example` menjadi `.env` dan isi konfigurasinya (DB, `JWT_SECRET`, SMTP Gmail App Password):
   ```powershell
   Copy-Item .env.example .env
   ```
2. Jalankan migrasi & seeding (otomatis membuat database `golang_restful_api_todolist` jika belum ada).

## Perintah Menjalankan

Semua perintah dijalankan dari root project (tanpa `make`):

```powershell
# Menjalankan Server API (http://localhost:8080)
go run cmd/api/main.go

# Migrasi Database (Up) — membuat tabel + mencatat di schema_migrations
go run cmd/migrate/main.go up

# Rollback Database (Down) — menghapus tabel (urutan terbalik) + catatan
go run cmd/migrate/main.go down

# Seeding Data — 4 prioritas Eisenhower, dummy users & todos (idempoten)
go run cmd/migrate/main.go seed
```

Urutan pertama kali: `go run cmd/migrate/main.go up` → `go run cmd/migrate/main.go seed` → `go run cmd/api/main.go`.

## Cek Cepat

```powershell
Invoke-RestMethod http://localhost:8080/api/v1/health
Invoke-RestMethod http://localhost:8080/api/v1/priorities
```
