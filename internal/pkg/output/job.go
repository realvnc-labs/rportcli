package output

import (
	"fmt"
	"io"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
)

type JobRenderer struct {
	Writer io.Writer
	Format string
}

func (jr *JobRenderer) RenderJob(j *models.Job) error {
	return RenderByFormat(
		jr.Format,
		jr.Writer,
		j,
		func() error {
			return jr.renderJobInHumanFormat(j)
		},
	)
}

func (jr *JobRenderer) renderJobInHumanFormat(j *models.Job) error {
	if j == nil {
		return nil
	}
	_, err := jr.Writer.Write([]byte("Command Execution Result\n"))
	if err != nil {
		return err
	}

	for _, kv := range j.KeyValues() {
		_, err = fmt.Fprintf(jr.Writer, "%s: %s\n", kv.Key, kv.Value)
		if err != nil {
			return err
		}
	}

	return nil
}
