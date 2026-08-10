package api_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProfile(t *testing.T) {
	tk := mustToken(t)

	code, resp := doJSON(t, "GET", "/profile", tk, nil)
	require.Equal(t, 200, code)
	_ = resp

	code, resp = doJSON(t, "GET", "/profile", "", nil)
	require.Equal(t, 401, code, "get profile tanpa token harus 401")
	_ = resp

	body, ctype := multipartBody(t, map[string]string{"nama": "Bu"}, "", "")
	code, raw := doReq(t, "PUT", "/profile", tk, ctype, body)
	ar := parseResp(t, raw)
	require.Equal(t, 400, code)
	require.NotEmpty(t, ar.Errors["nama"])
}
