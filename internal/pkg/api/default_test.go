package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		assert.Equal(t, "Basic bG9nMTQzMTUzOnBhc3MxMTM0MzE=", authHeader)

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
			login = "log143153"
			pass = "pass113431"
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
		assert.Equal(t, "Basic bG9nMTEzNDM6cGFzczEzNDEyMzQzMg==", authHeader)

		assert.Equal(t, "/api/v1/me", r.URL.String())
		jsonEnc := json.NewEncoder(rw)
		e := jsonEnc.Encode(UserResponse{Data: models.Me{
			Username: "someUser",
			Groups:   []string{"group1", "group2"},
		}})
		assert.NoError(t, e)
	}))
	defer srv.Close()

	cl := New(srv.URL, &utils.StorageBasicAuth{
		AuthProvider: func() (login, pass string, err error) {
			login = "log11343"
			pass = "pass134123432"
			return
		},
	})

	usrResp, err := cl.Me(context.Background())
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.Equal(t, "someUser", usrResp.Data.Username)
	assert.Equal(t, []string{"group1", "group2"}, usrResp.Data.Groups)
}

func TestStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		assert.Equal(t, "Basic bG9nMTMzMzpwYXNzMTM0MTIzNA==", authHeader)

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
			login = "log1333"
			pass = "pass1341234"
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
			login = "log1dfasf"
			pass = "pass134124"
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

func TestLogout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		actualToken := r.Header.Get("Authorization")
		assert.Equal(t, "Bearer some_tok", actualToken)
		rw.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	cl := New(srv.URL, &utils.BearerAuth{
		TokenProvider: func() (string, error) {
			return "some_tok", nil
		},
	})
	err := cl.Logout(context.Background())
	require.NoError(t, err)

	cl2 := New(srv.URL, &utils.BearerAuth{
		TokenProvider: func() (string, error) {
			return "", errors.New("some failed to get token")
		},
	})
	err = cl2.Logout(context.Background())
	require.EqualError(t, err, "some failed to get token")
}

func TestLogoutServerError(t *testing.T) {
	testCases := []struct {
		respCodeToGive int
		errRespToGive  string
		errToExpect    string
		name           string
	}{
		{
			respCodeToGive: http.StatusInternalServerError,
			errToExpect:    "operation failed",
			name:           "handle_500code",
		},
		{
			respCodeToGive: http.StatusCreated,
			errToExpect:    "unexpected response code 201, 204 is expected",
			name:           "handle_201code",
		},
		{
			respCodeToGive: http.StatusBadRequest,
			errToExpect:    "some problem, code: 400, details: some problem",
			errRespToGive:  "some problem",
			name:           "handle_400code+error_resp",
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(testCases[i].name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				rw.WriteHeader(tc.respCodeToGive)

				if tc.errRespToGive != "" {
					jsonEnc := json.NewEncoder(rw)
					e := jsonEnc.Encode(models.ErrorResp{
						Errors: []models.Error{
							{
								Code:   "400",
								Title:  tc.errRespToGive,
								Detail: tc.errRespToGive,
							},
						},
					})
					assert.NoError(t, e)
				}
			}))
			defer srv.Close()

			cl := New(srv.URL, &utils.BearerAuth{
				TokenProvider: func() (string, error) {
					return "123", nil
				},
			})
			err := cl.Logout(context.Background())
			require.EqualError(t, err, tc.errToExpect)
		})
	}
}
