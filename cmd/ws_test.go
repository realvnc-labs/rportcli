package cmd

import (
	"testing"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestShouldAddAuthHeaderWhenAPIToken(t *testing.T) {
	params := config.FromValues(map[string]string{
		"api_user":  "admin",
		"api_token": "123478-123478-123478-123478",
	})
	reqHeader, err := addAuthHeaderIfAPIToken(params)
	assert.NoError(t, err)
	assert.NotNil(t, reqHeader)
	header := reqHeader.Get("Authorization")
	assert.Equal(t, "Basic YWRtaW46MTIzNDc4LTEyMzQ3OC0xMjM0NzgtMTIzNDc4", header)
}
