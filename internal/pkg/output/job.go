package output

import (
	"io"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
)

type JobRenderer struct {
}

func (jr *JobRenderer) RenderJob(rw io.Writer, t *models.Job) error {
	if t == nil {
		return nil
	}

	err := RenderHeader(rw, "Command Execution Result")
	if err != nil {
		return err
	}

	RenderKeyValues(rw, t)

	return nil
}
