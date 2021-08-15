package controllers

import (
	"context"

	options "github.com/breathbath/go_utils/v2/pkg/config"
)

type CommandsController struct {
	*ExecutionHelper
}

func (cc *CommandsController) Start(ctx context.Context, params *options.ParameterBag) error {
	return cc.execute(ctx, params, "", params.ReadString(Interpreter, ""))
}
