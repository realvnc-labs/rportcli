package controllers

import (
	"context"

	options "github.com/breathbath/go_utils/v2/pkg/config"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
)

type ClientSearch interface {
	Search(ctx context.Context, term string, params *options.ParameterBag) (foundCls []models.Client, err error)
}
