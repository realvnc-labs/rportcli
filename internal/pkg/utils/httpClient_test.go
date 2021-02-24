package utils

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAuthRequest (t *testing.T) {
	req, err := http.NewRequest("post", "/", nil)
	assert.NoError(t, err)

	ba := &BasicAuth{
		Login: "root",
		Pass:  "root",
	}
	err = ba.AuthRequest(req)
	assert.NoError(t, err)

	assert.Equal(t, "Basic cm9vdDpyb290", req.Header.Get("Authorization"))
}
