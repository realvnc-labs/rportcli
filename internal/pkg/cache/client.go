package cache

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"

	"github.com/breathbath/go_utils/v2/pkg/io"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
)

const ClientsCacheFileName = "clients.json"
const DefaultCacheValidityHours = 24

type ClientsCacheModel struct {
	Clients   []*models.Client `json:"clients"`
	ValidTill time.Time       `json:"valid_till"`
}

type ClientsCache struct {
}

func (cc *ClientsCache) Store(ctx context.Context, cls []*models.Client, params *options.ParameterBag) error {
	filePath := cc.getFilePath(params)

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer io.CloseResourceSecure(ClientsCacheFileName, f)

	validHours := params.ReadInt(config.CacheValidityEnvVar, DefaultCacheValidityHours)
	dataToStore := ClientsCacheModel{
		Clients:   cls,
		ValidTill: time.Now().UTC().Add(time.Hour * time.Duration(validHours)),
	}

	jsonEnc := json.NewEncoder(f)
	err = jsonEnc.Encode(dataToStore)
	if err != nil {
		return err
	}

	return nil
}

func (cc *ClientsCache) Exists(ctx context.Context, params *options.ParameterBag) (bool, error) {
	clc, err := cc.loadFromFile(params)
	if err != nil {
		return false, err
	}

	if clc == nil {
		return false, nil
	}

	return clc.ValidTill.After(time.Now().UTC()), nil
}

func (cc *ClientsCache) Load(ctx context.Context, params *options.ParameterBag) (cls []*models.Client, err error) {
	clc, err := cc.loadFromFile(params)
	if err != nil {
		return nil, err
	}

	cls = append(cls, clc.Clients...)
	return cls, nil
}

func (cc *ClientsCache) getFilePath(params *options.ParameterBag) string {
	return filepath.Join(params.ReadString(config.CacheFolderEnvVar, ""), ClientsCacheFileName)
}

func (cc *ClientsCache) loadFromFile(params *options.ParameterBag) (clc *ClientsCacheModel, err error) {
	clc = &ClientsCacheModel{}
	filePath := cc.getFilePath(params)

	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer io.CloseResourceSecure(ClientsCacheFileName, jsonFile)

	jsonDecoder := json.NewDecoder(jsonFile)
	err = jsonDecoder.Decode(clc)
	if err != nil {
		return
	}

	return
}
