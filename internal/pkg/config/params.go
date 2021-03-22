package config

import (
	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/sirupsen/logrus"
)

var Params *options.ParameterBag

func init() {
	var err error
	Params, err = LoadConfig()
	if err != nil {
		logrus.Debugf("failed to read config: %v", err)
		Params = GetDefaultConfig()
	}
}
