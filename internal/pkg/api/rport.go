package api

import (
	"net/http"

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

	var errResp ErrorResp

	resp, err = cl.Call(req, target, &errResp)
	if err != nil {
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
