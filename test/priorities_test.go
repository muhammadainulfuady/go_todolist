package api_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPriorities(t *testing.T) {
	code, resp := doJSON(t, "GET", "/priorities", "", nil)
	require.Equal(t, 200, code)
	var items []struct {
		ID   int    `json:"id_priorities"`
		Name string `json:"name"`
	}
	parseData(t, resp.Data, &items)
	require.Len(t, items, 4, "priorities harus 4")
}
