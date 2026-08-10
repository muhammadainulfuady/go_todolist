package api_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHealth(t *testing.T) {
	code, resp := doJSON(t, "GET", "/health", "", nil)
	require.Equal(t, 200, code)
	require.Equal(t, "success", resp.Status)
}
