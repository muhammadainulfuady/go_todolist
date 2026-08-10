package http

import (
	"errors"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"go_todolist/internal/helper"
	"go_todolist/internal/model"
	"go_todolist/internal/usecase"
)

type TodoHandler struct {
	usecase usecase.TodoUsecase
}

func NewTodoHandler(uc usecase.TodoUsecase) *TodoHandler {
	return &TodoHandler{usecase: uc}
}

func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	claims := ClaimsFromContext(r.Context())
	if claims == nil {
		WriteError(w, http.StatusUnauthorized, "fail", "Autentikasi gagal")
		return
	}

	if err := r.ParseMultipartForm(helper.MaxImageSize); err != nil {
		WriteError(w, http.StatusBadRequest, "fail", "Format request tidak valid")
		return
	}

	req := model.TodoCreateRequest{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}
	idPriority, err := strconv.Atoi(r.FormValue("id_priorities"))
	if err != nil {
		req.IDPriorities = 0
	} else {
		req.IDPriorities = idPriority
	}

	var image *multipart.FileHeader
	if file, fh, err := r.FormFile("image"); err == nil {
		file.Close()
		image = fh
	} else if !errors.Is(err, http.ErrMissingFile) {
		WriteError(w, http.StatusBadRequest, "fail", "Format request tidak valid")
		return
	}

	resp, err := h.usecase.Create(r.Context(), claims.IDUsers, req, image)
	if err != nil {
		writeUsecaseError(w, "/api/v1/todos", err)
		return
	}

	WriteJSON(w, http.StatusCreated, "success", "Tugas berhasil dibuat", resp)
}

func (h *TodoHandler) List(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	claims := ClaimsFromContext(r.Context())
	if claims == nil {
		WriteError(w, http.StatusUnauthorized, "fail", "Autentikasi gagal")
		return
	}

	search := r.URL.Query().Get("search")

	var idPriority *int
	if raw := r.URL.Query().Get("id_priorities"); raw != "" {
		val, err := strconv.Atoi(raw)
		if err != nil {
			WriteError(w, http.StatusBadRequest, "fail", "Parameter id_priorities tidak valid")
			return
		}
		idPriority = &val
	}

	resp, err := h.usecase.List(r.Context(), claims.IDUsers, search, idPriority)
	if err != nil {
		writeUsecaseError(w, "/api/v1/todos", err)
		return
	}

	WriteJSON(w, http.StatusOK, "success", "Berhasil mengambil daftar tugas", resp)
}

func (h *TodoHandler) GetBySlug(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	claims := ClaimsFromContext(r.Context())
	if claims == nil {
		WriteError(w, http.StatusUnauthorized, "fail", "Autentikasi gagal")
		return
	}

	resp, err := h.usecase.GetBySlug(r.Context(), claims.IDUsers, params.ByName("slug"))
	if err != nil {
		writeUsecaseError(w, "/api/v1/todos/{slug}", err)
		return
	}

	WriteJSON(w, http.StatusOK, "success", "Berhasil mengambil detail tugas", resp)
}

func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	claims := ClaimsFromContext(r.Context())
	if claims == nil {
		WriteError(w, http.StatusUnauthorized, "fail", "Autentikasi gagal")
		return
	}

	if err := r.ParseMultipartForm(helper.MaxImageSize); err != nil {
		WriteError(w, http.StatusBadRequest, "fail", "Format request tidak valid")
		return
	}

	req := model.TodoUpdateRequest{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}
	if raw := r.FormValue("id_priorities"); raw != "" {
		idPriority, err := strconv.Atoi(raw)
		if err != nil {
			WriteError(w, http.StatusBadRequest, "fail", "Parameter id_priorities tidak valid")
			return
		}
		req.IDPriorities = idPriority
	}

	var image *multipart.FileHeader
	if file, fh, err := r.FormFile("image"); err == nil {
		file.Close()
		image = fh
	} else if !errors.Is(err, http.ErrMissingFile) {
		WriteError(w, http.StatusBadRequest, "fail", "Format request tidak valid")
		return
	}

	resp, err := h.usecase.Update(r.Context(), claims.IDUsers, params.ByName("slug"), req, image)
	if err != nil {
		writeUsecaseError(w, "/api/v1/todos/{slug}", err)
		return
	}

	WriteJSON(w, http.StatusOK, "success", "Tugas berhasil diperbarui", resp)
}

func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	claims := ClaimsFromContext(r.Context())
	if claims == nil {
		WriteError(w, http.StatusUnauthorized, "fail", "Autentikasi gagal")
		return
	}

	err := h.usecase.Delete(r.Context(), claims.IDUsers, params.ByName("slug"))
	if err != nil {
		writeUsecaseError(w, "/api/v1/todos/{slug}", err)
		return
	}

	WriteJSON(w, http.StatusOK, "success", "Tugas berhasil dihapus", nil)
}

func (h *TodoHandler) Toggle(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	claims := ClaimsFromContext(r.Context())
	if claims == nil {
		WriteError(w, http.StatusUnauthorized, "fail", "Autentikasi gagal")
		return
	}

	resp, err := h.usecase.Toggle(r.Context(), claims.IDUsers, params.ByName("slug"))
	if err != nil {
		writeUsecaseError(w, "/api/v1/todos/{slug}/toggle", err)
		return
	}

	WriteJSON(w, http.StatusOK, "success", "Status penyelesaian tugas berhasil diperbarui", resp)
}

