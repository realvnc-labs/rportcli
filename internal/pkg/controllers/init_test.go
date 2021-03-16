package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/sirupsen/logrus"

	options "github.com/breathbath/go_utils/utils/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestInitSuccess(t *testing.T) {
	statusRequested := false

	const tokenGiven = "some token"
	const tokenValidityVal = "90"
	err := os.Setenv(config.SessionValiditySecondsEnvVar, tokenValidityVal)
	assert.NoError(t, err)
	if err == nil {
		defer func() {
			e := os.Unsetenv(config.SessionValiditySecondsEnvVar)
			if e != nil {
				logrus.Error(e)
			}
		}()
	}

	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		resp := api.LoginResponse{
			Data: models.Token{
				Token: tokenGiven,
			},
		}
		statusRequested = true
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/v1/login?token-lifetime="+tokenValidityVal, r.URL.String())
		assert.Equal(t, "Basic bG9naW46cGFzc3dvcmRz", r.Header.Get("Authorization"))

		rw.WriteHeader(http.StatusOK)
		jsonEnc := json.NewEncoder(rw)
		e := jsonEnc.Encode(resp)
		assert.NoError(t, e)
	}))
	defer srv.Close()

	writtenParams := options.New(options.NewMapValuesProvider(map[string]interface{}{}))
	writtenParamsP := &writtenParams
	tController := InitController{
		ConfigWriter: func(params *options.ParameterBag) (err error) {
			writtenParamsP = params
			return nil
		},
		PromptReader: &PromptReaderMock{},
	}

	srvURL := srv.URL
	login := "login"
	pass := "passwords"

	err = tController.InitConfig(context.Background(), map[string]*string{
		"server_url": &srvURL,
		"login":      &login,
		"password":   &pass,
	})

	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, srvURL, writtenParamsP.ReadString(config.ServerURL, ""))
	assert.Equal(t, tokenGiven, writtenParamsP.ReadString(config.Token, ""))
	assert.True(t, statusRequested)
}

func TestInitFromPrompt(t *testing.T) {
	const login = "one"
	const pass = "two"
	const tokenToGive = "some tok"

	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		resp := api.LoginResponse{
			Data: models.Token{
				Token: tokenToGive,
			},
		}

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

	assert.Equal(t, srvURL, writtenParamsP.ReadString(config.ServerURL, ""))
	assert.Equal(t, tokenToGive, writtenParamsP.ReadString(config.Token, ""))
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
		PromptReader: &PromptReaderMock{},
	}

	srvURL := srv.URL
	login := "log1123"
	password := "pass111"

	err := tController.InitConfig(context.Background(), map[string]*string{
		config.ServerURL: &srvURL,
		config.Login:     &login,
		config.Password:  &password,
	})

	assert.EqualError(t, err, "config verification failed against the rport: operation failed")
}
