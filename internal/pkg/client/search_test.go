package client

import (
	"context"
	"errors"
	"testing"

	options "github.com/breathbath/go_utils/v2/pkg/config"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/stretchr/testify/assert"
)

var clientsList = []*models.Client{
	{
		ID:   "1",
		Name: "my tiny client",
	},
	{
		ID:   "2",
		Name: "my Tiny nice client",
	},
	{
		ID:   "3",
		Name: "$100 usd client",
	},
}

type DataProviderMock struct {
	clientsToGive []*models.Client
	errToGive     error
}

func (dpm *DataProviderMock) GetClients(ctx context.Context) (cls []*models.Client, err error) {
	return dpm.clientsToGive, dpm.errToGive
}

type CacheMock struct {
	clientsToStore  []*models.Client
	storeErrToGive  error
	existsErrToGive error
	clientsToLoad   []*models.Client
	loadClientsErr  error
}

func (cm *CacheMock) Store(ctx context.Context, cls []*models.Client, params *options.ParameterBag) error {
	cm.clientsToStore = cls
	return cm.storeErrToGive
}

func (cm *CacheMock) Exists(ctx context.Context, params *options.ParameterBag) (bool, error) {
	return len(cm.clientsToLoad) > 1, cm.existsErrToGive
}

func (cm *CacheMock) Load(ctx context.Context, params *options.ParameterBag) (cls []*models.Client, err error) {
	return cm.clientsToLoad, cm.loadClientsErr
}

func TestFindClientsFromDataProvider(t *testing.T) {
	cacheMock := &CacheMock{
		clientsToStore: []*models.Client{},
		clientsToLoad:  []*models.Client{},
	}
	search := Search{
		DataProvider: &DataProviderMock{
			clientsToGive: clientsList,
		},
		Cache: cacheMock,
	}

	foundCls, err := search.Search(context.Background(), "my tiny", &options.ParameterBag{})
	assert.NoError(t, err)
	assert.Len(t, foundCls, 2)
	assert.Equal(t, foundCls, []*models.Client{
		{
			ID:   "1",
			Name: "my tiny client",
		},
		{
			ID:   "2",
			Name: "my Tiny nice client",
		},
	})
	assert.Equal(t, clientsList, cacheMock.clientsToStore)

	foundCls2, err2 := search.Search(context.Background(), "my tiny client,$100 usd client", &options.ParameterBag{})
	assert.NoError(t, err2)
	assert.Equal(t, foundCls2, []*models.Client{
		{
			ID:   "1",
			Name: "my tiny client",
		},
		{
			ID:   "3",
			Name: "$100 usd client",
		},
	})
}

func TestFindClientsFromCache(t *testing.T) {
	cacheMock := &CacheMock{
		clientsToStore: []*models.Client{},
		clientsToLoad:  clientsList,
	}
	search := Search{
		DataProvider: &DataProviderMock{
			clientsToGive: clientsList[0:1],
		},
		Cache: cacheMock,
	}

	foundCls, err := search.Search(context.Background(), "$100", &options.ParameterBag{})
	assert.NoError(t, err)
	assert.Len(t, foundCls, 1)
	assert.Equal(t, foundCls, []*models.Client{
		{
			ID:   "3",
			Name: "$100 usd client",
		},
	})
	assert.Len(t, cacheMock.clientsToStore, 0)
}

func TestDataProviderError(t *testing.T) {
	search := Search{
		DataProvider: &DataProviderMock{
			clientsToGive: clientsList,
			errToGive:     errors.New("some load error"),
		},
		Cache: &CacheMock{
			clientsToStore: []*models.Client{},
			clientsToLoad:  []*models.Client{},
		},
	}

	_, err := search.Search(context.Background(), "$100", &options.ParameterBag{})
	assert.EqualError(t, err, "some load error")
}

func TestCacheStoreError(t *testing.T) {
	search := Search{
		DataProvider: &DataProviderMock{
			clientsToGive: clientsList,
		},
		Cache: &CacheMock{
			clientsToStore: []*models.Client{},
			storeErrToGive: errors.New("some store err"),
			clientsToLoad:  []*models.Client{},
		},
	}

	_, err := search.Search(context.Background(), "$100", &options.ParameterBag{})
	assert.EqualError(t, err, "some store err")
}

func TestCacheExistsError(t *testing.T) {
	search := Search{
		DataProvider: &DataProviderMock{
			clientsToGive: []*models.Client{},
		},
		Cache: &CacheMock{
			clientsToStore:  []*models.Client{},
			existsErrToGive: errors.New("some exists err"),
			clientsToLoad:   []*models.Client{},
		},
	}

	_, err := search.Search(context.Background(), "$100", &options.ParameterBag{})
	assert.EqualError(t, err, "some exists err")
}

func TestLoadCacheError(t *testing.T) {
	search := Search{
		DataProvider: &DataProviderMock{
			clientsToGive: []*models.Client{},
		},
		Cache: &CacheMock{
			clientsToStore: []*models.Client{},
			loadClientsErr: errors.New("some cache load err"),
			clientsToLoad:  clientsList,
		},
	}

	_, err := search.Search(context.Background(), "$100", &options.ParameterBag{})
	assert.EqualError(t, err, "some cache load err")
}

func TestFindByAmbiguousClientName(t *testing.T) {
	cacheMock := &CacheMock{
		clientsToStore: []*models.Client{},
		clientsToLoad:  []*models.Client{},
	}
	search := Search{
		DataProvider: &DataProviderMock{
			clientsToGive: clientsList,
		},
		Cache: cacheMock,
	}

	_, err := search.FindOne(context.Background(), "my tiny", &options.ParameterBag{})
	assert.EqualError(t, err, `client identified by 'my tiny' is ambiguous, use a more precise name or use the client id`)
}
