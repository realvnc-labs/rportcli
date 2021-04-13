package client

import (
	"context"
	"strings"

	options "github.com/breathbath/go_utils/v2/pkg/config"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
)

type DataProvider interface {
	GetClients(ctx context.Context) (cls []models.Client, err error)
}

type Cache interface {
	Store(ctx context.Context, cls []models.Client, params *options.ParameterBag) error
	Exists(ctx context.Context, params *options.ParameterBag) (bool, error)
	Load(ctx context.Context, cls *[]models.Client, params *options.ParameterBag) error
}

type Search struct {
	DataProvider DataProvider
	Cache        Cache
}

func (s *Search) Search(ctx context.Context, term string, params *options.ParameterBag) (foundCls []models.Client, err error) {
	cls, err := s.getClientsList(ctx, params)
	if err != nil {
		return foundCls, err
	}

	foundCls = s.findInClientsList(cls, term)
	return
}

func (s *Search) getClientsList(ctx context.Context, params *options.ParameterBag) (cls []models.Client, err error) {
	cacheExists, err := s.Cache.Exists(ctx, params)
	if err != nil {
		return cls, err
	}

	if !cacheExists {
		cls, err = s.DataProvider.GetClients(ctx)
		if err != nil {
			return
		}

		err = s.Cache.Store(ctx, cls, params)
		return
	}

	cls = []models.Client{}
	err = s.Cache.Load(ctx, &cls, params)
	return
}

func (s *Search) findInClientsList(cls []models.Client, term string) (foundCls []models.Client) {
	terms := strings.Split(term, ",")
	for i := range terms {
		terms[i] = strings.ToLower(terms[i])
	}

	foundCls = make([]models.Client, 0)
	for i := range cls {
		cl := cls[i]
		curClientName := strings.ToLower(cl.Name)

		for i := range terms {
			curTerm := terms[i]
			if strings.HasPrefix(curClientName, curTerm) {
				foundCls = append(foundCls, cl)
			}
		}
	}

	return
}
