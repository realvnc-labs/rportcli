package controllers

import (
	"context"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
)

type ClientSearch interface {
	Search(ctx context.Context, term string) (foundCls []models.Client, err error)
}
