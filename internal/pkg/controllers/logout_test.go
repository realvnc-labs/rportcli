package controllers

import (
	"context"
	"errors"
	"testing"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var defaultLogoutParams = config.FromValues(map[string]string{
	config.ServerURL: "some.srv",
	config.Token:     "some tok",
})

type LogoutAPIMock struct {
	mock.Mock
}

func (lam *LogoutAPIMock) Logout(ctx context.Context) (err error) {
	args := lam.Called(ctx)

	return args.Error(0)
}

type ConfigWriterMock struct {
	mock.Mock
	paramsGiven *options.ParameterBag
}

func (cwm *ConfigWriterMock) WriteConfig(params *options.ParameterBag) (err error) {
	args := cwm.Called(params)
	cwm.paramsGiven = params
	return args.Error(0)
}

func (cwm *ConfigWriterMock) DeleteConfig() (err error) {
	args := cwm.Called()
	return args.Error(0)
}

func TestLogoutSuccess(t *testing.T) {
	ctx := context.Background()

	apiMock := &LogoutAPIMock{}
	apiMock.On("Logout", ctx).Return(nil)

	configWriterMock := &ConfigWriterMock{}
	configWriterMock.On("DeleteConfig").Return(nil)

	cl := NewLogoutController(apiMock, configWriterMock.DeleteConfig)

	err := cl.Logout(ctx, defaultLogoutParams)
	require.NoError(t, err)

	apiMock.AssertCalled(t, "Logout", ctx)

	configWriterMock.AssertCalled(t, "DeleteConfig")

	err = cl.Logout(ctx, &options.ParameterBag{})
	require.NoError(t, err)
}

func TestLogoutWriteConfigError(t *testing.T) {
	ctx := context.Background()

	apiMock := &LogoutAPIMock{}
	apiMock.On("Logout", ctx).Return(nil)

	configWriterMock := &ConfigWriterMock{}
	configWriterMock.On("DeleteConfig").Return(errors.New("some config error"))

	cl := NewLogoutController(apiMock, configWriterMock.DeleteConfig)

	err := cl.Logout(ctx, defaultLogoutParams)
	require.EqualError(t, err, "some config error")
}

func TestLogoutAPIError(t *testing.T) {
	ctx := context.Background()

	apiMock := &LogoutAPIMock{}
	apiMock.On("Logout", ctx).Return(errors.New("some api error"))

	configWriterMock := &ConfigWriterMock{}
	configWriterMock.On("DeleteConfig").Return(nil)

	cl := NewLogoutController(apiMock, configWriterMock.DeleteConfig)

	err := cl.Logout(ctx, defaultLogoutParams)
	require.EqualError(t, err, "some api error")
}
