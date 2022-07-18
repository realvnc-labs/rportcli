package config

import (
	"errors"
	"fmt"
	"io"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"
	"github.com/fatih/color"
)

const maxPromptIterations = 100

type PromptReader interface {
	ReadString() (string, error)
	ReadConfirmation(prompt string) (confirmed bool, err error)
	ReadPassword() (string, error)
	Output(text string)
}

func PromptRequiredValues(
	missedRequirements []ParameterRequirement,
	targetKV map[string]interface{},
	promptReader PromptReader,
) error {
	var err error
	for i := range missedRequirements {
		readValue := ""
		missedReqP := &missedRequirements[i]
		defaultStr := ""
		if missedReqP.Default != nil {
			defaultStr = fmt.Sprint(missedReqP.Default)
		}
		if defaultStr != "" && missedReqP.Validate == nil {
			readValue, err = promptValue(missedReqP, promptReader)
			if err != nil {
				return err
			}
			if readValue != "" {
				targetKV[missedReqP.Field] = readValue
			}
			continue
		}

		if missedReqP.Validate == nil {
			continue
		}

		err = missedReqP.Validate(missedReqP.Field, readValue)
		promptCount := 0
		for err != nil {
			readValue, err = promptValue(missedReqP, promptReader)
			if err != nil {
				return err
			}

			err = missedReqP.Validate(missedReqP.Field, readValue)
			if err != nil {
				color.Red(err.Error())
			}
			promptCount++
			if promptCount > maxPromptIterations {
				return fmt.Errorf("max prompt attempts %d elapsed", maxPromptIterations)
			}
		}
		targetKV[missedReqP.Field] = readValue
	}

	return nil
}

func promptValue(req *ParameterRequirement, promptReader PromptReader) (string, error) {
	promptReader.Output(req.Help + "\n")
	if req.Default != nil {
		promptReader.Output(fmt.Sprintf("Default value: %v\n", req.Default))
	}

	promptReader.Output("-> ")

	var readValue string
	var err error
	if req.IsSecure {
		readValue, err = promptReader.ReadPassword()
	} else {
		readValue, err = promptReader.ReadString()
	}
	if err != nil {
		if err == io.EOF {
			return "", errors.New(utils.InterruptMessage)
		}
		return "", err
	}

	return readValue, nil
}
