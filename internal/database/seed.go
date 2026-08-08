package database

import (
	"database/sql"
	"errors"
	"fmt"
)

type prioritySeed struct {
	ID          int
	Name        string
	Description string
}

type userSeed struct {
	Nama        string
	Email       string
	FotoProfile string
}

type todoSeed struct {
	Title        string
	Description  string
	Slug         string
	IDPriorities int
	UserEmail    string
}

var prioritySeeds = []prioritySeed{
	{ID: 1, Name: "Penting & Mendesak", Description: "Kerjakan sekarang"},
	{ID: 2, Name: "Penting tapi Tidak Mendesak", Description: "Jadwalkan waktu khusus"},
	{ID: 3, Name: "Tidak Penting tapi Mendesak", Description: "Delegasikan jika memungkinkan"},
	{ID: 4, Name: "Tidak Penting & Tidak Mendesak", Description: "Hapus atau tunda dulu"},
}

var userSeeds = []userSeed{
	{Nama: "Budi Santoso", Email: "budi.santoso@gmail.com", FotoProfile: "/uploads/profiles/default.png"},
	{Nama: "Siti Rahayu", Email: "siti.rahayu@gmail.com", FotoProfile: "/uploads/profiles/default.png"},
	{Nama: "Andi Wijaya", Email: "andi.wijaya@gmail.com", FotoProfile: "/uploads/profiles/default.png"},
}

var todoSeeds = []todoSeed{
	{Title: "Mengerjakan Laporan Kinerja", Description: "Laporan bulanan untuk atasan", Slug: "mengerjakan-laporan-kinerja", IDPriorities: 1, UserEmail: "budi.santoso@gmail.com"},
	{Title: "Bayar Tagihan Listrik", Description: "Tagihan listrik bulan ini", Slug: "bayar-tagihan-listrik", IDPriorities: 1, UserEmail: "siti.rahayu@gmail.com"},
	{Title: "Balas Email Klien Penting", Description: "Menindaklanjuti penawaran kerjasama", Slug: "balas-email-klien-penting", IDPriorities: 3, UserEmail: "andi.wijaya@gmail.com"},
	{Title: "Jadwal Meeting Tim Bulanan", Description: "Koordinasi target bulan depan", Slug: "jadwal-meeting-tim-bulanan", IDPriorities: 2, UserEmail: "budi.santoso@gmail.com"},
	{Title: "Menyusun Strategi Q3", Description: "Perencanaan jangka panjang", Slug: "menyusun-strategi-q3", IDPriorities: 2, UserEmail: "siti.rahayu@gmail.com"},
	{Title: "Menata Ulang Arsip Dokumen", Description: "Arsip lama dirapikan", Slug: "menata-ulang-arsip-dokumen", IDPriorities: 4, UserEmail: "andi.wijaya@gmail.com"},
}

// RunSeed mengisi master prioritas, dummy users, dan dummy todos.
// Idempoten: data yang sudah ada (by email / slug) dilewati.
func RunSeed(db *sql.DB) error {
	if err := seedPriorities(db); err != nil {
		return err
	}

	userIDs, err := seedUsers(db)
	if err != nil {
		return err
	}

	return seedTodos(db, userIDs)
}

func seedPriorities(db *sql.DB) error {
	stmt, err := db.Prepare(`
		INSERT INTO priorities (id_priorities, name, description)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE name = VALUES(name), description = VALUES(description)
	`)
	if err != nil {
		return fmt.Errorf("gagal prepare insert priorities: %w", err)
	}
	defer stmt.Close()

	for _, p := range prioritySeeds {
		if _, err := stmt.Exec(p.ID, p.Name, p.Description); err != nil {
			return fmt.Errorf("gagal seed prioritas %d: %w", p.ID, err)
		}
		fmt.Printf("prioritas %d (%s) siap\n", p.ID, p.Name)
	}
	return nil
}

func seedUsers(db *sql.DB) (map[string]int64, error) {
	userIDs := make(map[string]int64)

	for _, u := range userSeeds {
		id, err := ensureUser(db, u)
		if err != nil {
			return nil, err
		}
		userIDs[u.Email] = id
		fmt.Printf("user %s (id=%d) siap\n", u.Email, id)
	}
	return userIDs, nil
}

func ensureUser(db *sql.DB, u userSeed) (int64, error) {
	var id int64
	err := db.QueryRow("SELECT id_users FROM users WHERE email = ?", u.Email).Scan(&id)
	if err == nil {
		return id, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("gagal cek user %s: %w", u.Email, err)
	}

	result, err := db.Exec(
		"INSERT INTO users (nama, email, foto_profile) VALUES (?, ?, ?)",
		u.Nama, u.Email, u.FotoProfile,
	)
	if err != nil {
		return 0, fmt.Errorf("gagal insert user %s: %w", u.Email, err)
	}
	return result.LastInsertId()
}

func seedTodos(db *sql.DB, userIDs map[string]int64) error {
	for _, t := range todoSeeds {
		userID, ok := userIDs[t.UserEmail]
		if !ok {
			return fmt.Errorf("user %s tidak ditemukan untuk todo %s", t.UserEmail, t.Slug)
		}
		if err := ensureTodo(db, t, userID); err != nil {
			return err
		}
		fmt.Printf("todo %s siap\n", t.Slug)
	}
	return nil
}

func ensureTodo(db *sql.DB, t todoSeed, userID int64) error {
	var id int
	err := db.QueryRow("SELECT id_todos FROM todos WHERE slug = ?", t.Slug).Scan(&id)
	if err == nil {
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("gagal cek todo %s: %w", t.Slug, err)
	}

	_, err = db.Exec(
		"INSERT INTO todos (id_users, id_priorities, title, slug, description) VALUES (?, ?, ?, ?, ?)",
		userID, t.IDPriorities, t.Title, t.Slug, t.Description,
	)
	if err != nil {
		return fmt.Errorf("gagal insert todo %s: %w", t.Slug, err)
	}
	return nil
}
