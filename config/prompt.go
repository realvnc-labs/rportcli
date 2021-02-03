package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

// PromptRequiredValues will ask user for the list of required values
func PromptRequiredValues(missedRequirements []ParameterRequirement, targetKV map[string]string) error {
	reader := bufio.NewReader(os.Stdin)
	var err error
	for _, missedReq := range missedRequirements {
		readValue := ""
		if missedReq.Default != "" && missedReq.Validate == nil {
			readValue, err = promptValue(missedReq, reader)
			if err != nil {
				return err
			}
			if readValue != "" {
				targetKV[missedReq.Field] = readValue
			}
			continue
		}

		if missedReq.Validate == nil {
			continue
		}

		err = missedReq.Validate(missedReq.Field, readValue)
		for err != nil {
			readValue, err = promptValue(missedReq, reader)
			if err != nil {
				return err
			}

			err = missedReq.Validate(missedReq.Field, readValue)
			if err != nil {
				color.Red(err.Error())
			}
		}
		targetKV[missedReq.Field] = readValue
	}

	return nil
}

func promptValue(req ParameterRequirement, reader *bufio.Reader) (string, error) {
	fmt.Println(req.Help)
	if req.Default != "" {
		fmt.Printf("Default value: %s\n", req.Default)
	}

	fmt.Print("-> ")
	readValue, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	readValue = strings.Replace(readValue, "\n", "", -1)

	return readValue, nil
}
