package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestClientsList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		assert.Equal(t, "Basic bG9nMTpwYXNzMQ==", authHeader)

		assert.Equal(t, ClientsURL, r.URL.String())
		jsonEnc := json.NewEncoder(rw)
		e := jsonEnc.Encode(ClientsResponse{Data: clientsStub})
		assert.NoError(t, e)
	}))
	defer srv.Close()

	cl := New(srv.URL, &utils.BasicAuth{
		Login: "log1",
		Pass:  "pass1",
	})

	clientsResp, err := cl.Clients(context.Background())
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.Equal(t, clientsStub, clientsResp.Data)
}
