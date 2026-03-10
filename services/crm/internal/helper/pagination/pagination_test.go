package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertPaginationToStrSql_WithAllParameters(t *testing.T) {
	pag := &PaginationDTI{
		Page:     "2",
		PerPage:  "10",
		SortBy:   "created_at",
		SortDesc: "true",
	}

	result, err := ConvertPaginationToStrSql(pag)

	assert.NoError(t, err)
	assert.Contains(t, result, "ORDER BY created_at")
	assert.Contains(t, result, "DESC")
	assert.Contains(t, result, "LIMIT 10 OFFSET 10")
}

func TestConvertPaginationToStrSql_WithAscendingSort(t *testing.T) {
	pag := &PaginationDTI{
		Page:     "1",
		PerPage:  "20",
		SortBy:   "name",
		SortDesc: "false",
	}

	result, err := ConvertPaginationToStrSql(pag)

	assert.NoError(t, err)
	assert.Contains(t, result, "ORDER BY name")
	assert.NotContains(t, result, "DESC")
	assert.Contains(t, result, "LIMIT 20 OFFSET 0")
}

func TestConvertPaginationToStrSql_WithoutSort(t *testing.T) {
	pag := &PaginationDTI{
		Page:     "1",
		PerPage:  "15",
		SortBy:   "",
		SortDesc: "",
	}

	result, err := ConvertPaginationToStrSql(pag)

	assert.NoError(t, err)
	assert.NotContains(t, result, "ORDER BY")
	assert.Contains(t, result, "LIMIT 15 OFFSET 0")
}

func TestConvertPaginationToStrSql_PageThree(t *testing.T) {
	pag := &PaginationDTI{
		Page:     "3",
		PerPage:  "10",
		SortBy:   "id",
		SortDesc: "false",
	}

	result, err := ConvertPaginationToStrSql(pag)

	assert.NoError(t, err)
	// Page 3 with 10 per page should have offset of 20
	assert.Contains(t, result, "LIMIT 10 OFFSET 20")
}

func TestConvertPaginationToStrSql_InvalidPage(t *testing.T) {
	pag := &PaginationDTI{
		Page:     "invalid",
		PerPage:  "10",
		SortBy:   "id",
		SortDesc: "false",
	}

	result, err := ConvertPaginationToStrSql(pag)

	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestConvertPaginationToStrSql_InvalidPerPage(t *testing.T) {
	pag := &PaginationDTI{
		Page:     "1",
		PerPage:  "invalid",
		SortBy:   "id",
		SortDesc: "false",
	}

	result, err := ConvertPaginationToStrSql(pag)

	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestConvertPaginationToStrSql_EmptyPagination(t *testing.T) {
	pag := &PaginationDTI{
		Page:     "",
		PerPage:  "",
		SortBy:   "id",
		SortDesc: "true",
	}

	result, err := ConvertPaginationToStrSql(pag)

	assert.NoError(t, err)
	assert.Contains(t, result, "ORDER BY id DESC")
	assert.NotContains(t, result, "LIMIT")
}

func TestGetResponsePagination_Success(t *testing.T) {
	pagDTI := &PaginationDTI{
		Page:    "2",
		PerPage: "10",
	}
	total := 25

	result, err := GetResponsePagination(pagDTI, total)

	assert.NoError(t, err)
	assert.Equal(t, 25, result.Total)
	assert.Equal(t, 2, result.Page)
	assert.Equal(t, 10, result.Limit)
	assert.True(t, result.HasMore) // 25 total > (2 * 10) = 20
}

func TestGetResponsePagination_NoMorePages(t *testing.T) {
	pagDTI := &PaginationDTI{
		Page:    "3",
		PerPage: "10",
	}
	total := 25

	result, err := GetResponsePagination(pagDTI, total)

	assert.NoError(t, err)
	assert.Equal(t, 25, result.Total)
	assert.Equal(t, 3, result.Page)
	assert.Equal(t, 10, result.Limit)
	assert.False(t, result.HasMore) // 25 total <= (3 * 10) = 30
}

func TestGetResponsePagination_ExactlyOnePage(t *testing.T) {
	pagDTI := &PaginationDTI{
		Page:    "1",
		PerPage: "20",
	}
	total := 20

	result, err := GetResponsePagination(pagDTI, total)

	assert.NoError(t, err)
	assert.Equal(t, 20, result.Total)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 20, result.Limit)
	assert.False(t, result.HasMore)
}

func TestGetResponsePagination_EmptyResult(t *testing.T) {
	pagDTI := &PaginationDTI{
		Page:    "1",
		PerPage: "10",
	}
	total := 0

	result, err := GetResponsePagination(pagDTI, total)

	assert.NoError(t, err)
	assert.Equal(t, 0, result.Total)
	assert.False(t, result.HasMore)
}

func TestGetResponsePagination_InvalidPage(t *testing.T) {
	pagDTI := &PaginationDTI{
		Page:    "invalid",
		PerPage: "10",
	}
	total := 50

	result, err := GetResponsePagination(pagDTI, total)

	assert.Error(t, err)
	assert.Equal(t, 50, result.Total)
}

func TestGetResponsePagination_InvalidPerPage(t *testing.T) {
	pagDTI := &PaginationDTI{
		Page:    "1",
		PerPage: "invalid",
	}
	total := 50

	result, err := GetResponsePagination(pagDTI, total)

	assert.Error(t, err)
	assert.Equal(t, 50, result.Total)
}

func TestGetResponsePagination_EmptyPagination(t *testing.T) {
	pagDTI := &PaginationDTI{
		Page:    "",
		PerPage: "",
	}
	total := 100

	result, err := GetResponsePagination(pagDTI, total)

	assert.NoError(t, err)
	assert.Equal(t, 100, result.Total)
	assert.Equal(t, 0, result.Page)
	assert.Equal(t, 0, result.Limit)
}

func TestIsHasMore_MorePagesAvailable(t *testing.T) {
	result := isHasMore(1, 10, 25)
	assert.True(t, result)

	result = isHasMore(2, 10, 25)
	assert.True(t, result)
}

func TestIsHasMore_NoMorePages(t *testing.T) {
	result := isHasMore(3, 10, 25)
	assert.False(t, result)

	result = isHasMore(5, 10, 25)
	assert.False(t, result)
}

func TestIsHasMore_ExactlyOnePage(t *testing.T) {
	result := isHasMore(1, 20, 20)
	assert.False(t, result)
}

func TestIsHasMore_EmptyResult(t *testing.T) {
	result := isHasMore(1, 10, 0)
	assert.False(t, result)
}
