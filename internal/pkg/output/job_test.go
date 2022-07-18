package output

import (
	"bytes"
	"testing"
	"time"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestRenderJob(t *testing.T) {
	timeToCheck, err := time.Parse("2006-01-02T15:04:05", "2021-01-01T00:00:01")
	assert.NoError(t, err)

	testCases := []struct {
		Format         string
		ExpectedOutput string
		IsFullOutput   bool
	}{
		{
			Format: FormatHuman,
			ExpectedOutput: `Client ID: cl123
Client Name: some cl name
    Command Execution Result
    Job ID: 123
    Status: success
    Command Output:
      some std
    Command Error Output:
      no
    Started at: 2021-01-01 00:00:01 +0000 UTC
    Finished at: 2021-01-01 00:00:01 +0000 UTC
    Command: ls
    Interpreter: cmd
    Pid: 123
    Timeout sec: 10
    Created By: me
    Multi Job ID: multi1234
    Cwd: here
    Is sudo: true
    Error: no
    Is script: true
    Status: success
`,
			IsFullOutput: true,
		},
		{
			Format: FormatHuman,
			ExpectedOutput: `some cl name
    some std
    no
`,
			IsFullOutput: false,
		},
		{
			Format: FormatJSON,
			ExpectedOutput: `{"jid":"123","status":"success","finished_at":"2021-01-01T00:00:01Z","client_id":"cl123","client_name":"some cl name","command":"ls","cwd":"here","pid":123,"started_at":"2021-01-01T00:00:01Z","created_by":"me","multi_job_id":"multi1234","timeout_sec":10,"error":"no","result":{"stdout":"some std","stderr":""},"is_sudo":true,"is_script":true,"interpreter":"cmd"}
`,
		},
		{
			Format: FormatJSONPretty,
			ExpectedOutput: `{
  "jid": "123",
  "status": "success",
  "finished_at": "2021-01-01T00:00:01Z",
  "client_id": "cl123",
  "client_name": "some cl name",
  "command": "ls",
  "cwd": "here",
  "pid": 123,
  "started_at": "2021-01-01T00:00:01Z",
  "created_by": "me",
  "multi_job_id": "multi1234",
  "timeout_sec": 10,
  "error": "no",
  "result": {
    "stdout": "some std",
    "stderr": ""
  },
  "is_sudo": true,
  "is_script": true,
  "interpreter": "cmd"
}
`,
		},
		{
			Format: FormatYAML,
			ExpectedOutput: `jid: "123"
status: success
finished_at: 2021-01-01T00:00:01Z
client_id: cl123
client_name: some cl name
command: ls
cwd: here
pid: 123
started_at: 2021-01-01T00:00:01Z
created_by: me
multi_job_id: multi1234
timeout_sec: 10
error: "no"
result:
  stdout: some std
  stderr: ""
is_sudo: true
is_script: true
interpreter: cmd
`,
		},
	}

	tunnel := &models.Job{
		Jid:         "123",
		Status:      "success",
		FinishedAt:  timeToCheck,
		ClientID:    "cl123",
		ClientName:  "some cl name",
		Command:     "ls",
		Cwd:         "here",
		Interpreter: "cmd",
		Pid:         123,
		StartedAt:   timeToCheck,
		CreatedBy:   "me",
		MultiJobID:  "multi1234",
		TimeoutSec:  10,
		Error:       "no",
		Result: models.JobResult{
			Stdout: "some std",
		},
		IsSudo:   true,
		IsScript: true,
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run("render_"+tc.Format, func(t *testing.T) {
			buf := &bytes.Buffer{}
			jr := &JobRenderer{
				Writer:       buf,
				Format:       tc.Format,
				IsFullOutput: tc.IsFullOutput,
			}

			err = jr.RenderJob(tunnel)
			assert.NoError(t, err)
			if err != nil {
				return
			}

			assert.Equal(
				t,
				tc.ExpectedOutput,
				buf.String(),
			)
		})
	}
}
