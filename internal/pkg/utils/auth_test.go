package utils

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type AuthMock struct {
	req       []*http.Request
	errToGive error
}

func (am *AuthMock) AuthRequest(req *http.Request) error {
	am.req = append(am.req, req)
	return am.errToGive
}

func TestAuthRequestBasic(t *testing.T) {
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

func TestAuthRequestBearer(t *testing.T) {
	req, err := http.NewRequestWithContext(context.Background(), "post", "/", nil)
	assert.NoError(t, err)

	ba := &BearerAuth{
		TokenProvider: func() (token string, err error) {
			token = "2§1234132423"
			return
		},
	}
	err = ba.AuthRequest(req)
	assert.NoError(t, err)

	assert.Equal(t, "Bearer 2§1234132423", req.Header.Get("Authorization"))
}

func TestFallbackAuth(t *testing.T) {
	primaryAuth := &AuthMock{
		req:       []*http.Request{},
		errToGive: errors.New("some error"),
	}
	fallbackAuth := &AuthMock{
		req: []*http.Request{},
	}

	fa := &FallbackAuth{
		PrimaryAuth:  primaryAuth,
		FallbackAuth: fallbackAuth,
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/some", nil)
	assert.NoError(t, err)

	err = fa.AuthRequest(req)
	assert.NoError(t, err)

	assert.Len(t, primaryAuth.req, 1)
	assert.Len(t, fallbackAuth.req, 1)
	assert.Equal(t, "/some", fallbackAuth.req[0].URL.String())
	assert.Equal(t, http.MethodPost, fallbackAuth.req[0].Method)
}

func TestFallbackAuthError(t *testing.T) {
	primaryAuth := &AuthMock{
		req:       []*http.Request{},
		errToGive: errors.New("some primaryAuth error"),
	}
	fallbackAuth := &AuthMock{
		req:       []*http.Request{},
		errToGive: errors.New("some fallbackAuth error"),
	}

	fa := &FallbackAuth{
		PrimaryAuth:  primaryAuth,
		FallbackAuth: fallbackAuth,
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/some", nil)
	assert.NoError(t, err)

	err = fa.AuthRequest(req)
	assert.EqualError(t, err, "some fallbackAuth error")
}

func TestFallbackAuthPrimarySuccess(t *testing.T) {
	primaryAuth := &AuthMock{
		req: []*http.Request{},
	}
	fallbackAuth := &AuthMock{
		req:       []*http.Request{},
		errToGive: errors.New("some error"),
	}

	fa := &FallbackAuth{
		PrimaryAuth:  primaryAuth,
		FallbackAuth: fallbackAuth,
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/other", nil)
	assert.NoError(t, err)

	err = fa.AuthRequest(req)
	assert.NoError(t, err)

	assert.Len(t, primaryAuth.req, 1)
	assert.Len(t, fallbackAuth.req, 0)
	assert.Equal(t, "/other", primaryAuth.req[0].URL.String())
	assert.Equal(t, http.MethodPost, primaryAuth.req[0].Method)
}
