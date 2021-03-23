package client

import (
	"context"
	"strings"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
)

type DataProvider interface {
	GetClients(ctx context.Context) (cls []models.Client, err error)
}

type Cache interface {
	Store(ctx context.Context, cls []models.Client) error
	Exists(ctx context.Context) (bool, error)
	Load(ctx context.Context, cls *[]models.Client) error
}

type Search struct {
	DataProvider DataProvider
	Cache        Cache
}

func (s *Search) Search(ctx context.Context, term string) (foundCls []models.Client, err error) {
	cls, err := s.getClientsList(ctx)
	if err != nil {
		return foundCls, err
	}

	foundCls = s.findInClientsList(cls, term)
	return
}

func (s *Search) getClientsList(ctx context.Context) (cls []models.Client, err error) {
	cacheExists, err := s.Cache.Exists(ctx)
	if err != nil {
		return cls, err
	}

	if !cacheExists {
		cls, err = s.DataProvider.GetClients(ctx)
		if err != nil {
			return
		}

		err = s.Cache.Store(ctx, cls)
		return
	}

	cls = []models.Client{}
	err = s.Cache.Load(ctx, &cls)
	return
}

func (s *Search) findInClientsList(cls []models.Client, term string) (foundCls []models.Client) {
	foundCls = make([]models.Client, 0)
	for i := range cls {
		cl := cls[i]
		curClientName := strings.ToLower(cl.Name)
		term = strings.ToLower(term)
		if strings.HasPrefix(curClientName, term) {
			foundCls = append(foundCls, cl)
		}
	}

	return
}
