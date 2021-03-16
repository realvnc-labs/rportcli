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
			login = "root"
			pass = "root"
			return
		},
	}
	err = ba.AuthRequest(req)
	assert.NoError(t, err)

	assert.Equal(t, "Basic cm9vdDpyb290", req.Header.Get("Authorization"))
}
