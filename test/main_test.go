package api_test

import (
	"database/sql"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"go_todolist/internal/config"
	deliveryhttp "go_todolist/internal/delivery/http"
)

var (
	ts      *httptest.Server
	db      *sql.DB
	cfg     *config.Config
	baseURL string
	token   string
)

func TestMain(m *testing.M) {
	if err := chdirToModuleRoot(); err != nil {
		fmt.Println("GAGAL cari root modul:", err)
		os.Exit(1)
	}

	var err error
	cfg, err = config.Load()
	if err != nil {
		fmt.Println("GAGAL load config (cek .env):", err)
		os.Exit(1)
	}

	db, err = cfg.ConnectDB()
	if err != nil {
		fmt.Println("GAGAL koneksi DB:", err)
		os.Exit(1)
	}

	router := deliveryhttp.NewRouter(db, cfg)
	ts = httptest.NewServer(deliveryhttp.LoggingMiddleware(router))
	baseURL = ts.URL + "/api/v1"

	code := m.Run()

	ts.Close()
	db.Close()
	os.Exit(code)
}

// chdirToModuleRoot berpindah ke direktori berisi go.mod agar .env dan
// folder upload (relative path) terbaca saat test dijalankan dari folder test/.
func chdirToModuleRoot() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			if dir != "." {
				return os.Chdir(dir)
			}
			return nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return fmt.Errorf("go.mod tidak ditemukan")
		}
		dir = parent
	}
}
