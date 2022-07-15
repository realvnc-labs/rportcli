package controllers

import (
	"context"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
)

type CommandsController struct {
	*ExecutionHelper
}

func (cc *CommandsController) Start(ctx context.Context,
	params *options.ParameterBag,
	promptReader config.PromptReader,
	hostInfo *config.HostInfo) error {
	return cc.execute(ctx, params, "", params.ReadString(config.Interpreter, ""), promptReader, hostInfo)
}
