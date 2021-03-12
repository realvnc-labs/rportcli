package output

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderByFormat(t *testing.T) {
	testCases := []struct {
		inputFormat    string
		source         interface{}
		expectedResult string
		expectedError  string
	}{
		{
			inputFormat:    "",
			source:         map[string]int{"one": 1, "two": 2, "three": 3},
			expectedResult: "default render result",
			expectedError:  "",
		},
		{
			inputFormat: FormatJSON,
			source:      []string{"one", "two", "three"},
			expectedResult: `["one","two","three"]
`,
			expectedError: "",
		},
		{
			inputFormat: FormatYAML,
			source:      []string{"one", "two", "three"},
			expectedResult: `- one
- two
- three
`,
			expectedError: "",
		},
		{
			inputFormat: FormatJSONPretty,
			source:      []string{"one", "two", "three"},
			expectedResult: `[
  "one",
  "two",
  "three"
]
`,
			expectedError: "",
		},
		{
			inputFormat:   "some unknown format",
			source:        []string{"one", "two", "three"},
			expectedError: "unknown rendering format: some unknown format",
		},
	}

	for _, testCase := range testCases {
		buf := &bytes.Buffer{}
		err := RenderByFormat(testCase.inputFormat, buf, testCase.source, func() error {
			_, e := buf.WriteString("default render result")
			return e
		})

		if testCase.expectedError != "" {
			assert.EqualError(t, err, testCase.expectedError)
			continue
		}

		assert.Equal(t, testCase.expectedResult, buf.String())
	}
}
