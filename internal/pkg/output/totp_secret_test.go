package output

import (
	"bytes"
	"testing"

	"github.com/cloudradar-monitoring/rportcli/internal/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestRenderTotPSecret(t *testing.T) {
	testCases := []struct {
		Format         string
		ExpectedOutput string
		ColCountToGive int
	}{
		{
			Format: FormatHuman,
			ExpectedOutput: `One time password secret

comment123
SECRET    QR FILE      
secret123 somefile.png 
`,
			ColCountToGive: 150,
		},
		{
			Format: FormatJSON,
			ExpectedOutput: `{"secret":"secret123","comment":"comment123","file":"somefile.png"}
`,
			ColCountToGive: 10,
		},
		{
			Format: FormatJSONPretty,
			ExpectedOutput: `{
  "secret": "secret123",
  "comment": "comment123",
  "file": "somefile.png"
}
`,
			ColCountToGive: 10,
		},
		{
			Format: FormatYAML,
			ExpectedOutput: `secret: secret123
comment: comment123
file: somefile.png
`,
			ColCountToGive: 10,
		},
	}

	key := &models.TotPSecretOutput{
		Secret:  "secret123",
		Comment: "comment123",
		File:    "somefile.png",
	}
	for _, testCase := range testCases {
		tc := testCase
		t.Run(testCase.Format, func(t *testing.T) {
			buf := &bytes.Buffer{}
			colCountToGive := tc.ColCountToGive
			tr := &TotPSecretRenderer{
				ColCountCalculator: func() int {
					return colCountToGive
				},
				Writer: buf,
				Format: tc.Format,
			}

			err := tr.RenderTotPSecret(key)
			assert.NoError(t, err)
			if err != nil {
				return
			}

			assert.Equal(
				t,
				tc.ExpectedOutput,
				buf.String(),
			)
		})
	}
}
