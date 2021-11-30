package controllers

import (
	"context"
	"testing"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/api"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/output"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MeAPIMock struct {
	mock.Mock
}

func (mam *MeAPIMock) Me(ctx context.Context) (user api.UserResponse, err error) {
	args := mam.Called(ctx)

	return args.Get(0).(api.UserResponse), args.Error(1)
}

type MeRendererMock struct {
	mock.Mock
}

func (mr *MeRendererMock) RenderMe(os output.KvProvider) error {
	args := mr.Called(os)

	return args.Error(0)
}

func TestMeSuccess(t *testing.T) {
	ctx := context.Background()

	renderMock := &MeRendererMock{}
	renderMock.On("RenderMe", mock.Anything).Return(nil)

	userResponseGiven := api.UserResponse{
		Data: models.Me{
			Username:    "login",
			Groups:      []string{"group1", "group2"},
			TwoFASendTo: "m@me.com",
		},
	}
	apiMock := &MeAPIMock{}
	apiMock.On("Me", ctx).Return(userResponseGiven, nil)

	cl := &MeController{
		Rport:      apiMock,
		MeRenderer: renderMock,
	}

	err := cl.Me(ctx)
	require.NoError(t, err)

	apiMock.AssertCalled(t, "Me", ctx)

	renderMock.AssertCalled(t, "RenderMe", &userResponseGiven.Data)
}
