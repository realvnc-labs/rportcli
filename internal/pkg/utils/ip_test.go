package utils

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIPFetching(t *testing.T) {
	const ipToCheck = "123.123.11.11"
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		_, err := rw.Write([]byte(ipToCheck))
		assert.NoError(t, err)
	}))

	defer srv.Close()

	ipProvider := APIIPProvider{
		URL: srv.URL,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	actualIP, err := ipProvider.GetIP(ctx)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.Equal(t, ipToCheck, actualIP)
}
