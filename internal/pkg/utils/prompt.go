package utils

import "fmt"

type PromptReader struct {
}

func (pr *PromptReader) ReadString() (string, error) {
	var val = new(string)
	_, err := fmt.Scanln(val)
	if err != nil && err.Error() == "unexpected newline" {
		return "", nil
	}
	return *val, err
}

func (pr *PromptReader) ReadPassword() (string, error) {
	inputBytes, err := ReadPassword()
	return string(inputBytes), err
}
