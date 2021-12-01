package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTotPSecret(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		assert.Equal(t, "Basic bG9nMTMzMzpwYXNzMTM0MTIzNA==", authHeader)

		assert.Equal(t, TotPSecretURL, r.URL.String())
		jsonEnc := json.NewEncoder(rw)
		e := jsonEnc.Encode(models.TotPSecretResp{
			Secret:        "123",
			QRImageBase64: "base52",
		})
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

	createdSecretResp, err := cl.CreateTotPSecret(context.Background())
	require.NoError(t, err)

	assert.Equal(t, "123", createdSecretResp.Secret)
	assert.Equal(t, "base52", createdSecretResp.QRImageBase64)
}
