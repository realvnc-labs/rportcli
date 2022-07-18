package controllers

import (
	"testing"

	options "github.com/breathbath/go_utils/v2/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestShouldGetNumOfClientsWithSimpleExecLog(t *testing.T) {
	sourceTestExecLog := "../../../testdata/execlog-1failed.yaml"
	logFilename := "n/a"

	params := makeBasicTestParams(logFilename)

	testExecutionInfo, err := ReadJobsFromYAML(sourceTestExecLog)
	assert.NoError(t, err)

	el := NewExecLog(params, logFilename, nil, nil)
	el.SetJobs(testExecutionInfo.Jobs)

	numClients := el.getNumClients()
	assert.Equal(t, 1, numClients)
	assert.NoError(t, nil)
}

func TestShouldFindFailedJobs(t *testing.T) {
	sourceTestExecLog := "../../../testdata/execlog-3jobs1failed.yaml"
	logFilename := "n/a"

	params := makeBasicTestParams(logFilename)

	testExecutionInfo, err := ReadJobsFromYAML(sourceTestExecLog)
	assert.NoError(t, err)

	el := NewExecLog(params, logFilename, nil, nil)
	el.SetJobs(testExecutionInfo.Jobs)

	numFailedJobs := el.getNumClientsWithFailedJobs()
	assert.Equal(t, 1, numFailedJobs)
	assert.NoError(t, nil)
}

func makeBasicTestParams(logFilename string) (params *options.ParameterBag) {
	params = config.FromValues(map[string]string{
		config.ClientIDs:    "1235",
		config.Command:      "cmd",
		config.Interpreter:  "bash",
		config.WriteExecLog: logFilename,
		config.NoPrompt:     "true",
		config.APIUser:      "Test API User",
		config.APIURL:       "http://test-server.com",
	})
	return params
}
