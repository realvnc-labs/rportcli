package output

import (
	"fmt"
	"io"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
)

type JobRenderer struct {
}

func (jr *JobRenderer) RenderJob(rw io.Writer, t *models.Job) error {
	if t == nil {
		return nil
	}
	_, err := rw.Write([]byte("Command Execution Result\n"))
	if err != nil {
		return err
	}

	for _, kv := range t.KeyValues() {
		_, err = fmt.Fprintf(rw, "%s: %s\n", kv.Key, kv.Value)
		if err != nil {
			return err
		}
	}

	return nil
}
