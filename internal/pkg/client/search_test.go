package client

import (
	"context"
	"errors"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

var clientsList = []models.Client{
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
	clientsToGive []models.Client
	errToGive     error
}

func (dpm *DataProviderMock) GetClients(ctx context.Context) (cls []models.Client, err error) {
	return dpm.clientsToGive, dpm.errToGive
}

type CacheMock struct {
	clientsToStore  []models.Client
	storeErrToGive  error
	existsErrToGive error
	clientsToLoad   []models.Client
	loadClientsErr  error
}

func (cm *CacheMock) Store(ctx context.Context, cls []models.Client) error {
	cm.clientsToStore = cls
	return cm.storeErrToGive
}

func (cm *CacheMock) Exists(ctx context.Context) (bool, error) {
	return len(cm.clientsToLoad) > 1, cm.existsErrToGive
}

func (cm *CacheMock) Load(ctx context.Context, cls *[]models.Client) error {
	*cls = append(*cls, cm.clientsToLoad...)
	return cm.loadClientsErr
}

func TestFindClientsFromDataProvider(t *testing.T) {
	cacheMock := &CacheMock{
		clientsToStore: []models.Client{},
		clientsToLoad:  []models.Client{},
	}
	search := Search{
		DataProvider: &DataProviderMock{
			clientsToGive: clientsList,
		},
		Cache: cacheMock,
	}

	foundCls, err := search.Search(context.Background(), "my tiny")
	assert.NoError(t, err)
	assert.Len(t, foundCls, 2)
	assert.Equal(t, foundCls, []models.Client{
		{
			ID:       "1",
			Name:     "my tiny client",
		},
		{
			ID:       "2",
			Name:     "my Tiny nice client",
		},
	})
	assert.Equal(t, clientsList, cacheMock.clientsToStore)
}

func TestFindClientsFromCache(t *testing.T) {
	cacheMock := &CacheMock{
		clientsToStore: []models.Client{},
		clientsToLoad:  clientsList,
	}
	search := Search{
		DataProvider: &DataProviderMock{
			clientsToGive: clientsList[0:1],
		},
		Cache: cacheMock,
	}

	foundCls, err := search.Search(context.Background(), "$100")
	assert.NoError(t, err)
	assert.Len(t, foundCls, 1)
	assert.Equal(t, foundCls, []models.Client{
		{
			ID:       "3",
			Name:     "$100 usd client",
		},
	})
	assert.Len(t, cacheMock.clientsToStore, 0)
}

func TestDataProviderError(t *testing.T) {
	search := Search{
		DataProvider: &DataProviderMock{
			clientsToGive: clientsList,
			errToGive: errors.New("some load error"),
		},
		Cache: &CacheMock{
			clientsToStore: []models.Client{},
			clientsToLoad:  []models.Client{},
		},
	}

	_, err := search.Search(context.Background(), "$100")
	assert.EqualError(t, err, "some load error")
}

func TestCacheStoreError(t *testing.T) {
	search := Search{
		DataProvider: &DataProviderMock{
			clientsToGive: clientsList,
		},
		Cache: &CacheMock{
			clientsToStore: []models.Client{},
			storeErrToGive: errors.New("some store err"),
			clientsToLoad:  []models.Client{},
		},
	}

	_, err := search.Search(context.Background(), "$100")
	assert.EqualError(t, err, "some store err")
}

func TestCacheExistsError(t *testing.T) {
	search := Search{
		DataProvider: &DataProviderMock{
			clientsToGive: []models.Client{},
		},
		Cache: &CacheMock{
			clientsToStore: []models.Client{},
			existsErrToGive: errors.New("some exists err"),
			clientsToLoad:  []models.Client{},
		},
	}

	_, err := search.Search(context.Background(), "$100")
	assert.EqualError(t, err, "some exists err")
}

func TestLoadCacheError(t *testing.T) {
	search := Search{
		DataProvider: &DataProviderMock{
			clientsToGive: []models.Client{},
		},
		Cache: &CacheMock{
			clientsToStore: []models.Client{},
			loadClientsErr: errors.New("some cache load err"),
			clientsToLoad:  clientsList,
		},
	}

	_, err := search.Search(context.Background(), "$100")
	assert.EqualError(t, err, "some cache load err")
}
