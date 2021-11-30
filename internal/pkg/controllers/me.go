package controllers

import (
	"context"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/output"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
)

type MeRenderer interface {
	RenderMe(os output.KvProvider) error
}

type MeAPI interface {
	Me(ctx context.Context) (user api.UserResponse, err error)
}

type MeController struct {
	Rport      MeAPI
	MeRenderer MeRenderer
}

func (tc *MeController) Me(ctx context.Context) error {
	userResp, err := tc.Rport.Me(ctx)
	if err != nil {
		return err
	}

	return tc.MeRenderer.RenderMe(&userResp.Data)
}
