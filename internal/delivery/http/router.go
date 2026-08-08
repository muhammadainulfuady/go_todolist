package http

import (
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func NewRouter(db *sql.DB) *httprouter.Router {
	router := httprouter.New()

	router.GET("/api/v1/health", HealthHandler(db))
	router.GET("/api/v1/priorities", PrioritiesHandler())

	router.ServeFiles("/uploads/*filepath", http.Dir("uploads"))

	return router
}
