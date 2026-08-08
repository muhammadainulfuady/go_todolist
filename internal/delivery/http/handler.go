package http

import (
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func HealthHandler(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if err := db.PingContext(r.Context()); err != nil {
			WriteError(w, http.StatusServiceUnavailable, "error", "Database tidak tersedia")
			return
		}
		WriteJSON(w, http.StatusOK, "success", "OK", map[string]string{"status": "healthy"})
	}
}

func PrioritiesHandler() httprouter.Handle {
	return func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		WriteJSON(w, http.StatusOK, "success", "OK", []any{})
	}
}
