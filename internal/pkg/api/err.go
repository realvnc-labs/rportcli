package api

import (
	"encoding/json"
	"fmt"
)

type ErrorResp struct {
	Errors []struct {
		Code   string `json:"code"`
		Title  string `json:"title"`
		Detail string `json:"detail"`
	} `json:"errors"`
}

func (er ErrorResp) Error() string {
	jsonStr, err := json.Marshal(er)
	if err != nil {
		jsonStr = []byte{}
	}

	return fmt.Sprintf("API error response: '%s'", jsonStr)
}
