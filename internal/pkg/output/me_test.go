package output

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestRenderMe(t *testing.T) {
	testCases := []struct {
		Format         string
		ExpectedOutput string
	}{
		{
			Format: FormatHuman,
			ExpectedOutput: `KEY                  VALUE          
Username:            myusr          
TwoFactorAuthSentTo: no@mail.me     
Groups:              group1, group2 
`,
		},
		{
			Format: FormatJSON,
			ExpectedOutput: `{"username":"myusr","groups":["group1","group2"],"two_fa_send_to":"no@mail.me"}
`,
		},
		{
			Format: FormatJSONPretty,
			ExpectedOutput: `{
  "username": "myusr",
  "groups": [
    "group1",
    "group2"
  ],
  "two_fa_send_to": "no@mail.me"
}
`,
		},
		{
			Format: FormatYAML,
			ExpectedOutput: `username: myusr
groups:
- group1
- group2
twofasendto: no@mail.me
`,
		},
	}

	me := &models.Me{
		Username:    "myusr",
		Groups:      []string{"group1", "group2"},
		TwoFASendTo: "no@mail.me",
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run("render_"+tc.Format, func(t *testing.T) {
			buf := &bytes.Buffer{}
			jr := &MeRenderer{
				Writer: buf,
				Format: tc.Format,
			}

			err := jr.RenderMe(me)
			require.NoError(t, err)

			assert.Equal(
				t,
				tc.ExpectedOutput,
				buf.String(),
			)
		})
	}
}
