package api

import (
	"net/http"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"
	"github.com/sirupsen/logrus"
)

type Auth interface {
	AuthRequest(r *http.Request) error
}

type Rport struct {
	BaseURL string
	Auth    Auth
}

func New(baseURL string, a Auth) *Rport {
	return &Rport{BaseURL: baseURL, Auth: a}
}

func (rp *Rport) CallBaseClient(req *http.Request, target interface{}) (resp *http.Response, err error) {
	cl := &utils.BaseClient{}
	cl.WithAuth(rp.Auth)

	var errResp models.ErrorResp

	resp, err = cl.Call(req, target, &errResp)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusUnauthorized {
			logrus.Warn("The request was unauthorized. The server might have 2FA enabled which prevents authentication " +
				"by environment variable RPORT_API_PASSWORD. Try with RPORT_API_TOKEN instead. RPORT_API_USER must also be " +
				"set when using RPORT_API_TOKEN.")
		}
		return nil, err
	}

	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			logrus.Error(closeErr)
		}
	}()

	return
}
