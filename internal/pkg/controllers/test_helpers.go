package controllers

import (
	"os"
	"testing"
	"time"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

type execPromptReaderMock struct {
	ConfirmationAnswer bool
}

func (prm *execPromptReaderMock) ReadString() (string, error) {
	return "", nil
}

func (prm *execPromptReaderMock) ReadPassword() (string, error) {
	return "", nil
}

func (prm *execPromptReaderMock) Output(text string) {
}

func (prm *execPromptReaderMock) ReadConfirmation(prompt string) (confirmed bool, err error) {
	return prm.ConfirmationAnswer, nil
}

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

type JobRendererMock struct {
	jobToRender *models.Job
	err         error
}

func (jrm *JobRendererMock) RenderJob(j *models.Job) error {
	jrm.jobToRender = j
	return jrm.err
}

func ReadJobsFromYAML(sourceJobsFilename string) (prevExecutionLogInfo *ExecutionLogInfo, err error) {
	fileContents, err := os.ReadFile(sourceJobsFilename)
	if err != nil {
		return nil, err
	}

	prevExecutionLogInfo = &ExecutionLogInfo{}
	err = yaml.Unmarshal(fileContents, &prevExecutionLogInfo)
	if err != nil {
		return nil, err
	}

	return prevExecutionLogInfo, nil
}

func CopyFile(src, dst string) (err error) {
	srcContents, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	err = os.WriteFile(dst, srcContents, 0600)
	if err != nil {
		return err
	}
	return nil
}

func MakeTime(t *testing.T, timeStr string) (atTime time.Time) {
	t.Helper()
	atTime, err := time.Parse(time.RFC3339, timeStr)
	assert.NoError(t, err)
	return atTime
}
