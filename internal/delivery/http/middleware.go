package http

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"

	"go_todolist/internal/security"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

type ctxKey string

const claimsKey ctxKey = "user_claims"

func ClaimsFromContext(ctx context.Context) *security.Claims {
	if claims, ok := ctx.Value(claimsKey).(*security.Claims); ok {
		return claims
	}
	return nil
}

func AuthMiddleware(secret string, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		header := r.Header.Get("Authorization")
		token, ok := strings.CutPrefix(header, "Bearer ")
		if !ok || strings.TrimSpace(token) == "" {
			logrus.WithField("path", r.URL.Path).Warn("Autentikasi gagal: token tidak ada")
			WriteError(w, http.StatusUnauthorized, "fail", "Autentikasi gagal")
			return
		}

		claims, err := security.ParseToken(secret, strings.TrimSpace(token))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"path":   r.URL.Path,
				"detail": err.Error(),
			}).Warn("Autentikasi gagal: token tidak valid")
			WriteError(w, http.StatusUnauthorized, "fail", "Autentikasi gagal")
			return
		}

		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next(w, r.WithContext(ctx), p)
	}
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	return r.ResponseWriter.Write(b)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w}

		next.ServeHTTP(rec, r)

		if rec.status == 0 {
			rec.status = http.StatusOK
		}

		fields := logrus.Fields{
			"method":   r.Method,
			"path":     r.URL.Path,
			"status":   rec.status,
			"duration": time.Since(start).String(),
		}

		switch {
		case rec.status >= 500:
			logrus.WithFields(fields).Error("HTTP request gagal")
		case rec.status >= 400:
			logrus.WithFields(fields).Warn("HTTP request ditolak")
		default:
			logrus.WithFields(fields).Info("HTTP request masuk")
		}
	})
}
