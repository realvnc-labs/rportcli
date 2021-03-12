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
	}{
		{
			Format: FormatHuman,
			ExpectedOutput: `Command Execution Result
Job ID: 123
Status: success
Command Output: some std
Command Error Output: 
Started at: 2021-01-01T00:00:01Z
Finished at: 2021-01-01T00:00:01Z
Client ID: cl123
Command: ls
Shell: cmd
Pid: 123
Timeout sec: 10
Created By: me
Multi Job ID: 
`,
		},
		{
			Format: FormatJSON,
			ExpectedOutput: `{"jid":"123","status":"success","finished_at":"2021-01-01T00:00:01Z","client_id":"cl123","command":"ls","shell":"cmd","pid":123,"started_at":"2021-01-01T00:00:01Z","created_by":"me","multi_job_id":"","timeout_sec":10,"error":"","result":{"stdout":"some std","stderr":""}}
`,
		},
		{
			Format: FormatJSONPretty,
			ExpectedOutput: `{
  "jid": "123",
  "status": "success",
  "finished_at": "2021-01-01T00:00:01Z",
  "client_id": "cl123",
  "command": "ls",
  "shell": "cmd",
  "pid": 123,
  "started_at": "2021-01-01T00:00:01Z",
  "created_by": "me",
  "multi_job_id": "",
  "timeout_sec": 10,
  "error": "",
  "result": {
    "stdout": "some std",
    "stderr": ""
  }
}
`,
		},
		{
			Format: FormatYAML,
			ExpectedOutput: `jid: "123"
status: success
finishedat: 2021-01-01T00:00:01Z
clientid: cl123
command: ls
shell: cmd
pid: 123
startedat: 2021-01-01T00:00:01Z
createdby: me
multijobid: ""
timeoutsec: 10
error: ""
result:
  stdout: some std
  stderr: ""
`,
		},
	}

	tunnel := &models.Job{
		Jid:        "123",
		Status:     "success",
		FinishedAt: timeToCheck,
		ClientID:   "cl123",
		Command:    "ls",
		Shell:      "cmd",
		Pid:        123,
		StartedAt:  timeToCheck,
		CreatedBy:  "me",
		TimeoutSec: 10,
		Result: models.JobResult{
			Stdout: "some std",
		},
	}

	for _, testCase := range testCases {
		buf := &bytes.Buffer{}
		jr := &JobRenderer{
			Writer: buf,
			Format: testCase.Format,
		}

		err = jr.RenderJob(tunnel)
		assert.NoError(t, err)
		if err != nil {
			return
		}

		assert.Equal(
			t,
			testCase.ExpectedOutput,
			buf.String(),
		)
	}
}
