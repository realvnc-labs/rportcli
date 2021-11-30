package models

import (
	"strings"

	"github.com/breathbath/go_utils/v2/pkg/testing"
)

type Me struct {
	Username    string   `json:"username"`
	Groups      []string `json:"groups"`
	TwoFASendTo string   `json:"two_fa_send_to"`
}

func (m *Me) KeyValues() []testing.KeyValueStr {
	return []testing.KeyValueStr{
		{
			Key:   "Username",
			Value: m.Username,
		},
		{
			Key:   "TwoFactorAuthSentTo",
			Value: m.TwoFASendTo,
		},
		{
			Key:   "Groups",
			Value: strings.Join(m.Groups, ", "),
		},
	}
}
