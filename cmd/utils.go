package cmd

import (
	"bufio"
	"context"
	"os"
	"os/signal"
	"syscall"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/breathbath/go_utils/v2/pkg/env"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/auth"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/controllers"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/output"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"
	"github.com/spf13/cobra"
)

func makeRunContext() (ctx context.Context, cancel context.CancelFunc, sigs chan os.Signal) {
	ctx, cancel = buildContext(context.Background())

	sigs = make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	return ctx, cancel, sigs
}

func loadParams(cmd *cobra.Command,
	sigs chan os.Signal,
	reqs []config.ParameterRequirement) (params *options.ParameterBag, err error) {
	promptReader := &utils.PromptReader{
		Sc:              bufio.NewScanner(os.Stdin),
		SigChan:         sigs,
		PasswordScanner: utils.ReadPassword,
	}

	params, err = config.LoadParamsFromFileAndEnvAndFlagsAndPrompt(cmd, reqs, promptReader)
	return params, err
}

func newWsURLProvider(params *options.ParameterBag, baseRportURL string) (p *api.WsURLProvider) {
	tokenValidity := env.ReadEnvInt(config.SessionValiditySecondsEnvVar, api.DefaultTokenValiditySeconds)

	p = &api.WsURLProvider{
		BaseURL: baseRportURL,
		TokenProvider: func() (token string, err error) {
			return auth.GetToken(params)
		},
		TokenValiditySeconds: tokenValidity,
	}
	return p
}

func newWsClient(ctx context.Context, urlBuilder utils.WsURLBuilder) (wsc *utils.WsClient, err error) {
	wsc, err = utils.NewWsClient(ctx, urlBuilder)
	if err != nil {
		return nil, err
	}

	return wsc, nil
}

func newExecutionHelper(params *options.ParameterBag,
	wsc *utils.WsClient,
	rportAPI *api.Rport) (helper *controllers.ExecutionHelper) {
	isFullJobOutput := params.ReadBool(config.IsFullOutput, false)
	helper = &controllers.ExecutionHelper{
		ReadWriter: wsc,
		JobRenderer: &output.JobRenderer{
			Writer:       os.Stdout,
			Format:       getOutputFormat(),
			IsFullOutput: isFullJobOutput,
		},
		Rport: rportAPI,
	}
	return helper
}
