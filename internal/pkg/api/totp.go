package api

import (
	"context"
	"net/http"

	"github.com/breathbath/go_utils/v2/pkg/url"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
)

const (
	TotPSecretURL  = "/api/v1/me/totp-secret" //nolint:gosec
	TotPKeyPending = "pending"
	TotPKeyExists  = "exists"
)

func (rp *Rport) CreateTotPSecret(ctx context.Context) (key *models.TotPSecretResp, err error) {
	var req *http.Request
	req, err = http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url.JoinURL(rp.BaseURL, TotPSecretURL),
		nil,
	)
	if err != nil {
		return
	}

	key = &models.TotPSecretResp{}
	_, err = rp.CallBaseClient(req, key)

	return key, err
}
