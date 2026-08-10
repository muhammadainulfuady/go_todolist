package api_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTodoUploadInvalid(t *testing.T) {
	tk := mustToken(t)

	fakeTxt := writeTempFile(t, "fake.txt", []byte("bukan gambar"))
	body, ctype := multipartBody(t, map[string]string{"title": "Tugas TXT", "id_priorities": "1"}, "image", fakeTxt)
	code, raw := doReq(t, "POST", "/todos", tk, ctype, body)
	ar := parseResp(t, raw)
	require.Equal(t, 400, code)
	require.NotEmpty(t, ar.Errors["image"])

	oversize := writeTempFile(t, "oversize.png", bytes.Repeat([]byte{0}, (2<<20)+10))
	body, ctype = multipartBody(t, map[string]string{"title": "Tugas Oversize", "id_priorities": "1"}, "image", oversize)
	code, raw = doReq(t, "POST", "/todos", tk, ctype, body)
	ar = parseResp(t, raw)
	require.Equal(t, 400, code, "image oversize harus 400")
}

func TestTodoCreateValidation(t *testing.T) {
	tk := mustToken(t)

	body, ctype := multipartBody(t, map[string]string{"title": "", "id_priorities": "1"}, "", "")
	code, raw := doReq(t, "POST", "/todos", tk, ctype, body)
	ar := parseResp(t, raw)
	require.Equal(t, 400, code)
	require.NotEmpty(t, ar.Errors["title"])

	body, ctype = multipartBody(t, map[string]string{"title": "Validasi", "id_priorities": "7"}, "", "")
	code, raw = doReq(t, "POST", "/todos", tk, ctype, body)
	ar = parseResp(t, raw)
	require.Equal(t, 400, code)
	require.NotEmpty(t, ar.Errors["id_priorities"])

	body, ctype = multipartBody(t, map[string]string{"title": "Validasi"}, "", "")
	code, raw = doReq(t, "POST", "/todos", tk, ctype, body)
	ar = parseResp(t, raw)
	require.Equal(t, 400, code)
	require.NotEmpty(t, ar.Errors["id_priorities"])
}

func TestTodoCRUD(t *testing.T) {
	tk := mustToken(t)
	stamp := fmt.Sprintf("%d", time.Now().UnixNano())

	// create (dengan image)
	png := writeTempFile(t, "img.png", tinyPNG)
	body, ctype := multipartBody(t, map[string]string{
		"title":         "Tugas Uji " + stamp,
		"description":   "dibuat oleh integration test",
		"id_priorities": "2",
	}, "image", png)
	code, raw := doReq(t, "POST", "/todos", tk, ctype, body)
	require.Equal(t, 201, code, "create harus 201")
	var created todoData
	parseData(t, parseResp(t, raw).Data, &created)
	slug := created.Slug
	require.Equal(t, 2, created.IDPriorities)
	require.NotEmpty(t, created.Priority.Name)
	require.NotNil(t, created.Image)

	// duplicate title -> suffix -1
	body, ctype = multipartBody(t, map[string]string{"title": "Tugas Uji " + stamp, "id_priorities": "1"}, "", "")
	code, raw = doReq(t, "POST", "/todos", tk, ctype, body)
	require.Equal(t, 201, code, "create duplikat harus 201")
	var dup todoData
	parseData(t, parseResp(t, raw).Data, &dup)
	require.Equal(t, slug+"-1", dup.Slug, "slug duplikat harus %s-1", slug)

	// list + search + filter
	code, resp := doJSON(t, "GET", "/todos", tk, nil)
	require.Equal(t, 200, code, "list harus 200")
	var list []todoData
	parseData(t, resp.Data, &list)
	require.True(t, containsSlug(list, slug), "list tidak memuat slug %s", slug)

	code, resp = doJSON(t, "GET", "/todos?search=Tugas%20Uji", tk, nil)
	require.Equal(t, 200, code, "search harus 200")
	parseData(t, resp.Data, &list)
	require.GreaterOrEqual(t, len(list), 2, "search 'Tugas Uji' harus >= 2")

	code, resp = doJSON(t, "GET", "/todos?id_priorities=2", tk, nil)
	require.Equal(t, 200, code, "filter harus 200")
	parseData(t, resp.Data, &list)
	for _, td := range list {
		require.Equal(t, 2, td.IDPriorities, "filter bocor: ada id_priorities %d", td.IDPriorities)
	}

	// detail
	code, resp = doJSON(t, "GET", "/todos/"+slug, tk, nil)
	require.Equal(t, 200, code, "detail harus 200")
	var detail todoData
	parseData(t, resp.Data, &detail)
	require.Equal(t, slug, detail.Slug)
	require.NotEmpty(t, detail.Priority.Name)

	// detail lintas user -> 404
	otherToken := genToken(t, 2, time.Hour)
	code, _ = doJSON(t, "GET", "/todos/"+slug, otherToken, nil)
	require.Equal(t, 404, code, "detail lintas-user harus 404")

	// update title -> slug baru
	body, ctype = multipartBody(t, map[string]string{"title": "Tugas Uji Update " + stamp}, "", "")
	code, raw = doReq(t, "PUT", "/todos/"+slug, tk, ctype, body)
	require.Equal(t, 200, code, "update harus 200")
	var updated struct {
		Slug string `json:"slug"`
	}
	parseData(t, parseResp(t, raw).Data, &updated)
	require.NotEqual(t, slug, updated.Slug, "slug harus berubah setelah update judul")
	slug = updated.Slug

	// update partial description -> title & prioritas tetap
	body, ctype = multipartBody(t, map[string]string{"description": "catatan revisi"}, "", "")
	code, raw = doReq(t, "PUT", "/todos/"+slug, tk, ctype, body)
	require.Equal(t, 200, code, "update partial harus 200")
	code, resp = doJSON(t, "GET", "/todos/"+slug, tk, nil)
	require.Equal(t, 200, code)
	parseData(t, resp.Data, &detail)
	require.Equal(t, "Tugas Uji Update "+stamp, detail.Title)
	require.Equal(t, 2, detail.IDPriorities)
	require.NotNil(t, detail.Description)

	// toggle
	code, resp = doJSON(t, "PATCH", "/todos/"+slug+"/toggle", tk, nil)
	require.Equal(t, 200, code, "toggle harus 200")
	var toggled struct {
		IsCompleted bool `json:"is_completed"`
	}
	parseData(t, resp.Data, &toggled)
	require.True(t, toggled.IsCompleted, "toggle harus true")
	code, resp = doJSON(t, "PATCH", "/todos/"+slug+"/toggle", tk, nil)
	require.Equal(t, 200, code)
	parseData(t, resp.Data, &toggled)
	require.False(t, toggled.IsCompleted, "toggle kedua harus false")

	// delete -> 200, lalu get 404
	code, resp = doJSON(t, "DELETE", "/todos/"+slug, tk, nil)
	require.Equal(t, 200, code, "delete harus 200")
	code, _ = doJSON(t, "GET", "/todos/"+slug, tk, nil)
	require.Equal(t, 404, code, "get setelah delete harus 404")
}

func TestTodoAuth401(t *testing.T) {
	code, _ := doJSON(t, "GET", "/todos", "", nil)
	require.Equal(t, 401, code, "list tanpa token harus 401")

	code, _ = doJSON(t, "POST", "/todos", "Bearer-invalid", nil)
	require.Equal(t, 401, code, "token invalid harus 401")

	expired := genToken(t, 1, -time.Minute)
	code, _ = doJSON(t, "GET", "/todos", expired, nil)
	require.Equal(t, 401, code, "token expired harus 401")

	code, _ = doJSON(t, "GET", "/todos", "eyJhbGciOiJIUzI1NiJ9.abc.def", nil)
	require.Equal(t, 401, code, "token malformed harus 401")

	wrongSig := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZF91c2VycyI6MSwibmFtYSI6IkJ1ZGkiLCJleHAiOjk5OTk5OTk5OTl9.ZmFrZS1zaWduYXR1cmU"
	code, _ = doJSON(t, "GET", "/todos", wrongSig, nil)
	require.Equal(t, 401, code, "token signature salah harus 401")
}

func TestTodoDeleteNotFound(t *testing.T) {
	tk := mustToken(t)
	code, _ := doJSON(t, "DELETE", "/todos/slug-tidak-ada", tk, nil)
	require.Equal(t, 404, code, "delete slug tak ada harus 404")
	code, _ = doJSON(t, "PATCH", "/todos/slug-tidak-ada/toggle", tk, nil)
	require.Equal(t, 404, code, "toggle slug tak ada harus 404")
}
