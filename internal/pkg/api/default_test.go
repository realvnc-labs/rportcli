package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		assert.Equal(t, "Basic bG9nMTpwYXNzMQ==", authHeader)

		assert.Equal(t, "/api/v1/login?token-lifetime=10", r.URL.String())
		jsonEnc := json.NewEncoder(rw)
		e := jsonEnc.Encode(LoginResponse{Data: models.Token{
			Token: "token123",
		}})
		assert.NoError(t, e)
	}))
	defer srv.Close()

	cl := New(srv.URL, &utils.StorageBasicAuth{
		AuthProvider: func() (login, pass string, err error) {
			login = "log1"
			pass = "pass1"
			return
		},
	})

	loginInfo, err := cl.GetToken(context.Background(), 10)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.Equal(t, "token123", loginInfo.Data.Token)
}

func TestMe(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		assert.Equal(t, "Basic bG9nMTpwYXNzMQ==", authHeader)

		assert.Equal(t, "/api/v1/me", r.URL.String())
		jsonEnc := json.NewEncoder(rw)
		e := jsonEnc.Encode(UserResponse{Data: models.User{
			User:   "someUser",
			Groups: []string{"group1", "group2"},
		}})
		assert.NoError(t, e)
	}))
	defer srv.Close()

	cl := New(srv.URL, &utils.StorageBasicAuth{
		AuthProvider: func() (login, pass string, err error) {
			login = "log1"
			pass = "pass1"
			return
		},
	})

	usrResp, err := cl.Me(context.Background())
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.Equal(t, "someUser", usrResp.Data.User)
	assert.Equal(t, []string{"group1", "group2"}, usrResp.Data.Groups)
}

func TestStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		assert.Equal(t, "Basic bG9nMTpwYXNzMQ==", authHeader)

		assert.Equal(t, "/api/v1/status", r.URL.String())
		jsonEnc := json.NewEncoder(rw)
		e := jsonEnc.Encode(StatusResponse{Data: models.Status{
			ClientsConnected:    3,
			ClientsDisconnected: 1,
			Version:             "v123",
			Fingerprint:         "fp123",
			ConnectURL:          "conn",
		}})
		assert.NoError(t, e)
	}))
	defer srv.Close()

	cl := New(srv.URL, &utils.StorageBasicAuth{
		AuthProvider: func() (login, pass string, err error) {
			login = "log1"
			pass = "pass1"
			return
		},
	})

	statusResp, err := cl.Status(context.Background())
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.Equal(t, "v123", statusResp.Data.Version)
	assert.Equal(t, "conn", statusResp.Data.ConnectURL)
	assert.Equal(t, "fp123", statusResp.Data.Fingerprint)
	assert.Equal(t, 3, statusResp.Data.ClientsConnected)
	assert.Equal(t, 1, statusResp.Data.ClientsDisconnected)
}

func TestErrorResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusBadRequest)
		jsonEnc := json.NewEncoder(rw)
		e := jsonEnc.Encode(models.ErrorResp{
			Errors: []models.Error{
				{
					Code:   "400",
					Title:  "some title",
					Detail: "unauthorized",
				},
			},
		})
		assert.NoError(t, e)
	}))
	defer srv.Close()

	cl := New(srv.URL, &utils.StorageBasicAuth{
		AuthProvider: func() (login, pass string, err error) {
			login = "log1"
			pass = "pass1"
			return
		},
	})

	_, err := cl.Status(context.Background())
	assert.Error(t, err)
	if err == nil {
		return
	}
	errResp, ok := err.(*models.ErrorResp)
	assert.True(t, ok)
	if !ok {
		return
	}

	expectedErrors := []models.Error{
		{
			Code:   "400",
			Title:  "some title",
			Detail: "unauthorized",
		},
	}
	assert.Equal(t, expectedErrors, errResp.Errors)
}
