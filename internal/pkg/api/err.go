package api

import (
	"encoding/json"
	"fmt"
)

type Error struct {
	Code   string `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

type ErrorResp struct {
	Errors []Error `json:"errors"`
}

func (er ErrorResp) Error() string {
	jsonStr, err := json.Marshal(er)
	if err != nil {
		jsonStr = []byte{}
	}

	return fmt.Sprintf("API error response: '%s'", jsonStr)
}
