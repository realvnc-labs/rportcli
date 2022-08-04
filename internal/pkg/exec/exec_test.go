package exec

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecutor(t *testing.T) {
	stdOut := &bytes.Buffer{}
	const filePath = "file123"
	e := &Executor{
		CommandProvider: func(fp string) (cmd string, args []string) {
			assert.Equal(t, filePath, fp)
			return "echo", []string{"123"}
		},
		StdOut: stdOut,
	}

	err := e.StartDefaultApp(filePath)
	assert.NoError(t, err)
	assert.Equal(t, "123\n", stdOut.String())
}
