package api

import (
	"net/url"
	"testing"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestPaginationFromParams(t *testing.T) {
	pagination := NewPaginationFromParams(options.New(options.NewMapValuesProvider(map[string]interface{}{
		"limit":  30,
		"offset": 90,
	})))

	q := url.Values{}
	pagination.Apply(q)

	assert.Equal(t, q.Get("page[limit]"), "30")
	assert.Equal(t, q.Get("page[offset]"), "90")
}
