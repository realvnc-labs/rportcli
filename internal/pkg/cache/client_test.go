package cache

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	options "github.com/breathbath/go_utils/v2/pkg/config"

	"github.com/breathbath/go_utils/v2/pkg/io"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/stretchr/testify/assert"
)

// we put all in one method as in different tests it might give file access conflicts
func TestCache(t *testing.T) {
	assertExists(t)
	assertStoreAndLoad(t)
}

func assertStoreAndLoad(t *testing.T) {
	cc := &ClientsCache{}
	providedClients := []*models.Client{
		{
			ID:   "1",
			Name: "client 1",
		},
		{
			ID:   "2",
			Name: "client 2",
		},
	}
	err := cc.Store(context.Background(), providedClients, &options.ParameterBag{})
	assert.NoError(t, err)

	actualClients, err := cc.Load(context.Background(), &options.ParameterBag{})
	assert.NoError(t, err)
	assert.Equal(t, providedClients, actualClients)

	err = os.Remove(ClientsCacheFileName)
	assert.NoError(t, err)
}

func assertExists(t *testing.T) {
	cc := &ClientsCache{}

	exists, err := cc.Exists(context.Background(), &options.ParameterBag{})
	assert.NoError(t, err)
	assert.False(t, exists)

	err = storeModelToFile(&ClientsCacheModel{
		Clients:   []*models.Client{},
		ValidTill: time.Now().UTC().Add(-1 * time.Hour),
	})
	assert.NoError(t, err)

	exists, err = cc.Exists(context.Background(), &options.ParameterBag{})
	assert.NoError(t, err)
	assert.False(t, exists)

	// cleanup
	err = os.Remove(ClientsCacheFileName)
	assert.NoError(t, err)

	err = storeModelToFile(&ClientsCacheModel{
		Clients:   []*models.Client{},
		ValidTill: time.Now().UTC().Add(time.Hour),
	})
	assert.NoError(t, err)
	exists, err = cc.Exists(context.Background(), &options.ParameterBag{})
	assert.NoError(t, err)
	assert.True(t, exists)

	// cleanup
	err = os.Remove(ClientsCacheFileName)
	assert.NoError(t, err)
}

func storeModelToFile(m *ClientsCacheModel) error {
	f, err := os.OpenFile(ClientsCacheFileName, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer io.CloseResourceSecure(ClientsCacheFileName, f)

	jsonEnc := json.NewEncoder(f)
	err = jsonEnc.Encode(m)

	return err
}
