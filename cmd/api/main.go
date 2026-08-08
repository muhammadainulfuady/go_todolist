package main

import (
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	"go_todolist/internal/config"
	deliveryhttp "go_todolist/internal/delivery/http"
	"go_todolist/internal/helper"
)

func main() {
	helper.InitLogger()

	cfg, err := config.Load()
	if err != nil {
		logrus.WithField("detail", err.Error()).Fatal("Gagal load config")
	}

	db, err := cfg.ConnectDB()
	if err != nil {
		logrus.WithField("detail", err.Error()).Fatal("Gagal koneksi database")
	}
	defer db.Close()

	router := deliveryhttp.NewRouter(db, cfg)

	addr := ":" + cfg.Server.Port
	logrus.WithFields(logrus.Fields{
		"url":  cfg.Server.BaseURL + "/api/v1/health",
		"port": cfg.Server.Port,
	}).Info("Server berjalan")
	if err := http.ListenAndServe(addr, deliveryhttp.LoggingMiddleware(router)); err != nil {
		logrus.WithField("detail", err.Error()).Fatal("Server error")
	}

	os.Exit(0)
}
