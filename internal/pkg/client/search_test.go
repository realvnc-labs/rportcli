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

func TestFindClientsFromDataProvider(t *testing.T) {
	search := Search{
		DataProvider: &DataProviderMock{
			clientsToGive: clientsList,
		},
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

func TestDataProviderError(t *testing.T) {
	search := Search{
		DataProvider: &DataProviderMock{
			clientsToGive: clientsList,
			errToGive:     errors.New("some load error"),
		},
	}

	_, err := search.Search(context.Background(), "$100", &options.ParameterBag{})
	assert.EqualError(t, err, "some load error")
}

func TestFindByAmbiguousClientName(t *testing.T) {
	search := Search{
		DataProvider: &DataProviderMock{
			clientsToGive: clientsList,
		},
	}

	_, err := search.FindOne(context.Background(), "my tiny", &options.ParameterBag{})
	assert.EqualError(t, err, `client identified by 'my tiny' is ambiguous, use a more precise name or use the client id`)
}
