package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/fatih/color"
)

type JobRenderer struct {
	Writer       io.Writer
	Format       string
	IsFullOutput bool
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

func (jr *JobRenderer) genShiftedMultilineStr(input, shiftStr string) string {
	input = strings.Trim(input, "\n")
	input = strings.TrimSpace(input)
	if input == "" {
		return input
	}

	inputLines := strings.Split(input, "\n")
	for k, inputLine := range inputLines {
		inputLine = strings.Trim(inputLine, "\r")
		inputLines[k] = shiftStr + inputLine
	}

	return strings.Join(inputLines, "\n")
}

func (jr *JobRenderer) formatError(j *models.Job, shiftStr string) string {
	errOutput := ""
	if j.Error != "" {
		errOutput = j.Error
	}
	if j.Result.Stderr != "" {
		sep := ""
		if errOutput != "" {
			sep = " "
		}
		errOutput += sep + j.Result.Stderr
	}

	errOutput = strings.Trim(errOutput, "\n")
	errOutput = strings.TrimSpace(errOutput)

	if errOutput == "" {
		return errOutput
	}

	inputLines := strings.Split(errOutput, "\n")
	for k := range inputLines {
		l := inputLines[k]
		l = strings.Trim(l, "\r")
		inputLines[k] = shiftStr + l
	}

	return strings.Join(inputLines, "\n")
}

func (jr *JobRenderer) extractClientNameOrID(j *models.Job) string {
	if j.ClientName != "" {
		return j.ClientName
	}

	return j.ClientID
}

func (jr *JobRenderer) renderJobInHumanFormat(j *models.Job) error {
	if j == nil {
		return nil
	}

	var outputs []string

	if !jr.IsFullOutput {
		_, err := fmt.Fprintln(jr.Writer, jr.extractClientNameOrID(j))
		if err != nil {
			return err
		}

		stdOut := jr.genShiftedMultilineStr(j.Result.Stdout, "    ")
		if stdOut != "" {
			co := color.New(color.FgGreen)
			_, err = co.Fprintln(jr.Writer, stdOut)
			if err != nil {
				return err
			}
		}
		stdErr := jr.formatError(j, "    ")
		if stdErr != "" {
			co := color.New(color.FgRed)
			_, err = co.Fprintln(jr.Writer, stdErr)
			if err != nil {
				return err
			}
		}
		return nil
	}

	outputs = []string{
		fmt.Sprintf("Client ID: %s", j.ClientID),
		fmt.Sprintf("Client Name: %s", j.ClientName),
		"    Command Execution Result",
		fmt.Sprintf("    Job ID: %s", j.Jid),
		fmt.Sprintf("    Status: %s", j.Status),
		"    Command Output:",
	}

	stdOut := jr.genShiftedMultilineStr(j.Result.Stdout, "      ")
	if stdOut != "" {
		outputs = append(outputs, stdOut)
	}

	outputs = append(outputs, "    Command Error Output:")
	errOut := jr.formatError(j, "      ")
	if errOut != "" {
		outputs = append(outputs, errOut)
	}

	outputs2 := []string{
		fmt.Sprintf("    Started at: %v", j.StartedAt),
		fmt.Sprintf("    Finished at: %v", j.FinishedAt),
		fmt.Sprintf("    Command: %s", j.Command),
		fmt.Sprintf("    Interpreter: %s", j.Interpreter),
		fmt.Sprintf("    Pid: %d", j.Pid),
		fmt.Sprintf("    Timeout sec: %d", j.TimeoutSec),
		fmt.Sprintf("    Created By: %s", j.CreatedBy),
		fmt.Sprintf("    Multi Job ID: %s", j.MultiJobID),
		fmt.Sprintf("    Cwd: %s", j.Cwd),
		fmt.Sprintf("    Is sudo: %v", j.IsSudo),
		fmt.Sprintf("    Error: %v", j.Error),
		fmt.Sprintf("    Is script: %v", j.IsScript),
		fmt.Sprintf("    Status: %v", j.Status),
	}
	outputs = append(outputs, outputs2...)

	for _, output := range outputs {
		_, err := fmt.Fprintln(jr.Writer, output)
		if err != nil {
			return err
		}
	}

	return nil
}
