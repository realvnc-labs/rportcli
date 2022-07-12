package applog

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestShouldHaveCleanWarningMessages(t *testing.T) {
	f := Init()
	logEntry := &logrus.Entry{
		Level:   logrus.WarnLevel,
		Message: "Test Message",
	}
	expectedLine := "WARNING: " + logEntry.Message + "\n"

	result, err := f.Format(logEntry)

	assert.NoError(t, err)
	assert.Equal(t, expectedLine, string(result))
}

func TestShouldHaveCleanFatalMessage(t *testing.T) {
	f := Init()
	logEntry := &logrus.Entry{
		Level:   logrus.FatalLevel,
		Message: "Test Message",
	}

	expectedLine := "FATAL: " + logEntry.Message + "\n"

	result, err := f.Format(logEntry)

	assert.NoError(t, err)
	assert.Equal(t, expectedLine, string(result))
}

func TestShouldHaveDetailedDebugMessage(t *testing.T) {
	f := Init()
	logEntry := &logrus.Entry{
		Level:   logrus.DebugLevel,
		Message: "Test Message",
	}

	result, err := f.Format(logEntry)

	assert.NoError(t, err)
	assert.Contains(t, string(result), "time=")
	assert.Contains(t, string(result), "level=")
	assert.Contains(t, string(result), logEntry.Message)
}
