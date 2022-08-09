package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func PrettifyJSON(input []byte) string {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, input, "", "  ")
	if err == nil {
		return prettyJSON.String()
	}
	return fmt.Sprintf("%s: %s", err, input)
}
