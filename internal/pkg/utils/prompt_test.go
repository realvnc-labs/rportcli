package utils

import (
	"errors"
	"io"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockScanner struct {
	Txt     string
	ScanRes bool
	Error   error
}

func (ms MockScanner) Scan() bool {
	return ms.ScanRes
}

func (ms MockScanner) Text() string {
	return ms.Txt
}

func (ms MockScanner) Err() error {
	return ms.Error
}

func TestReadStringSuccess(t *testing.T) {
	sigChan := make(chan os.Signal, 1)
	s := MockScanner{
		Txt:     "some text",
		ScanRes: true,
	}

	pr := PromptReader{
		Sc:      s,
		SigChan: sigChan,
		PasswordScanner: func() (b []byte, err error) {
			return
		},
	}

	acualStr, err := pr.ReadString()
	assert.NoError(t, err)
	assert.Equal(t, "some text", acualStr)
}

func TestReadStringError(t *testing.T) {
	sigChan := make(chan os.Signal, 1)
	s := MockScanner{
		Txt:   "some text",
		Error: errors.New("some error"),
	}

	pr := PromptReader{
		Sc:      s,
		SigChan: sigChan,
		PasswordScanner: func() (b []byte, err error) {
			return
		},
	}

	acualStr, err := pr.ReadString()
	assert.EqualError(t, err, "some error")
	assert.Equal(t, "", acualStr)
}

func TestTermination(t *testing.T) {
	sigChan := make(chan os.Signal, 1)
	sigChan <- syscall.SIGTERM

	s := MockScanner{}

	pr := PromptReader{
		Sc:      s,
		SigChan: sigChan,
		PasswordScanner: func() (b []byte, err error) {
			return
		},
	}

	_, err := pr.ReadString()
	assert.EqualError(t, err, io.EOF.Error())
}

func TestPasswordReadingSuccess(t *testing.T) {
	sigChan := make(chan os.Signal, 1)
	s := MockScanner{}

	pr := PromptReader{
		Sc:      s,
		SigChan: sigChan,
		PasswordScanner: func() (b []byte, err error) {
			return []byte("somepass"), nil
		},
	}

	pass, err := pr.ReadPassword()
	assert.NoError(t, err)
	assert.Equal(t, "somepass", pass)
}

func TestPasswordReadingError(t *testing.T) {
	sigChan := make(chan os.Signal, 1)
	s := MockScanner{}

	pr := PromptReader{
		Sc:      s,
		SigChan: sigChan,
		PasswordScanner: func() (b []byte, err error) {
			return b, errors.New("some error")
		},
	}

	_, err := pr.ReadPassword()
	assert.EqualError(t, err, "some error")
}

func TestReadConfirmation(t *testing.T) {
	cases := []struct {
		name           string
		inputString    string
		expectedResult bool
	}{
		{
			name:           "Empty Response",
			inputString:    "",
			expectedResult: false,
		},
		{
			name:           "Simple Y",
			inputString:    "Y",
			expectedResult: true,
		},
		{
			name:           "Simple y",
			inputString:    "y",
			expectedResult: true,
		},
		{
			name:           "Simple N",
			inputString:    "N",
			expectedResult: false,
		},
		{
			name:           "Simple n",
			inputString:    "n",
			expectedResult: false,
		},
		{
			name:           "Complex y",
			inputString:    "yhhhnfhhfhfh",
			expectedResult: true,
		},
		{
			name:           "Complex n",
			inputString:    "nkshdfkh",
			expectedResult: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := MockScanner{
				Txt:     tc.inputString,
				ScanRes: true,
			}

			pr := PromptReader{
				Sc:              s,
				SigChan:         nil,
				PasswordScanner: nil,
			}

			result, err := pr.ReadConfirmation("Test Prompt: ")
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}
