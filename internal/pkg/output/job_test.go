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

	jr := &JobRenderer{}

	buf := &bytes.Buffer{}
	err = jr.RenderJob(buf, tunnel)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	actualRenderResult := RemoveEmptySpaces(buf.String())
	assert.Equal(
		t,
		"Command Execution Result Job ID: 123 Status: success Command Output: some std Command Error Output: Started at: 2021-01-01T00:00:01Z Finished at: 2021-01-01T00:00:01Z Client ID: cl123 Command: ls Shell: cmd Pid: 123 Timeout sec: 10 Created By: me Multi Job ID: ",
		actualRenderResult,
	)
}
