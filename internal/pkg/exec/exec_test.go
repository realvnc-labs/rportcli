package exec

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/recorder"

	"github.com/stretchr/testify/assert"
)

func TestExecutor(t *testing.T) {
	r := recorder.NewCmdRecorder()
	err := StartDefaultApp("info.txt")
	require.NoError(t, err)
	cmdRecords := r.GetRecords()

	assert.Len(t, cmdRecords, 1)
	assert.Equal(t, OpenCmd+" info.txt", cmdRecords[0])
}
