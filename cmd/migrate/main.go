package main

import (
	"fmt"
	"os"

	"go_todolist/internal/config"
	"go_todolist/internal/database"
)

const migrationsDir = "db/migrations"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Penggunaan: go run cmd/migrate/main.go <up|down|seed>")
		os.Exit(1)
	}

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

	switch os.Args[1] {
	case "up":
		err = database.RunMigrations(db, migrationsDir)
	case "down":
		err = database.DropMigrations(db, migrationsDir)
	case "seed":
		err = database.RunSeed(db)
	default:
		fmt.Println("Perintah tidak dikenal. Gunakan: up | down | seed")
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Gagal: %v\n", err)
		os.Exit(1)
	}
}
