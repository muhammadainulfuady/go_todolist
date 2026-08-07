# Todo List REST API - Project Specification & Plan

## 🎯 Permasalahan & Tujuan Aplikasi

Aplikasi ini adalah **RESTful API** berbasis **Go (Golang)** yang dirancang untuk membantu pengguna mengelola tugas sehari-hari secara terorganisir berdasarkan skala prioritas, sehingga pengguna tidak lupa dan fokus pada hal yang paling penting.

---

## 🚀 Fitur Utama

### 1. Autentikasi (Passwordless via Email OTP)

- Pengguna mendaftar/login menggunakan Nama Lengkap dan Email.
- Sistem mengirimkan kode OTP ke email pengguna.
- Pengguna melakukan verifikasi OTP untuk mendapatkan token akses (JWT/Session).

### 2. Profil Pengguna

- Menampilkan nama lengkap dan email.
- Mengunggah foto profil kustom.
- Menyediakan _default avatar_ secara otomatis jika pengguna belum mengunggah foto profil.

### 3. Manajemen Prioritas (Matriks Eisenhower - 4 Skala)

Tugas dikelompokkan ke dalam 4 skala prioritas:

1. **Penting & Mendesak**: Kerjakan segera (tenggat waktu dekat, dampak besar).
2. **Penting tapi Tidak Mendesak**: Jadwalkan waktu khusus (dampak jangka panjang).
3. **Tidak Penting tapi Mendesak**: Delegasikan atau selesaikan dengan cepat.
4. **Tidak Penting & Tidak Mendesak**: Tunda atau abaikan.

### 4. Manajemen Tugas (Todo CRUD)

- **Create**: Membuat tugas baru dengan judul, deskripsi, tingkat prioritas, dan opsional lampiran gambar.
- **Read**: Menampilkan daftar tugas milik pengguna yang sedang login.
- **Search & Filter**: Cari tugas berdasarkan judul dan saring (_filter_) berdasarkan prioritas.
- **Update**: Mengubah detail tugas.
- **Toggle Complete (Undo/Redo)**: Mengubah status tugas antara _Selesai_ (`true`) atau _Belum Selesai_ (`false`).
- **Delete**: Menghapus tugas.

---

## 🗄️ Database Design (Skema Tabel)

### 1. Tabel `users`

| Kolom            | Tipe Data    | Constraint                  | Keterangan               |
| :--------------- | :----------- | :-------------------------- | :----------------------- |
| `id_users`       | INT / BIGINT | PRIMARY KEY, AUTO_INCREMENT | Identifier unik pengguna |
| `nama`           | VARCHAR(100) | NOT NULL                    | Nama lengkap pengguna    |
| `email`          | VARCHAR(100) | UNIQUE, NOT NULL            | Email pengguna           |
| `foto_profile`   | VARCHAR(255) | NULLABLE                    | Path/URL foto profil     |
| `otp_code`       | VARCHAR(6)   | NULLABLE                    | Kode OTP sementara       |
| `otp_expires_at` | DATETIME     | NULLABLE                    | Waktu kedaluwarsa OTP    |
| `created_at`     | DATETIME     | DEFAULT CURRENT_TIMESTAMP   | Waktu pendaftaran        |
| `updated_at`     | DATETIME     | DEFAULT CURRENT_TIMESTAMP   | Waktu pembaruan akun     |

### 2. Tabel `priorities` _(Master Table)_

| Kolom           | Tipe Data   | Constraint  | Keterangan                                   |
| :-------------- | :---------- | :---------- | :------------------------------------------- |
| `id_priorities` | INT         | PRIMARY KEY | ID Prioritas (1, 2, 3, 4)                    |
| `name`          | VARCHAR(50) | NOT NULL    | Nama prioritas (misal: "Penting & Mendesak") |
| `description`   | TEXT        | NULLABLE    | Penjelasan skala prioritas                   |

### 3. Tabel `todos`

| Kolom           | Tipe Data    | Constraint                              | Keterangan                                      |
| :-------------- | :----------- | :-------------------------------------- | :---------------------------------------------- |
| `id_todos`      | INT / BIGINT | PRIMARY KEY, AUTO_INCREMENT             | Identifier unik tugas                           |
| `id_users`      | INT / BIGINT | FOREIGN KEY (`users.id_users`)          | Pemilik tugas                                   |
| `id_priorities` | INT          | FOREIGN KEY (`priorities.id_priorities`) | Skala prioritas (1 - 4)                         |
| `title`         | VARCHAR(150) | NOT NULL                                | Judul tugas                                     |
| `slug`          | VARCHAR(200) | UNIQUE, NOT NULL                        | Slug URL (misal: `mengerjakan-pr`, `mengerjakan-pr-1`) |
| `description`   | TEXT         | NULLABLE                                | Detail/catatan tugas                            |
| `image`         | VARCHAR(255) | NULLABLE                                | Path/URL gambar lampiran                        |
| `is_completed`  | BOOLEAN      | DEFAULT FALSE                           | Status penyelesaian tugas                       |
| `created_at`    | DATETIME     | DEFAULT CURRENT_TIMESTAMP               | Waktu tugas dibuat                              |
| `updated_at`    | DATETIME     | DEFAULT CURRENT_TIMESTAMP               | Waktu tugas diperbarui                          |

## 📡 API Endpoint Design (Kontrak API)

Spesifikasi kontrak API lengkap (OpenAPI 3.0) berisi *request body*, *response payload*, parameter, dan skema autentikasi telah didefinisikan secara resmi di file [`api-spec.json`](file:///D:/LATIHAN%20CODING%20SENDIRI/GOLANG/go_todolist/api-spec.json).

#### Konvensi Response Error

Response error (`4xx`/`5xx`) memakai `{ code, status, message }`, dengan tambahan `errors` **hanya** untuk validasi yang melibatkan banyak field sekaligus:

- **Error tunggal / non-field** (404, 401, 500): cukup `message` saja, tanpa `errors`.
  ```json
  { "code": 404, "status": "fail", "message": "Tugas tidak ditemukan" }
  ```
- **Error validasi multi-field**: tambahkan `errors` berupa objek `map[string]string`, satu key per kolom yang gagal.
  ```json
  { "code": 400, "status": "fail", "message": "Validasi gagal",
    "errors": { "title": "Judul tugas tidak boleh kosong",
                "id_priorities": "Prioritas harus dipilih antara skala 1 sampai 4" } }
  ```

Response sukses (`2xx`) tidak berubah: `{ code, status, message, data }`.

---

## ⚙️ Keputusan Teknologi (Final)

| Aspek           | Pilihan                                                            |
| :-------------- | :----------------------------------------------------------------- |
| Bahasa          | Go 1.26                                                              |
| Router          | **`github.com/julienschmidt/httprouter`** (path param `:slug`)        |
| Database        | **MySQL** (driver `github.com/go-sql-driver/mysql`)                  |
| Autentikasi     | OTP email (SMTP) → JWT HS256 (`github.com/golang-jwt/jwt/v5`)        |
| Pengiriman Email| `net/smtp` via konfigurasi `.env` (host, port, username, password)   |
| Slug Generator  | Helper internal (lowercase, tanpa aksara khusus, auto suffix `-1`/`-2`) |
| Upload Media    | Folder `uploads/` dilayani static, max 2MB JPG/PNG (validasi `mimetype`) |
| Validasi Input  | `github.com/go-playground/validator/v10`                              |
| Testing         | `github.com/stretchr/testify`                                        |

### Catatan Endpoint
- **Semua endpoint** (Profile & Todos) dilindungi JWT, kecuali `/auth/*` dan `GET /priorities`.
- `GET /priorities` sengaja **publik**: data master statis (4 skala Eisenhower) yang sama untuk semua pengguna dan berguna sebelum login.

## 📁 Struktur Folder & File

```
go_todolist/
├── main.go                    # Entry point & routing (httprouter)
├── config/                    # config.go — load .env, koneksi MySQL
├── database/                  # migrate.go, schema.sql, seed.sql
├── models/                    # struct User, Priority, Todo
├── repository/                # query DB per entity
├── service/                   # logika bisnis (OTP, slug, upload, email)
├── handler/                   # HTTP handler per endpoint
├── middleware/                # AuthMiddleware (verifikasi JWT)
├── utils/                     # response helper, JWT, slug, validasi, upload
├── uploads/                   # static file:
│   ├── profiles/              #   foto profil pengguna
│   └── todos/                 #   gambar lampiran tugas
├── .env                       # konfigurasi (JANGAN di-commit)
├── .env.example               # template .env yang aman untuk di-commit
├── .gitignore                 # mengecualikan .env, uploads/, *.exe, dll
└── go.mod
```

### File Non-Go yang Dibutuhkan
| File                  | Peran                                                          |
| :-------------------- | :------------------------------------------------------------- |
| `.env`                | Konfigurasi runtime (secret MySQL, JWT, SMTP) — tidak di-commit |
| `.env.example`        | Contoh konfigurasi ringkas (nilai placeholder) untuk developer  |
| `.gitignore`          | Mencegah secret & artefak masuk ke version control              |
| `database/schema.sql` | DDL migrasi tabel `users`, `priorities`, `todos` (idempotent)   |
| `database/seed.sql`   | Seed 4 data master `priorities` (skala Eisenhower 1–4)          |
| `uploads/.gitkeep`    | Menjaga folder upload tetap ada walau kosong di repo            |

## 🔐 Konfigurasi `.env`

| Variabel                | Keterangan                                   |
| :---------------------- | :------------------------------------------- |
| `APP_PORT`              | Port server (default `8080`)                 |
| `DB_HOST`, `DB_PORT`    | Host & port MySQL                            |
| `DB_USER`, `DB_PASSWORD`| Kredensial MySQL                             |
| `DB_NAME`               | Nama database (`go_todolist`)                |
| `JWT_SECRET`            | Secret untuk signing JWT (wajib dirahasiakan)|
| `JWT_EXPIRES_IN`        | Durasi token (contoh `24h`)                  |
| `SMTP_HOST`, `SMTP_PORT`| Host & port SMTP                             |
| `SMTP_USER`, `SMTP_PASS`| Kredensial email pengirim                    |
| `OTP_EXPIRES_IN`        | Masa berlaku OTP (contoh `5m`)               |

> **Keamanan**: `.env` tidak boleh di-commit (masuk `.gitignore`). JWT secret & kredensial SMTP hanya dibaca dari environment.

---

## 🛠️ Rencana Tahapan Eksekusi Development

1. **Fase 1: Database & Core Setup** (Inisialisasi project Go, dependency, koneksi **MySQL**, migrasi DB + seed 4 prioritas)
2. **Fase 2: Master Priority & Auth OTP** (Endpoint OTP via SMTP, token JWT, middleware auth)
3. **Fase 3: Profile & Media Handling** (Upload foto profil & default avatar)
4. **Fase 4: Core Todo CRUD** (Create, Read, Update, Delete, Toggle Status)
5. **Fase 5: Search, Filter, & Testing** (Pengujian seluruh API via HTTP Client / Postman)

