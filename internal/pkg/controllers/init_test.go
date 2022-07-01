package controllers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/require"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/sirupsen/logrus"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
	"github.com/stretchr/testify/assert"
)

const validLoginTokenWithout2Fa = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFwIiwic2NvcGVzIjpbeyJ1cmkiOiIqIiwibWV0aG9kIjoiKiJ9LHsidXJpIjoiL2FwaS92MS92ZXJpZnktMmZhIiwibWV0aG9kIjoiKiIsImV4Y2x1ZGUiOnRydWV9XSwianRpIjoiMTExODQ0MjU4NTQyNjAyNTgwODYifQ.0-OghMxBzf7CB2nTHkB0hOcteF1P9S20nsWY9gSttzE"                               //nolint:gosec
const validLoginTokenWith2FaWithoutTotPSecret = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFwIiwic2NvcGVzIjpbeyJ1cmkiOiIvYXBpL3YxL3ZlcmlmeS0yZmEiLCJtZXRob2QiOiJQT1NUIn0seyJ1cmkiOiIvYXBpL3YxL21lL3RvdHAtc2VjcmV0IiwibWV0aG9kIjoiUE9TVCJ9XSwianRpIjoiMTUwMDE4MjEzMTM3MDAyNTYxNTgifQ.cueDSueE2he35Y262NKBoUHkxJ4FpQcWZacIQABPf1s" //nolint:gosec
const validLoginTokenWith2FaWithTotPSecret = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFwIiwic2NvcGVzIjpbeyJ1cmkiOiIvYXBpL3YxL3ZlcmlmeS0yZmEiLCJtZXRob2QiOiJQT1NUIn1dLCJqdGkiOiIxNDI5ODAxMjQyOTUyNDk4MDkzOSJ9.whujLHERis49xoodPhuRUWrazSwJMLnNgIEZbE3kAe4"                                                                      //nolint:gosec

type TotPSecretRendererMock struct {
	mock.Mock
}

func (tsrm *TotPSecretRendererMock) RenderTotPSecret(key *models.TotPSecretOutput) error {
	args := tsrm.Called(key)

	return args.Error(0)
}

func TestInitSuccess(t *testing.T) {
	statusRequested := false

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
				Token: validLoginTokenWithout2Fa,
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
	tController := InitController{
		ConfigWriter: func(params *options.ParameterBag) (err error) {
			writtenParams = params
			return nil
		},
		PromptReader:       &PromptReaderMock{},
		TotPSecretRenderer: &TotPSecretRendererMock{},
	}

	params := config.FromValues(map[string]string{
		config.APIURL:      srv.URL,
		config.APIUser:     "login",
		config.APIPassword: "passwords",
	})
	err = tController.InitConfig(context.Background(), params)

	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, srv.URL, config.ReadAPIURLWithDefault(writtenParams, ""))
	assert.Equal(t, validLoginTokenWithout2Fa, writtenParams.ReadString(config.Token, ""))
	assert.True(t, statusRequested)
}

func TestInitSuccessWithLegacyConfig(t *testing.T) {
	statusRequested := false

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
				Token: validLoginTokenWithout2Fa,
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
	tController := InitController{
		ConfigWriter: func(params *options.ParameterBag) (err error) {
			writtenParams = params
			return nil
		},
		PromptReader:       &PromptReaderMock{},
		TotPSecretRenderer: &TotPSecretRendererMock{},
	}

	params := config.FromValues(map[string]string{
		config.ServerURL: srv.URL,
		config.Login:     "login",
		config.Password:  "passwords",
	})
	err = tController.InitConfig(context.Background(), params)

	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, srv.URL, config.ReadAPIURLWithDefault(writtenParams, ""))
	assert.Equal(t, validLoginTokenWithout2Fa, writtenParams.ReadString(config.Token, ""))
	assert.True(t, statusRequested)
}

func assertBasicLoginPass(t *testing.T, expectedLogin, expectedPass string, r *http.Request) {
	login, pass, err := utils.ExtractBasicAuthLoginAndPassFromRequest(r)
	require.NoError(t, err)
	assert.Equal(t, expectedLogin, login)
	assert.Equal(t, expectedPass, pass)
}

func getTwoFaLoginFromRequestBody(r *http.Request) (*api.TwoFaLogin, error) {
	res := &api.TwoFaLogin{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(res)

	return res, err
}

func assertTwoFaLoginFromRequest(t *testing.T, r *http.Request, expectedUsername, expectedToken string) {
	twoFALogin, err := getTwoFaLoginFromRequestBody(r)
	require.NoError(t, err)
	assert.Equal(t, expectedUsername, twoFALogin.Username)
	assert.Equal(t, expectedToken, twoFALogin.Token)
}

func TestInit2FASuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		curURL := r.URL.String()
		resp := api.LoginResponse{
			Data: models.Token{
				Token: validLoginTokenWith2FaWithTotPSecret,
				TwoFA: models.TwoFA{
					SentTo: "no@mail.me",
				},
			},
		}

		if strings.HasPrefix(curURL, "/api/v1/login") {
			assertBasicLoginPass(t, "log1", "pass1", r)

			rw.WriteHeader(http.StatusOK)
			jsonEnc := json.NewEncoder(rw)
			e := jsonEnc.Encode(resp)
			assert.NoError(t, e)
			return
		}

		if !strings.HasPrefix(curURL, "/api/v1/verify-2fa") {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		assertTwoFaLoginFromRequest(t, r, "log1", "someCode")
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "Bearer "+validLoginTokenWith2FaWithTotPSecret, r.Header.Get("Authorization"))

		resp.Data.Token = "someTok"
		resp.Data.TwoFA = models.TwoFA{}

		rw.WriteHeader(http.StatusOK)
		jsonEnc := json.NewEncoder(rw)
		e := jsonEnc.Encode(resp)
		assert.NoError(t, e)
	}))
	defer srv.Close()

	writtenParams := options.New(options.NewMapValuesProvider(map[string]interface{}{}))
	tController := InitController{
		ConfigWriter: func(params *options.ParameterBag) (err error) {
			writtenParams = params
			return nil
		},
		PromptReader: &PromptReaderMock{
			ReadOutputs: []string{
				"someCode",
			},
		},
		TotPSecretRenderer: &TotPSecretRendererMock{},
	}

	params := config.FromValues(map[string]string{
		config.APIURL:      srv.URL,
		config.APIUser:     "log1",
		config.APIPassword: "pass1",
	})
	err := tController.InitConfig(context.Background(), params)

	require.NoError(t, err)

	assert.Equal(t, srv.URL, config.ReadAPIURLWithDefault(writtenParams, ""))
	assert.Equal(t, "someTok", writtenParams.ReadString(config.Token, ""))
}

func TestInitTotPWithoutSecretSuccess(t *testing.T) {
	const qrCodeImageContent = "qrCodeContent"
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		curURL := r.URL.String()
		resp := api.LoginResponse{
			Data: models.Token{
				Token: validLoginTokenWith2FaWithoutTotPSecret,
				TwoFA: models.TwoFA{
					DeliveryMethod: "totp_authenticator_app",
					TotPKeyStatus:  api.TotPKeyPending,
				},
			},
		}

		if strings.HasPrefix(curURL, "/api/v1/login") {
			assertBasicLoginPass(t, "log1", "pass1", r)

			rw.WriteHeader(http.StatusOK)
			jsonEnc := json.NewEncoder(rw)
			e := jsonEnc.Encode(resp)
			assert.NoError(t, e)
			return
		}

		if strings.HasPrefix(curURL, "/api/v1/me/totp-secret") {
			assert.Equal(t, "Bearer "+validLoginTokenWith2FaWithoutTotPSecret, r.Header.Get("Authorization"))

			qrCodeContent := base64.StdEncoding.EncodeToString([]byte(qrCodeImageContent))
			totPResp := &models.TotPSecretResp{
				Secret:        "some secret",
				QRImageBase64: qrCodeContent,
			}
			rw.WriteHeader(http.StatusOK)
			jsonEnc := json.NewEncoder(rw)
			e := jsonEnc.Encode(totPResp)
			assert.NoError(t, e)
			return
		}

		if !strings.HasPrefix(curURL, "/api/v1/verify-2fa") {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		assertTwoFaLoginFromRequest(t, r, "log1", "123456")

		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "Bearer "+validLoginTokenWith2FaWithoutTotPSecret, r.Header.Get("Authorization"))

		resp.Data.Token = validLoginTokenWithout2Fa
		resp.Data.TwoFA = models.TwoFA{}

		rw.WriteHeader(http.StatusOK)
		jsonEnc := json.NewEncoder(rw)
		e := jsonEnc.Encode(resp)
		assert.NoError(t, e)
	}))
	defer srv.Close()

	totpSecretRenderer := &TotPSecretRendererMock{}

	totpSecretRenderer.On("RenderTotPSecret", mock.Anything).Return(nil)

	writtenParams := options.New(options.NewMapValuesProvider(map[string]interface{}{}))

	qrCodeBuf := &bytes.Buffer{}
	qrCodeFilePatternGiven := ""
	tController := InitController{
		ConfigWriter: func(params *options.ParameterBag) (err error) {
			writtenParams = params
			return nil
		},
		PromptReader: &PromptReaderMock{
			ReadOutputs: []string{
				"123456",
			},
		},
		TotPSecretRenderer: totpSecretRenderer,
		QrImageWriterProvider: func(namePattern string) (writer io.Writer, closr io.Closer, name string, err error) {
			qrCodeFilePatternGiven = namePattern
			return qrCodeBuf, nil, "qr-file.png", nil
		},
	}

	params := config.FromValues(map[string]string{
		config.APIURL:      srv.URL,
		config.APIUser:     "log1",
		config.APIPassword: "pass1",
	})
	err := tController.InitConfig(context.Background(), params)
	require.NoError(t, err)

	totpSecretRenderer.AssertCalled(
		t,
		"RenderTotPSecret",
		mock.MatchedBy(
			func(actualOutput *models.TotPSecretOutput) bool {
				return actualOutput.Secret == "some secret" &&
					actualOutput.File == "qr-file.png"
			},
		),
	)

	assert.Equal(t, qrCodeImageContent, qrCodeBuf.String())
	assert.Equal(t, "qr-*.png", qrCodeFilePatternGiven)
	assert.Equal(t, srv.URL, config.ReadAPIURLWithDefault(writtenParams, ""))
	assert.Equal(t, validLoginTokenWithout2Fa, writtenParams.ReadString(config.Token, ""))
}

func TestInitTotPWithSecretSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		curURL := r.URL.String()

		if strings.HasPrefix(curURL, "/api/v1/login") {
			rw.WriteHeader(http.StatusOK)
			jsonEnc := json.NewEncoder(rw)
			e := jsonEnc.Encode(api.LoginResponse{
				Data: models.Token{
					Token: validLoginTokenWith2FaWithTotPSecret,
					TwoFA: models.TwoFA{
						DeliveryMethod: "totp_authenticator_app",
						TotPKeyStatus:  api.TotPKeyExists,
					},
				},
			})
			assert.NoError(t, e)
			return
		}

		if !strings.HasPrefix(curURL, "/api/v1/verify-2fa") {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		assertTwoFaLoginFromRequest(t, r, "somelogin", "341234")

		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "Bearer "+validLoginTokenWith2FaWithTotPSecret, r.Header.Get("Authorization"))

		rw.WriteHeader(http.StatusOK)
		jsonEnc := json.NewEncoder(rw)
		e := jsonEnc.Encode(api.LoginResponse{
			Data: models.Token{
				Token: validLoginTokenWithout2Fa,
				TwoFA: models.TwoFA{},
			},
		})
		assert.NoError(t, e)
	}))
	defer srv.Close()

	writtenParams := options.New(options.NewMapValuesProvider(map[string]interface{}{}))

	tController := InitController{
		ConfigWriter: func(params *options.ParameterBag) (err error) {
			writtenParams = params
			return nil
		},
		PromptReader: &PromptReaderMock{
			ReadOutputs: []string{
				"341234",
			},
		},
		TotPSecretRenderer: &TotPSecretRendererMock{},
		QrImageWriterProvider: func(namePattern string) (writer io.Writer, closr io.Closer, name string, err error) {
			return nil, nil, "", errors.New("should not be called")
		},
	}

	params := config.FromValues(map[string]string{
		config.APIURL:      srv.URL,
		config.APIUser:     "somelogin",
		config.APIPassword: "somepass",
	})
	err := tController.InitConfig(context.Background(), params)
	require.NoError(t, err)

	assert.Equal(t, srv.URL, config.ReadAPIURLWithDefault(writtenParams, ""))
	assert.Equal(t, validLoginTokenWithout2Fa, writtenParams.ReadString(config.Token, ""))
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
		PromptReader:       &PromptReaderMock{},
		TotPSecretRenderer: &TotPSecretRendererMock{},
	}

	params := config.FromValues(map[string]string{
		config.APIURL:      srv.URL,
		config.APIUser:     "log1123",
		config.APIPassword: "pass111",
	})
	err := tController.InitConfig(context.Background(), params)

	assert.EqualError(t, err, "config verification failed: operation failed")
}

func TestInitErrorWithLegacyConfig(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	tController := InitController{
		ConfigWriter: func(params *options.ParameterBag) (err error) {
			return nil
		},
		PromptReader:       &PromptReaderMock{},
		TotPSecretRenderer: &TotPSecretRendererMock{},
	}

	params := config.FromValues(map[string]string{
		config.ServerURL: srv.URL,
		config.Login:     "log1123",
		config.Password:  "pass111",
	})
	err := tController.InitConfig(context.Background(), params)

	assert.EqualError(t, err, "config verification failed: operation failed")
}

func TestInitErrorWithApiTokenAndPassword(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	tController := InitController{
		ConfigWriter: func(params *options.ParameterBag) (err error) {
			return nil
		},
		PromptReader:       &PromptReaderMock{},
		TotPSecretRenderer: &TotPSecretRendererMock{},
	}

	params := config.FromValues(map[string]string{
		config.APIURL:   srv.URL,
		config.Login:    "log1123",
		config.Password: "pass111",
		config.APIToken: "123123",
	})
	err := tController.InitConfig(context.Background(), params)

	assert.EqualError(t, err, "config verification failed: RPORT_API_TOKEN and a password cannot be set at the same time. Please choose one and remove use of the other.")
}
