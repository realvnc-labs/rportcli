package config

import (
	options "github.com/breathbath/go_utils/utils/config"
	"github.com/sirupsen/logrus"
)

var Params *options.ParameterBag

func init() {
	var err error
	Params, err = GetConfig()
	if err != nil {
		logrus.Debugf("failed to read config: %v", err)
		Params = GetDefaultConfig()
	}
}
