package recorder

import (
	"os/exec"
	"strings"

	"bou.ke/monkey"
)

type CmdRecorder struct {
	records []string
}

func NewCmdRecorder() *CmdRecorder {
	cr := &CmdRecorder{}
	var fakeResponse = exec.Command("true")
	monkey.Patch(exec.Command, func(cmd string, args ...string) *exec.Cmd {
		cr.records = append(cr.records, cmd+" "+strings.Join(args, " "))
		return fakeResponse
	})
	return cr
}

func (cr *CmdRecorder) Stop() {
	monkey.Unpatch(exec.Command)
	cr.records = []string{}
}

func (cr *CmdRecorder) GetRecords() (records []string) {
	records = cr.records
	cr.records = []string{}
	return records
}
