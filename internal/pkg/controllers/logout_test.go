package controllers

import (
	"context"
	"errors"
	"testing"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
	"github.com/magiconair/properties/assert"
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

func TestLogoutSuccess(t *testing.T) {
	ctx := context.Background()

	apiMock := &LogoutAPIMock{}
	apiMock.On("Logout", ctx).Return(nil)

	configWriterMock := &ConfigWriterMock{}
	configWriterMock.On("WriteConfig", mock.Anything).Return(nil)

	cl := NewLogoutController(apiMock, configWriterMock.WriteConfig)

	err := cl.Logout(ctx, defaultLogoutParams)
	require.NoError(t, err)

	apiMock.AssertCalled(t, "Logout", ctx)
	require.NotNil(t, configWriterMock.paramsGiven)

	assert.Equal(t, configWriterMock.paramsGiven.ReadString(config.Token, "some1"), "")
	assert.Equal(t, configWriterMock.paramsGiven.ReadString(config.ServerURL, "some2"), "some.srv")

	err = cl.Logout(ctx, &options.ParameterBag{})
	require.NoError(t, err)
}

func TestLogoutWriteConfigError(t *testing.T) {
	ctx := context.Background()

	apiMock := &LogoutAPIMock{}
	apiMock.On("Logout", ctx).Return(nil)

	configWriterMock := &ConfigWriterMock{}
	configWriterMock.On("WriteConfig", mock.Anything).Return(errors.New("some config error"))

	cl := NewLogoutController(apiMock, configWriterMock.WriteConfig)

	err := cl.Logout(ctx, defaultLogoutParams)
	require.EqualError(t, err, "some config error")
}

func TestLogoutAPIError(t *testing.T) {
	ctx := context.Background()

	apiMock := &LogoutAPIMock{}
	apiMock.On("Logout", ctx).Return(errors.New("some api error"))

	configWriterMock := &ConfigWriterMock{}
	configWriterMock.On("WriteConfig", mock.Anything).Return(nil)

	cl := NewLogoutController(apiMock, configWriterMock.WriteConfig)

	err := cl.Logout(ctx, defaultLogoutParams)
	require.EqualError(t, err, "some api error")
}
