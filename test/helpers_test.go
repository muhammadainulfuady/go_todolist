package api_test

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go_todolist/internal/entity"
	"go_todolist/internal/security"
)

type apiResponse struct {
	Code    int               `json:"code"`
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Data    json.RawMessage   `json:"data"`
	Errors  map[string]string `json:"errors"`
}

type todoData struct {
	IDTodos      int64  `json:"id_todos"`
	IDUsers      int64  `json:"id_users"`
	IDPriorities int    `json:"id_priorities"`
	Priority     struct {
		IDPriorities int    `json:"id_priorities"`
		Name         string `json:"name"`
	} `json:"priority"`
	Title       string  `json:"title"`
	Slug        string  `json:"slug"`
	Description *string `json:"description"`
	Image       *string `json:"image"`
	IsCompleted bool    `json:"is_completed"`
}

func doReq(t *testing.T, method, path, token, contentType string, body io.Reader) (int, []byte) {
	t.Helper()
	req, err := http.NewRequest(method, baseURL+path, body)
	require.NoError(t, err, "gagal membuat request")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err, "gagal request %s %s", method, path)
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "gagal baca body")
	return resp.StatusCode, data
}

func doJSON(t *testing.T, method, path, token string, payload any) (int, *apiResponse) {
	t.Helper()
	var body io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		require.NoError(t, err, "gagal marshal payload")
		body = bytes.NewReader(b)
	}
	code, raw := doReq(t, method, path, token, "application/json", body)
	resp := &apiResponse{}
	require.NoError(t, json.Unmarshal(raw, resp), "response bukan JSON valid (%d): %s", code, string(raw))
	return code, resp
}

func parseData(t *testing.T, raw json.RawMessage, dst any) {
	t.Helper()
	require.NoError(t, json.Unmarshal(raw, dst), "gagal parse data")
}

func parseResp(t *testing.T, raw []byte) *apiResponse {
	t.Helper()
	resp := &apiResponse{}
	require.NoError(t, json.Unmarshal(raw, resp), "response bukan JSON valid: %s", string(raw))
	return resp
}

func mustToken(t *testing.T) string {
	t.Helper()
	if token != "" {
		return token
	}
	token = doLogin(t, cfg.SeedEmail)
	return token
}

func doLogin(t *testing.T, email string) string {
	t.Helper()
	_, resp := doJSON(t, "POST", "/auth/request-otp", "", map[string]string{
		"email": email,
	})
	require.Equal(t, 200, resp.Code, resp.Message)

	otp := readOTP(t, email)

	code, vresp := doJSON(t, "POST", "/auth/verify-otp", "", map[string]string{
		"email":    email,
		"otp_code": otp,
	})
	require.Equal(t, 200, code, vresp.Message)
	var auth struct {
		Token string `json:"token"`
	}
	parseData(t, vresp.Data, &auth)
	require.NotEmpty(t, auth.Token, "token kosong setelah verify-otp")
	return auth.Token
}

func readOTP(t *testing.T, email string) string {
	t.Helper()
	var code sql.NullString
	require.NoError(t, db.QueryRow("SELECT otp_code FROM users WHERE email = ?", email).Scan(&code), "gagal baca OTP dari DB")
	require.True(t, code.Valid && code.String != "", "OTP kosong di DB")
	return code.String
}

func genToken(t *testing.T, userID int64, expires time.Duration) string {
	t.Helper()
	tok, err := security.GenerateToken(cfg.JWT.Secret, expires, entity.User{
		IDUsers: userID,
		Nama:    "Budi Santoso",
		Email:   cfg.SeedEmail,
	})
	require.NoError(t, err, "gagal generate token")
	return tok
}

func multipartBody(t *testing.T, fields map[string]string, fileField, filePath string) (io.Reader, string) {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	for k, v := range fields {
		require.NoError(t, w.WriteField(k, v), "gagal tulis field %s", k)
	}
	if fileField != "" {
		fw, err := w.CreateFormFile(fileField, filepath.Base(filePath))
		require.NoError(t, err, "gagal buat form file")
		data, err := os.ReadFile(filePath)
		require.NoError(t, err, "gagal baca file %s", filePath)
		_, err = fw.Write(data)
		require.NoError(t, err, "gagal tulis file")
	}
	require.NoError(t, w.Close(), "gagal close multipart")
	return &buf, w.FormDataContentType()
}

func writeTempFile(t *testing.T, name string, data []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	require.NoError(t, os.WriteFile(path, data, 0o644), "gagal tulis file temp %s", name)
	return path
}

func containsSlug(list []todoData, slug string) bool {
	for _, td := range list {
		if td.Slug == slug {
			return true
		}
	}
	return false
}

var tinyPNG = func() []byte {
	b, err := base64.StdEncoding.DecodeString(
		"iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg==",
	)
	if err != nil {
		panic(err)
	}
	return b
}()
