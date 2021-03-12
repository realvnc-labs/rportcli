package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	options "github.com/breathbath/go_utils/utils/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestInitSuccess(t *testing.T) {
	statusRequested := false
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		resp := api.StatusResponse{}
		statusRequested = true
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/v1/status", r.URL.String())
		assert.Equal(t, "Basic bG9naW46cGFzc3dvcmRz", r.Header.Get("Authorization"))

		rw.WriteHeader(http.StatusOK)
		jsonEnc := json.NewEncoder(rw)
		err := jsonEnc.Encode(resp)
		assert.NoError(t, err)
	}))
	defer srv.Close()

	writtenParams := options.New(options.NewMapValuesProvider(map[string]interface{}{}))
	writtenParamsP := &writtenParams
	tController := InitController{
		ConfigWriter: func(params *options.ParameterBag) (err error) {
			writtenParamsP = params
			return nil
		},
	}

	srvURL := srv.URL
	login := "login"
	pass := "passwords"

	err := tController.InitConfig(context.Background(), map[string]*string{
		"server_url": &srvURL,
		"login":      &login,
		"password":   &pass,
	})

	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, srvURL, writtenParamsP.ReadString("server_url", ""))
	assert.Equal(t, login, writtenParamsP.ReadString("login", ""))
	assert.Equal(t, pass, writtenParamsP.ReadString("password", ""))
	assert.True(t, statusRequested)
}

func TestInitFromPrompt(t *testing.T) {
	login := "one"
	pass := "two"

	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		resp := api.StatusResponse{}
		rw.WriteHeader(http.StatusOK)
		jsonEnc := json.NewEncoder(rw)
		err := jsonEnc.Encode(resp)
		assert.NoError(t, err)
	}))
	defer srv.Close()
	srvURL := srv.URL

	writtenParams := options.New(options.NewMapValuesProvider(map[string]interface{}{}))
	writtenParamsP := &writtenParams
	tController := InitController{
		ConfigWriter: func(params *options.ParameterBag) (err error) {
			writtenParamsP = params
			return nil
		},
		PromptReader: &PromptReaderMock{
			ReadOutputs: []string{
				srvURL,
				login,
			},
			PasswordReadOutputs: []string{
				pass,
			},
		},
	}

	err := tController.InitConfig(context.Background(), map[string]*string{})

	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.Equal(t, srvURL, writtenParamsP.ReadString("server_url", ""))
	assert.Equal(t, login, writtenParamsP.ReadString("login", ""))
	assert.Equal(t, pass, writtenParamsP.ReadString("password", ""))
}

func TestInitError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	tController := InitController{
		ConfigWriter: func(params *options.ParameterBag) (err error) {
			return nil
		},
	}

	srvURL := srv.URL
	login := "log1"
	pass := "pass1"

	err := tController.InitConfig(context.Background(), map[string]*string{
		"server_url": &srvURL,
		"login":      &login,
		"password":   &pass,
	})

	assert.EqualError(t, err, "config verification failed against the rport: operation failed")
}
