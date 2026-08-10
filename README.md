# Todo List REST API

RESTful API berbasis **Go** untuk mengelola tugas sehari-hari berdasarkan skala prioritas **Matriks Eisenhower** (4 skala). Autentikasi *passwordless* via **OTP email (Resend API)** + token **JWT**, database **MySQL**. Kontrak API lengkap ada di [`api/api-spec.json`](api/api-spec.json) dan detail teknis di [`PLANT.md`](PLANT.md).

## Prasyarat

- Go 1.26+
- MySQL (server aktif di `127.0.0.1:3306`)

## Setup

1. Salin `.env.example` menjadi `.env` dan isi konfigurasinya (DB, `JWT_SECRET`, `RESEND_API_KEY`):
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

## Testing

Integration test ditulis dalam **Go** (`test/`) dengan **testify**, dipecah per-tag mengikuti `api/api-spec.json` (`auth_test.go`, `profile_test.go`, `priorities_test.go`, `health_test.go`, `todos_test.go` + `main_test.go` & `helpers_test.go`).

```powershell
go test ./test/ -v
```

Prasyarat:
- `.env` terisi benar (sama dengan langkah Setup) — test otomatis pindah ke root modul untuk membaca `.env`
- MySQL aktif serta migrasi & seeding sudah dijalankan (`go run cmd/migrate/main.go up` dan `go run cmd/migrate/main.go seed`)

Catatan: `request-otp` pada test login mengirim **1 email OTP nyata** via Resend ke `SEED_EMAIL` (`.env`) setiap kali `mustToken`/`TestAuthFlow` dijalankan.

## Endpoint API

Base path: `/api/v1`. Semua endpoint bertanda 🔒 membutuhkan header `Authorization: Bearer <JWT>`.

| Metode | Path | Keterangan |
| :----- | :--- | :--------- |
| POST | `/auth/request-otp` | Kirim kode OTP via email (auto-create user) |
| POST | `/auth/verify-otp` | Verifikasi OTP → terbitkan JWT |
| GET | `/priorities` | Master 4 prioritas Eisenhower (publik) |
| GET | `/profile` 🔒 | Detail profil dari claims JWT |
| PUT | `/profile` 🔒 | Update nama / upload foto profil (multipart, max 2MB) |
| POST | `/todos` 🔒 | Buat tugas (multipart, image opsional max 2MB) |
| GET | `/todos` 🔒 | Daftar tugas `?search=&id_priorities=` |
| GET | `/todos/{slug}` 🔒 | Detail 1 tugas |
| PUT | `/todos/{slug}` 🔒 | Update partial, regenerasi slug, ganti image |
| DELETE | `/todos/{slug}` 🔒 | Hapus tugas + file gambar |
| PATCH | `/todos/{slug}/toggle` 🔒 | Flip status `is_completed` (undo/redo) |

Kontrak lengkap (request/response, validasi, contoh) ada di [`api/api-spec.json`](api/api-spec.json).
