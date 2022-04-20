package api

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilters(t *testing.T) {
	filters := NewFilters(
		"name", "johny",
		"id", "",
		"*", "*abc*",
	)

	q := url.Values{}
	filters.Apply(q)

	assert.Equal(t, q.Get("filter[name]"), "johny")
	assert.False(t, q.Has("id"))
	assert.Equal(t, q.Get("filter[*]"), "*abc*")
}
