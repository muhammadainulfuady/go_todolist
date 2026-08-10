package http

import (
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"

	"go_todolist/internal/helper"
	"go_todolist/internal/usecase"
)

type ProfileHandler struct {
	usecase usecase.ProfileUsecase
}

func NewProfileHandler(uc usecase.ProfileUsecase) *ProfileHandler {
	return &ProfileHandler{usecase: uc}
}

func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	claims := ClaimsFromContext(r.Context())
	if claims == nil {
		WriteError(w, http.StatusUnauthorized, "fail", "Autentikasi gagal")
		return
	}

	resp, err := h.usecase.GetProfile(r.Context(), claims.IDUsers)
	if err != nil {
		writeUsecaseError(w, "/api/v1/profile", err)
		return
	}

	WriteJSON(w, http.StatusOK, "success", "Berhasil mengambil data profil", resp)
}

func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	claims := ClaimsFromContext(r.Context())
	if claims == nil {
		WriteError(w, http.StatusUnauthorized, "fail", "Autentikasi gagal")
		return
	}

	if err := r.ParseMultipartForm(helper.MaxImageSize); err != nil {
		WriteError(w, http.StatusBadRequest, "fail", "Format request tidak valid")
		return
	}

	nama := r.FormValue("nama")
	var foto *multipart.FileHeader
	if file, fh, err := r.FormFile("foto_profile"); err == nil {
		file.Close()
		foto = fh
	} else if !errors.Is(err, http.ErrMissingFile) {
		WriteError(w, http.StatusBadRequest, "fail", "Format request tidak valid")
		return
	}

	resp, err := h.usecase.UpdateProfile(r.Context(), claims.IDUsers, nama, foto)
	if err != nil {
		writeUsecaseError(w, "/api/v1/profile", err)
		return
	}

	WriteJSON(w, http.StatusOK, "success", "Profil berhasil diperbarui", resp)
}

func writeUsecaseError(w http.ResponseWriter, path string, err error) {
	var verr *helper.ValidationError
	if errors.As(err, &verr) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code":    http.StatusBadRequest,
			"status":  "fail",
			"message": "Validasi gagal",
			"data":    nil,
			"errors":  verr.FieldErrors,
		})
		return
	}

	switch {
	case errors.Is(err, usecase.ErrEmailNotFound):
		WriteError(w, http.StatusNotFound, "fail", "Email tidak terdaftar")
	case errors.Is(err, usecase.ErrInvalidOtp):
		WriteError(w, http.StatusBadRequest, "fail", "Kode OTP salah atau telah kedaluwarsa")
	case errors.Is(err, usecase.ErrUserNotFound):
		WriteError(w, http.StatusNotFound, "fail", "Profil tidak ditemukan")
	case errors.Is(err, usecase.ErrTodoNotFound):
		WriteError(w, http.StatusNotFound, "fail", "Tugas tidak ditemukan")
	default:
		logrus.WithFields(logrus.Fields{
			"path":   path,
			"detail": err.Error(),
		}).Error("Internal server error pada handler")
		WriteError(w, http.StatusInternalServerError, "error", "Terjadi kesalahan server")
	}
}
