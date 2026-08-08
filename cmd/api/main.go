package main

import (
	"fmt"
	"net/http"
	"os"

	"go_todolist/internal/config"
	deliveryhttp "go_todolist/internal/delivery/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Gagal load config: %v\n", err)
		os.Exit(1)
	}

	db, err := cfg.ConnectDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Gagal koneksi database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	router := deliveryhttp.NewRouter(db)

	addr := ":" + cfg.Server.Port
	fmt.Printf("Server berjalan di %s%s\n", cfg.Server.BaseURL, "/api/v1/health")
	if err := http.ListenAndServe(addr, router); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
