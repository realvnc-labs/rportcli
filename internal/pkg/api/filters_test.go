package api

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestFilters(t *testing.T) {
	filters, err := NewFilters(
		"name", "johny",
		"id", "",
		"*", "*abc*",
	)
	require.NoError(t, err)

	q := url.Values{}
	filters.Apply(q)

	assert.Equal(t, q.Get("filter[name]"), "johny")
	assert.False(t, q.Has("id"))
	assert.Equal(t, q.Get("filter[*]"), "*abc*")
}

func TestCombinedFilter(t *testing.T) {
	filters, err := NewFilterFromCombinedSearchString("name=ok&os_kernel=linux")
	require.NoError(t, err)

	q := url.Values{}
	filters.Apply(q)

	assert.Equal(t, "filter%5Bname%5D=ok&filter%5Bos_kernel%5D=linux", q.Encode())
}
