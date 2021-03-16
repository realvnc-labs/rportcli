package utils

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthRequest(t *testing.T) {
	req, err := http.NewRequestWithContext(context.Background(), "post", "/", nil)
	assert.NoError(t, err)

	ba := &StorageBasicAuth{
		AuthProvider: func() (login, pass string, err error) {
			login = "looooaaaa"
			pass = "paoooo"
			return
		},
	}
	err = ba.AuthRequest(req)
	assert.NoError(t, err)

	assert.Equal(t, "Basic bG9vb29hYWFhOnBhb29vbw==", req.Header.Get("Authorization"))
}
