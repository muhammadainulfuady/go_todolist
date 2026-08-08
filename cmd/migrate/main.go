package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"go_todolist/internal/config"
	"go_todolist/internal/database"
	"go_todolist/internal/helper"
)

const migrationsDir = "db/migrations"

func main() {
	helper.InitLogger()

	if len(os.Args) < 2 {
		logrus.Info("Penggunaan: go run cmd/migrate/main.go <up|down|seed>")
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		logrus.WithField("detail", err.Error()).Fatal("Gagal load config")
	}

	db, err := cfg.ConnectDB()
	if err != nil {
		logrus.WithField("detail", err.Error()).Fatal("Gagal koneksi database")
	}
	defer db.Close()

	switch os.Args[1] {
	case "up":
		err = database.RunMigrations(db, migrationsDir)
	case "down":
		err = database.DropMigrations(db, migrationsDir)
	case "seed":
		err = database.RunSeed(db, cfg.SeedEmail)
	default:
		logrus.Warn("Perintah tidak dikenal. Gunakan: up | down | seed")
		os.Exit(1)
	}

	if err != nil {
		logrus.WithField("detail", err.Error()).Fatal("Gagal")
	}
}
