package utils

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ErrorMock struct {
	Err string `json:"error"`
}

func (em *ErrorMock) Error() string {
	return em.Err
}

type SomeModel struct {
	Color string `json:"color"`
}

func TestBaseClient(t *testing.T) {
	testCases := []struct {
		respCodeToGive    int
		bodyToGive        string
		auth              *AuthMock
		target            interface{}
		errTarget         error
		expectedError     string
		expectedTarget    interface{}
		expectedErrTarget interface{}
	}{
		{
			respCodeToGive:    http.StatusOK,
			bodyToGive:        `{"color":"red"}`,
			auth:              &AuthMock{},
			target:            &SomeModel{},
			errTarget:         nil,
			expectedError:     "",
			expectedTarget:    &SomeModel{Color: "red"},
			expectedErrTarget: nil,
		},
		{
			respCodeToGive:    http.StatusInternalServerError,
			bodyToGive:        `{"error":"some error"}`,
			target:            &SomeModel{},
			errTarget:         &ErrorMock{},
			expectedError:     "some error",
			expectedTarget:    nil,
			expectedErrTarget: &ErrorMock{
				Err: "some error",
			},
		},
		{
			respCodeToGive:    http.StatusBadRequest,
			bodyToGive:        `{"error":"some error"}`,
			target:            &SomeModel{},
			errTarget:         nil,
			expectedError:     "invalid input provided",
			expectedTarget:    nil,
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(tc.respCodeToGive)
			_, e := rw.Write([]byte(tc.bodyToGive))
			if e != nil {
				assert.NoError(t, e)
			}
		}))

		bc := &BaseClient{}
		if testCase.auth != nil {
			bc.WithAuth(testCase.auth)
		}

		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, srv.URL, nil)
		assert.NoError(t, err)

		resp, err := bc.Call(req, testCase.target, testCase.errTarget)
		if testCase.expectedError != "" {
			assert.Error(t, err)
			assert.Contains(t, err.Error(), testCase.expectedError)
			continue
		}

		assert.NoError(t, err)
		if err != nil {
			continue
		}

		assert.Equal(t, testCase.respCodeToGive, resp.StatusCode)
		if testCase.auth != nil {
			assert.Len(t, testCase.auth.req, 1)
		}
		assert.Equal(t, testCase.expectedTarget, testCase.target)
	}
}
