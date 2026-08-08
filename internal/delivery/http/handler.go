package http

import (
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"

	"go_todolist/internal/repository"
)

func HealthHandler(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if err := db.PingContext(r.Context()); err != nil {
			logrus.WithField("detail", err.Error()).Error("Database tidak tersedia")
			WriteError(w, http.StatusServiceUnavailable, "error", "Database tidak tersedia")
			return
		}
		WriteJSON(w, http.StatusOK, "success", "OK", map[string]string{"status": "healthy"})
	}
}

func PrioritiesHandler(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		priorities, err := repository.NewPriorityRepository(db).FindAll(r.Context())
		if err != nil {
			logrus.WithField("detail", err.Error()).Error("Gagal mengambil daftar prioritas")
			WriteError(w, http.StatusInternalServerError, "error", "Terjadi kesalahan server")
			return
		}
		WriteJSON(w, http.StatusOK, "success", "Berhasil mengambil daftar prioritas", priorities)
	}
}
