package cli_test

import (
	"net/http"
	"testing"

	"github.com/Dyleme/image-coverter/internal/cli"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	testCases := []struct {
		testName   string
		nickname   string
		password   string
		respStatus int
		respBody   string
		wantErr    error
	}{
		{
			testName:   "basic request",
			nickname:   "me",
			password:   "me",
			respStatus: 200,
			respBody:   "",
			wantErr:    nil,
		},
		{
			testName:   "already exists",
			nickname:   "me",
			password:   "me",
			respStatus: 409,
			respBody:   "",
			wantErr:    cli.WrongStatusError{409},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			httpmock.RegisterResponder("POST", "http://localhost:8080/auth/register",
				func(req *http.Request) (*http.Response, error) {
					resp := httpmock.NewStringResponse(tc.respStatus, tc.respBody)
					return resp, nil
				})

			gotErr := cli.Register("", tc.nickname, tc.password)
			assert.ErrorIs(t, gotErr, tc.wantErr)
		})
	}
}
