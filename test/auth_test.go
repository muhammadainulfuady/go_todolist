package api_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequestOtpValidation(t *testing.T) {
	_, resp := doJSON(t, "POST", "/auth/request-otp", "", map[string]string{
		"email": "bukanemail",
	})
	require.Equal(t, 400, resp.Code)
	require.NotEmpty(t, resp.Errors["email"])

	_, resp = doJSON(t, "POST", "/auth/request-otp", "", map[string]string{
		"email": cfg.SeedEmail,
	})
	require.Equal(t, 200, resp.Code, resp.Message)
}

func TestAuthFlow(t *testing.T) {
	email := cfg.SeedEmail

	_, reqResp := doJSON(t, "POST", "/auth/request-otp", "", map[string]string{
		"email": email,
	})
	require.Equal(t, 200, reqResp.Code, reqResp.Message)

	otp := readOTP(t, email)

	code, vresp := doJSON(t, "POST", "/auth/verify-otp", "", map[string]string{
		"email":    email,
		"otp_code": otp,
	})
	require.Equal(t, 200, code, vresp.Message)
	var auth struct {
		Token string `json:"token"`
		User  struct {
			Email string `json:"email"`
		} `json:"user"`
	}
	parseData(t, vresp.Data, &auth)
	require.NotEmpty(t, auth.Token)
	require.Equal(t, email, auth.User.Email)

	code, _ = doJSON(t, "POST", "/auth/verify-otp", "", map[string]string{
		"email":    email,
		"otp_code": otp,
	})
	require.Equal(t, 400, code, "OTP reuse harus 400")

	code, _ = doJSON(t, "POST", "/auth/verify-otp", "", map[string]string{
		"email":    email,
		"otp_code": "999999",
	})
	require.Equal(t, 400, code, "OTP salah harus 400")

	code, _ = doJSON(t, "POST", "/auth/verify-otp", "", map[string]string{
		"email":    "tidakada@example.com",
		"otp_code": "123456",
	})
	require.Equal(t, 404, code, "email tak terdaftar harus 404")
}
