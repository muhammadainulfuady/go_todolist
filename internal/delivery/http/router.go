package http

import (
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"go_todolist/internal/config"
	mailgate "go_todolist/internal/gateway/email"
	"go_todolist/internal/repository"
	"go_todolist/internal/usecase"
)

func NewRouter(db *sql.DB, cfg *config.Config) *httprouter.Router {
	router := httprouter.New()

	userRepo := repository.NewUserRepository(db)
	mailer := mailgate.NewMailer(cfg.ResendAPIKey)
	authUsecase := usecase.NewAuthUsecase(
		userRepo, mailer, cfg.JWT.Secret, cfg.JWT.ExpiresIn, cfg.OTP.ExpiresIn,
	)
	authHandler := NewAuthHandler(authUsecase)

	profileUsecase := usecase.NewProfileUsecase(userRepo, cfg.Server.BaseURL, "uploads/profiles")
	profileHandler := NewProfileHandler(profileUsecase)

	router.GET("/api/v1/health", HealthHandler(db))
	router.GET("/api/v1/priorities", PrioritiesHandler(db))
	router.POST("/api/v1/auth/request-otp", authHandler.RequestOtp)
	router.POST("/api/v1/auth/verify-otp", authHandler.VerifyOtp)
	router.GET("/api/v1/profile", AuthMiddleware(cfg.JWT.Secret, profileHandler.GetProfile))
	router.PUT("/api/v1/profile", AuthMiddleware(cfg.JWT.Secret, profileHandler.UpdateProfile))

	router.ServeFiles("/uploads/*filepath", http.Dir("uploads"))

	return router
}
