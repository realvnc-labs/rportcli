package controllers

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/config"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestShouldErrIfOverwriteNotConfirmed(t *testing.T) {
	ic := &CommandsController{
		ExecutionHelper: &ExecutionHelper{
			ReadWriter:  nil,
			JobRenderer: nil,
		},
	}

	prm := &execPromptReaderMock{
		ConfirmationAnswer: false,
	}

	params := config.FromValues(map[string]string{
		config.ClientIDs:    "1235",
		config.Command:      "cmd",
		config.Interpreter:  "bash",
		config.WriteExecLog: "../../../testdata/execlog-simple1.yaml",
	})

	err := ic.Start(context.Background(), params, prm, nil)
	assert.ErrorIs(t, err, ErrOverwriteNotConfirmed)
}

func TestShouldErrorWhenSourceExecLogMissing(t *testing.T) {
	ic := &CommandsController{
		ExecutionHelper: &ExecutionHelper{
			ReadWriter:  nil,
			JobRenderer: nil,
		},
	}

	params := config.FromValues(map[string]string{
		config.ClientIDs:   "1235",
		config.Command:     "cmd",
		config.Interpreter: "bash",
		config.ReadExecLog: "missing.log",
	})

	err := ic.Start(context.Background(), params, nil, nil)
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestShouldSaveSimpleExecLog(t *testing.T) {
	eh, jobResp := makeExecutionHelperWithSimpleJob(t)

	ic := &CommandsController{
		ExecutionHelper: eh,
	}

	logFilename := "./simple-execlog.yaml"
	defer os.Remove(logFilename)

	params := config.FromValues(map[string]string{
		config.ClientIDs:    "1235",
		config.Command:      "cmd",
		config.Interpreter:  "bash",
		config.WriteExecLog: logFilename,
		config.NoPrompt:     "true",
		config.APIUser:      "Test API User",
		config.APIURL:       "http://test-server.com",
	})

	hostInfo := makeBasicTestHostInfo(t)

	err := ic.Start(context.Background(), params, nil, hostInfo)
	assert.NoError(t, err)

	jobRespAsArray := []*models.Job{
		jobResp,
	}

	execLogInfo := &ExecutionLogInfo{
		APIUser:    config.ReadAPIUser(params),
		APIURL:     config.ReadAPIURL(params),
		APIAuth:    "bearer+jwt",
		ExecutedAt: hostInfo.FetchedAt,
		ExecutedBy: hostInfo.Username,
		ExecutedOn: hostInfo.Machine,
		NumClients: 1,
		Failed:     1,
		Jobs:       jobRespAsArray,
	}

	originalJobContents, err := yaml.Marshal(execLogInfo)
	assert.NoError(t, err)

	jobContents, err := os.ReadFile(logFilename)
	assert.NoError(t, err)

	assert.Equal(t, string(originalJobContents), string(jobContents))
}

func TestShouldSaveSimpleExecLogForScript(t *testing.T) {
	eh, jobResp := makeExecutionHelperWithSimpleScriptJob(t)

	ic := &ScriptsController{
		ExecutionHelper: eh,
	}

	logFilename := "./simple-execlog.yaml"
	defer os.Remove(logFilename)

	params := config.FromValues(map[string]string{
		config.ClientIDs:    "1235",
		config.Script:       "../../../testdata/test.sh",
		config.Interpreter:  "bash",
		config.WriteExecLog: logFilename,
		config.NoPrompt:     "true",
		config.APIUser:      "Test API User",
		config.APIURL:       "http://test-server.com",
	})

	hostInfo := makeBasicTestHostInfo(t)

	err := ic.Start(context.Background(), params, nil, hostInfo)
	assert.NoError(t, err)

	jobRespAsArray := []*models.Job{
		jobResp,
	}

	execLogInfo := &ExecutionLogInfo{
		APIUser:    config.ReadAPIUser(params),
		APIURL:     config.ReadAPIURL(params),
		APIAuth:    "bearer+jwt",
		ExecutedAt: hostInfo.FetchedAt,
		ExecutedBy: hostInfo.Username,
		ExecutedOn: hostInfo.Machine,
		NumClients: 1,
		Failed:     1,
		Jobs:       jobRespAsArray,
	}

	originalJobContents, err := yaml.Marshal(execLogInfo)
	assert.NoError(t, err)

	jobContents, err := os.ReadFile(logFilename)
	assert.NoError(t, err)

	assert.Equal(t, string(originalJobContents), string(jobContents))
}

func TestShouldSaveSimpleExecLogWithConfirmation(t *testing.T) {
	eh, jobResp := makeExecutionHelperWithSimpleJob(t)

	ic := &CommandsController{
		ExecutionHelper: eh,
	}

	srcFilename := "../../../testdata/execlog-simple2.yaml"
	logFilename := "./simple-execlog.yaml"

	err := CopyFile(srcFilename, logFilename)
	assert.NoError(t, err)
	defer os.Remove(logFilename)

	params := config.FromValues(map[string]string{
		config.ClientIDs:    "1235",
		config.Command:      "cmd",
		config.Interpreter:  "bash",
		config.WriteExecLog: logFilename,
		config.APIUser:      "Test API User",
		config.APIURL:       "http://test-server.com",
	})

	prm := &execPromptReaderMock{
		ConfirmationAnswer: true,
	}

	hostInfo := makeBasicTestHostInfo(t)

	err = ic.Start(context.Background(), params, prm, hostInfo)
	assert.NoError(t, err)

	jobRespAsArray := []*models.Job{
		jobResp,
	}

	execLogInfo := &ExecutionLogInfo{
		APIUser:    config.ReadAPIUser(params),
		APIURL:     config.ReadAPIURL(params),
		APIAuth:    "bearer+jwt",
		ExecutedAt: hostInfo.FetchedAt,
		ExecutedBy: hostInfo.Username,
		ExecutedOn: hostInfo.Machine,
		NumClients: 1,
		Failed:     1,
		Jobs:       jobRespAsArray,
	}

	originalJobContents, err := yaml.Marshal(execLogInfo)
	assert.NoError(t, err)

	jobContents, err := os.ReadFile(logFilename)
	assert.NoError(t, err)

	assert.Equal(t, string(originalJobContents), string(jobContents))
}

func TestShouldSaveMoreComplexExecLog(t *testing.T) {
	sourceLogFilename := "../../../testdata/execlog-3jobs1failed.yaml"
	eh := makeExecutionHelperWithJobsFromFile(t, sourceLogFilename)

	ic := &CommandsController{
		ExecutionHelper: eh,
	}

	logFilename := "./simple-execlog.yaml"
	defer os.Remove(logFilename)

	params := config.FromValues(map[string]string{
		config.ClientIDs:    "1235",
		config.Command:      "cmd",
		config.Interpreter:  "bash",
		config.WriteExecLog: logFilename,
		config.NoPrompt:     "true",
		config.APIUser:      "Test API User",
		config.APIURL:       "http://test-server.com",
		config.APIToken:     "12345678",
	})

	hostInfo := makeBasicTestHostInfo(t)

	err := ic.Start(context.Background(), params, nil, hostInfo)
	assert.NoError(t, err)

	ec, err := os.ReadFile(sourceLogFilename)
	assert.NoError(t, err)
	ecy := ExecutionLogInfo{}
	err = yaml.Unmarshal(ec, &ecy)
	assert.NoError(t, err)

	ac, err := os.ReadFile(logFilename)
	assert.NoError(t, err)
	acy := ExecutionLogInfo{}
	err = yaml.Unmarshal(ac, &acy)
	assert.NoError(t, err)

	assert.Exactly(t, ecy, acy)
}

func TestShouldErrWhenClientIDsNotConfirmed(t *testing.T) {
	sourceLogFilename := "../../../testdata/execlog-3jobs1failed.yaml"
	eh := makeExecutionHelperWithJobsFromFile(t, sourceLogFilename)

	ic := &CommandsController{
		ExecutionHelper: eh,
	}

	logFilename := "./simple-execlog.yaml"
	defer os.Remove(logFilename)

	params := config.FromValues(map[string]string{
		config.ClientIDs:    "1235",
		config.Command:      "cmd",
		config.Interpreter:  "bash",
		config.WriteExecLog: logFilename,
		config.ReadExecLog:  "../../../testdata/execlog-3jobs1failed.yaml",
		config.APIUser:      "Test API User",
		config.APIURL:       "http://test-server.com",
		config.APIToken:     "12345678",
	})

	hostInfo := makeBasicTestHostInfo(t)

	prm := &execPromptReaderMock{
		ConfirmationAnswer: false,
	}

	err := ic.Start(context.Background(), params, prm, hostInfo)
	assert.ErrorIs(t, err, ErrClientIDsNotConfirmed)
}

func TestShouldSaveExecLogForFailedClientIDs(t *testing.T) {
	sourceLogFilename := "../../../testdata/execlog-3jobs1failed.yaml"
	expectedLogResultsFilename := "../../../testdata/execlog-1failed.yaml"
	eh := makeExecutionHelperWithJobsFromFile(t, expectedLogResultsFilename)

	ic := &CommandsController{
		ExecutionHelper: eh,
	}

	logFilename := "./simple-execlog.yaml"
	defer os.Remove(logFilename)

	params := config.FromValues(map[string]string{
		config.ClientIDs:    "1235",
		config.Command:      "cmd",
		config.Interpreter:  "bash",
		config.WriteExecLog: logFilename,
		config.ReadExecLog:  sourceLogFilename,
		config.APIUser:      "Test API User",
		config.APIURL:       "http://test-server.com",
		config.APIToken:     "12345678",
	})

	hostInfo := makeBasicTestHostInfo(t)

	prm := &execPromptReaderMock{
		ConfirmationAnswer: true,
	}

	err := ic.Start(context.Background(), params, prm, hostInfo)
	assert.NoError(t, err)

	ec, err := os.ReadFile(expectedLogResultsFilename)
	assert.NoError(t, err)
	ecy := ExecutionLogInfo{}
	err = yaml.Unmarshal(ec, &ecy)
	assert.NoError(t, err)

	ac, err := os.ReadFile(logFilename)
	assert.NoError(t, err)
	acy := ExecutionLogInfo{}
	err = yaml.Unmarshal(ac, &acy)
	assert.NoError(t, err)

	assert.Exactly(t, ecy, acy)
}

func makeSimpleJob(t *testing.T) (j *models.Job) {
	return &models.Job{
		Jid:        "2c31a80f-389f-40aa-9eb0-98f907dee16e",
		Status:     "failed",
		FinishedAt: MakeTime(t, "2022-07-13T05:57:07.778244133Z"),
		ClientID:   "76376f704ab0429eb6cb141e6f34ed75",
		ClientName: "Kathy-Daniels",
		Command:    "pwd",
		Cwd:        "",
		Pid:        124058,
		StartedAt:  MakeTime(t, "2022-07-13T05:57:07.777202262Z"),
		CreatedBy:  "rs",
		MultiJobID: "350e22cb-79f7-4032-9155-c8a647f76020",
		TimeoutSec: 30,
		Error:      "",
		Result: models.JobResult{
			Stdout: "/",
			Stderr: "",
		},
		IsSudo:      false,
		IsScript:    false,
		Interpreter: "",
	}
}

func makeSimpleScriptJob(t *testing.T) (j *models.Job) {
	return &models.Job{
		Jid:        "2c31a80f-389f-40aa-9eb0-98f907dee16e",
		Status:     "failed",
		FinishedAt: MakeTime(t, "2022-07-13T05:57:07.778244133Z"),
		ClientID:   "76376f704ab0429eb6cb141e6f34ed75",
		ClientName: "Kathy-Daniels",
		Command:    "script.sh",
		Cwd:        "",
		Pid:        124058,
		StartedAt:  MakeTime(t, "2022-07-13T05:57:07.777202262Z"),
		CreatedBy:  "rs",
		MultiJobID: "350e22cb-79f7-4032-9155-c8a647f76020",
		TimeoutSec: 30,
		Error:      "",
		Result: models.JobResult{
			Stdout: "/",
			Stderr: "",
		},
		IsSudo:      false,
		IsScript:    true,
		Interpreter: "",
	}
}

func makeBasicTestHostInfo(t *testing.T) (hostInfo *config.HostInfo) {
	return &config.HostInfo{
		Username:  "test_user",
		Machine:   "test_machine",
		FetchedAt: MakeTime(t, "2022-07-14T15:22:57.39771+07:00"),
		// FetchedAt: MakeTime(t, "2022-07-13T05:57:06.777202262Z"),
	}
}

func makeExecutionHelperWithSimpleJob(t *testing.T) (eh *ExecutionHelper, jobResp *models.Job) {
	jobResp = makeSimpleJob(t)
	rw := makeReadWriterMockFromJobs(t, []*models.Job{jobResp})

	jr := &JobRendererMock{}
	eh = &ExecutionHelper{
		ReadWriter:  rw,
		JobRenderer: jr,
	}

	return eh, jobResp
}

func makeExecutionHelperWithSimpleScriptJob(t *testing.T) (eh *ExecutionHelper, jobResp *models.Job) {
	jobResp = makeSimpleScriptJob(t)
	rw := makeReadWriterMockFromJobs(t, []*models.Job{jobResp})

	jr := &JobRendererMock{}
	eh = &ExecutionHelper{
		ReadWriter:  rw,
		JobRenderer: jr,
	}

	return eh, jobResp
}

func makeExecutionHelperWithJobsFromFile(t *testing.T, sampleLogFilename string) (eh *ExecutionHelper) {
	rw, _ := makeReadWriterMockFromSampleExecLog(t, sampleLogFilename)

	jr := &JobRendererMock{}
	eh = &ExecutionHelper{
		ReadWriter:  rw,
		JobRenderer: jr,
	}

	return eh
}

func makeReadWriterMockFromSampleExecLog(t *testing.T, sampleLogFilename string) (m *ReadWriterMock, jobs []*models.Job) {
	logInfo, err := ReadJobsFromYAML(sampleLogFilename)
	assert.NoError(t, err)

	rwMock := makeReadWriterMockFromJobs(t, logInfo.Jobs)
	return rwMock, logInfo.Jobs
}

func makeReadWriterMockFromJobs(t *testing.T, jobs []*models.Job) (m *ReadWriterMock) {
	readChunks := []ReadChunk{}
	for _, job := range jobs {
		jobBytes, err := json.Marshal(job)
		assert.NoError(t, err)

		rc := ReadChunk{
			Output: jobBytes,
		}

		readChunks = append(readChunks, rc)
	}

	readChunks = append(readChunks, ReadChunk{
		Err: io.EOF,
	})

	rwMock := &ReadWriterMock{
		itemsToRead:  readChunks,
		writtenItems: []string{},
		isClosed:     false,
	}
	return rwMock
}
