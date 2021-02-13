package config

import (
	"fmt"
	"strings"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/cli"

	"github.com/fatih/color"
)

type PromptReader interface {
	ReadString(delim byte) (string, error)
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
	readValue, err := promptReader.ReadString('\n')
	if err != nil {
		return "", err
	}

	readValue = strings.Replace(readValue, "\n", "", -1)

	return readValue, nil
}
