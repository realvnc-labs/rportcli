package controllers

import (
	"context"

	options "github.com/breathbath/go_utils/v2/pkg/config"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
)

type PromptReaderMock struct {
	ReadCount           int
	PasswordReadCount   int
	ReadOutputs         []string
	PasswordReadOutputs []string
	ErrToGive           error
	Inputs              []string
}

func (prm *PromptReaderMock) ReadString() (string, error) {
	prm.ReadCount++

	if len(prm.ReadOutputs) < prm.ReadCount {
		return "", prm.ErrToGive
	}

	return prm.ReadOutputs[prm.ReadCount-1], prm.ErrToGive
}

func (prm *PromptReaderMock) ReadPassword() (string, error) {
	prm.PasswordReadCount++

	if len(prm.PasswordReadOutputs) < prm.PasswordReadCount {
		return "", prm.ErrToGive
	}

	return prm.PasswordReadOutputs[prm.PasswordReadCount-1], prm.ErrToGive
}

func (prm *PromptReaderMock) Output(text string) {
	prm.Inputs = append(prm.Inputs, text)
}

type ClientSearchMock struct {
	searchTermGiven string
	clientsToGive   []*models.Client
	errorToGive     error
}

func (csm *ClientSearchMock) Search(ctx context.Context, term string, params *options.ParameterBag) (foundCls []*models.Client, err error) {
	csm.searchTermGiven = term
	return csm.clientsToGive, csm.errorToGive
}

func (csm *ClientSearchMock) FindOne(ctx context.Context, searchTerm string, params *options.ParameterBag) (*models.Client, error) {
	csm.searchTermGiven = searchTerm
	if len(csm.clientsToGive) > 0 {
		return csm.clientsToGive[0], csm.errorToGive
	}
	return &models.Client{}, csm.errorToGive
}
