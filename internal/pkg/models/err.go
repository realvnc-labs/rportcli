package models

import (
	"fmt"
	"strings"
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
	errs := make([]string, 0, len(er.Errors))
	for _, err := range er.Errors {
		if err.Code == "" && err.Detail == "" {
			errs = append(errs, err.Title)
		} else {
			errs = append(errs, fmt.Sprintf("%s, code: %s, details: %s", err.Title, err.Code, err.Detail))
		}
	}

	return strings.Join(errs, "\n")
}
