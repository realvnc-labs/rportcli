package controllers

import (
	"context"
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/stretchr/testify/assert"
)

type ReadChunk struct {
	Output []byte
	Err    error
}

type ReadWriterMock struct {
	itemsToRead   []ReadChunk
	itemReadIndex int
	writtenItems  []string
	writeError    error
	isClosed      bool
	closeError    error
}

func (rwm *ReadWriterMock) Read() (msg []byte, err error) {
	item := rwm.itemsToRead[rwm.itemReadIndex]

	msg = item.Output
	err = item.Err

	rwm.itemReadIndex++

	return
}

func (rwm *ReadWriterMock) Write(inputMsg []byte) (n int, err error) {
	rwm.writtenItems = append(rwm.writtenItems, string(inputMsg))
	return 0, rwm.writeError
}

func (rwm *ReadWriterMock) Close() error {
	rwm.isClosed = true
	return rwm.closeError
}

type SpinnerMock struct {
	startMsgs       []string
	updateMsgs      []string
	stopSuccessMsgs []string
	stopErrorMsgs   []string
}

func (sm *SpinnerMock) Start(msg string) {
	sm.startMsgs = append(sm.startMsgs, msg)
}

func (sm *SpinnerMock) Update(msg string) {
	sm.updateMsgs = append(sm.updateMsgs, msg)
}

func (sm *SpinnerMock) StopSuccess(msg string) {
	sm.stopSuccessMsgs = append(sm.stopSuccessMsgs, msg)
}

func (sm *SpinnerMock) StopError(msg string) {
	sm.stopErrorMsgs = append(sm.stopErrorMsgs, msg)
}

type JobRendererMock struct {
	jobToRender *models.Job
	err         error
}

func (jrm *JobRendererMock) RenderJob(j *models.Job) error {
	jrm.jobToRender = j
	return jrm.err
}

func TestInteractiveCommandExecutionSuccess(t *testing.T) {
	jobResp := models.Job{
		Jid:        "123",
		Status:     "done",
		FinishedAt: time.Now(),
		ClientID:   "123",
		Command:    "ls",
		Shell:      "sh",
		Pid:        12,
		StartedAt:  time.Now(),
		CreatedBy:  "admin",
		TimeoutSec: 1,
		Result: models.JobResult{
			Stdout: "some out",
			Stderr: "some err",
		},
	}
	jobRespBytes, err := json.Marshal(jobResp)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	rw := &ReadWriterMock{
		itemsToRead: []ReadChunk{
			{
				Output: jobRespBytes,
			},
			{
				Err: io.EOF,
			},
		},
		writtenItems: []string{},
		isClosed:     false,
	}

	pr := &PromptReaderMock{
		ReadOutputs:         []string{},
		PasswordReadOutputs: []string{},
	}

	s := &SpinnerMock{
		startMsgs:       []string{},
		updateMsgs:      []string{},
		stopSuccessMsgs: []string{},
		stopErrorMsgs:   []string{},
	}

	jr := &JobRendererMock{}

	ic := &InteractiveCommandsController{
		ReadWriter:   rw,
		PromptReader: pr,
		Spinner:      s,
		JobRenderer:  jr,
	}

	cids := "1235"
	cmd := "cmd"
	to := "1"
	gi := "333"
	ec := "1"
	err = ic.Start(context.Background(), map[string]*string{
		clientIDs:        &cids,
		command:          &cmd,
		timeout:          &to,
		groupIDs:         &gi,
		execConcurrently: &ec,
	})

	assert.NoError(t, err)

	assert.Equal(t, pr.PasswordReadCount, 0)
	assert.Equal(t, pr.ReadCount, 0)

	assert.Len(t, rw.writtenItems, 1)
	expectedCommandInput := `{"command":"cmd","client_ids":["1235"],"group_ids":["333"],"timeout_sec":1,"execute_concurrently":true}`
	assert.Equal(t, expectedCommandInput, rw.writtenItems[0])

	assert.NotNil(t, jr.jobToRender)
	actualJobRenderResult, err := json.Marshal(jr.jobToRender)
	assert.NoError(t, err)
	assert.Equal(t, string(jobRespBytes), string(actualJobRenderResult))
	assert.True(t, rw.isClosed)
	assert.Len(t, s.stopErrorMsgs, 0)
	assert.Len(t, s.stopSuccessMsgs, 2)
}

func TestInteractiveCommandExecutionWithPromptParams(t *testing.T) {
	jobResp := models.Job{
		Jid:    "123",
		Status: "done",
	}
	jobRespBytes, err := json.Marshal(jobResp)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	rw := &ReadWriterMock{
		itemsToRead: []ReadChunk{
			{
				Output: jobRespBytes,
			},
			{
				Err: io.EOF,
			},
		},
		writtenItems: []string{},
		isClosed:     false,
	}

	pr := &PromptReaderMock{
		ReadOutputs: []string{
			"123",
			"dir",
		},
		PasswordReadOutputs: []string{},
	}

	s := &SpinnerMock{
		startMsgs:       []string{},
		updateMsgs:      []string{},
		stopSuccessMsgs: []string{},
		stopErrorMsgs:   []string{},
	}

	jr := &JobRendererMock{}

	ic := &InteractiveCommandsController{
		ReadWriter:   rw,
		PromptReader: pr,
		Spinner:      s,
		JobRenderer:  jr,
	}

	err = ic.Start(context.Background(), map[string]*string{})

	assert.NoError(t, err)

	assert.Equal(t, pr.PasswordReadCount, 0)
	assert.Equal(t, pr.ReadCount, 2)

	assert.Len(t, rw.writtenItems, 1)
	expectedCommandInput := `{"command":"dir","client_ids":["123"],"timeout_sec":30,"execute_concurrently":false}`
	assert.Equal(t, expectedCommandInput, rw.writtenItems[0])

	assert.NotNil(t, jr.jobToRender)
	actualJobRenderResult, err := json.Marshal(jr.jobToRender)
	assert.NoError(t, err)
	assert.Equal(t, string(jobRespBytes), string(actualJobRenderResult))
	assert.True(t, rw.isClosed)
	assert.Len(t, s.stopErrorMsgs, 0)
	assert.Len(t, s.stopSuccessMsgs, 2)
}

func TestInteractiveCommandExecutionWithInvalidResponse(t *testing.T) {
	resp := models.ErrorResp{
		Errors: []models.Error{
			{
				Code:   "500",
				Title:  "some error",
				Detail: "some error detail",
			},
		},
	}
	jobRespBytes, err := json.Marshal(resp)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	rw := &ReadWriterMock{
		itemsToRead: []ReadChunk{
			{
				Output: jobRespBytes,
			},
			{
				Err: io.EOF,
			},
		},
		writtenItems: []string{},
		isClosed:     false,
	}

	pr := &PromptReaderMock{
		ReadOutputs:         []string{},
		PasswordReadOutputs: []string{},
	}

	s := &SpinnerMock{
		startMsgs:       []string{},
		updateMsgs:      []string{},
		stopSuccessMsgs: []string{},
		stopErrorMsgs:   []string{},
	}

	jr := &JobRendererMock{}

	ic := &InteractiveCommandsController{
		ReadWriter:   rw,
		PromptReader: pr,
		Spinner:      s,
		JobRenderer:  jr,
	}

	cids := "123"
	cmd := "ls"
	err = ic.Start(context.Background(), map[string]*string{
		clientIDs: &cids,
		command:   &cmd,
	})

	assert.Error(t, err)
	if err == nil {
		return
	}
	assert.Contains(t, err.Error(), "some error, code: 500, details: some error detail")
}
