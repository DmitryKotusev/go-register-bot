package login

import (
	"bot-main/globalvars"
	"bot-main/models"
	modelerrors "bot-main/models/errors"
	test_utils "bot-main/tests/utils"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	globalvars.LoginRequestUrl = "https://fake/login"
	globalvars.LoginPageUrl = "https://fake/login/page"

	testCases := []struct {
		name        string
		client      *http.Client
		loginData   models.LoginData
		mockResp    *http.Response
		mockErr     error
		wantToken   string
		wantErrStr  string
		wantErrType any
	}{
		{
			name:       "nil client",
			client:     nil,
			loginData:  models.LoginData{},
			wantErrStr: "Login, HTTP client is nil",
		},
		{
			name: "successful login",
			client: test_utils.NewTestClient(func(req *http.Request) *http.Response {
				resp := models.LoginResponse{
					IsAuthSuccessful: true,
					Token:            "abc123",
				}
				respBytes, _ := json.Marshal(resp)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(respBytes)),
					Header:     make(http.Header),
				}
			}),
			loginData: models.LoginData{Email: "user@example.com", Password: "pass"},
			wantToken: "abc123",
		},
		{
			name: "invalid credentials",
			client: test_utils.NewTestClient(func(req *http.Request) *http.Response {
				code := "User_Email_Password-NotExists"
				resp := models.LoginResponse{
					IsAuthSuccessful: false,
					Code:             &code,
				}
				respBytes, _ := json.Marshal(resp)
				return &http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       io.NopCloser(bytes.NewReader(respBytes)),
					Header:     make(http.Header),
				}
			}),
			loginData:   models.LoginData{Email: "wrong@example.com", Password: "bad"},
			wantErrStr:  "Login failed because of wrong credentials, code User_Email_Password-NotExists",
			wantErrType: &modelerrors.InvalidCredentailsError{},
		},
		{
			name: "server error response",
			client: test_utils.NewTestClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     "500 Internal Server Error",
					Body:       io.NopCloser(bytes.NewReader([]byte("error"))),
				}
			}),
			loginData:  models.LoginData{},
			wantErrStr: "Login request failed with status: 500 Internal Server Error",
		},
		{
			name: "bad json response",
			client: test_utils.NewTestClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte("{bad json"))),
				}
			}),
			loginData:  models.LoginData{},
			wantErrStr: "Login body JSON parcing error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token, err := Login(tc.client, tc.loginData)

			if tc.wantErrStr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErrStr)
				if tc.wantErrType != nil {
					assert.Error(t, err)
					assert.Truef(t,
						errors.As(err, tc.wantErrType),
						"expected error of type %T but got %T",
						tc.wantErrType, err,
					)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantToken, token)
			}
		})
	}
}
