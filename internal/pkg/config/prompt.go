package config

import (
	"fmt"
	"strings"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/cli"
	"github.com/cloudradar-monitoring/rportcli/internal/pkg/utils"

	"github.com/fatih/color"
)

type PromptReader interface {
	ReadString() (string, error)
	ReadPassword() (string, error)
}

type DefaultReader struct {
}

func (dr *DefaultReader) ReadString() (string, error) {
	var val = new(string)
	_, err := fmt.Scanln(val)
	if err != nil && err.Error() == "unexpected newline" {
		return "", nil
	}
	return *val, err
}

func (dr *DefaultReader) ReadPassword() (string, error) {
	inputBytes, err := utils.ReadPassword()
	return string(inputBytes), err
}

// PromptRequiredValues will ask user for the list of required values
func PromptRequiredValues(
	missedRequirements []cli.ParameterRequirement,
	targetKV map[string]string,
	promptReader PromptReader,
) error {
	var err error
	for i := range missedRequirements {
		readValue := ""
		missedReqP := &missedRequirements[i]
		if missedReqP.Default != "" && missedReqP.Validate == nil {
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
		for err != nil {
			readValue, err = promptValue(missedReqP, promptReader)
			if err != nil {
				return err
			}

			err = missedReqP.Validate(missedReqP.Field, readValue)
			if err != nil {
				color.Red(err.Error())
			}
		}
		targetKV[missedReqP.Field] = readValue
	}

	return nil
}

func promptValue(req *cli.ParameterRequirement, promptReader PromptReader) (string, error) {
	fmt.Println(req.Help)
	if req.Default != "" {
		fmt.Printf("Default value: %s\n", req.Default)
	}

	fmt.Print("-> ")

	var readValue string
	var err error
	if req.IsSecure {
		readValue, err = promptReader.ReadPassword()
	} else {
		readValue, err = promptReader.ReadString()
	}
	if err != nil {
		if strings.Contains(err.Error(), "The handle is invalid") {
			return "", fmt.Errorf("your terminal does not support password promting " +
				"please use PowerShell or CMD or specify -p parameter explicitly")
		}
		return "", err
	}

	return readValue, nil
}
