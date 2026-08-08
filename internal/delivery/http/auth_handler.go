package http

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"go_todolist/internal/model"
	"go_todolist/internal/usecase"
)

type AuthHandler struct {
	usecase usecase.AuthUsecase
}

func NewAuthHandler(uc usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{usecase: uc}
}

func (h *AuthHandler) RequestOtp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req model.RequestOtpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "fail", "Format request tidak valid")
		return
	}

	resp, err := h.usecase.RequestOtp(r.Context(), &req)
	if err != nil {
		writeUsecaseError(w, "/api/v1/auth", err)
		return
	}

	WriteJSON(w, http.StatusOK, "success", "Kode OTP telah berhasil dikirimkan ke email anda", resp)
}

func (h *AuthHandler) VerifyOtp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req model.VerifyOtpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "fail", "Format request tidak valid")
		return
	}

	resp, err := h.usecase.VerifyOtp(r.Context(), &req)
	if err != nil {
		writeUsecaseError(w, "/api/v1/auth", err)
		return
	}

	WriteJSON(w, http.StatusOK, "success", "Login berhasil", resp)
}
